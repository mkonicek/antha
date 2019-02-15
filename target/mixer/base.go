package mixer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/composer"
	driver "github.com/antha-lang/antha/driver/antha_driver_v1"
	"github.com/antha-lang/antha/driver/liquidhandling/client"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	lhdriver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
	"google.golang.org/grpc"
)

var subTypeToConnDriverFun = map[target.MixerDriverSubType](func(*grpc.ClientConn) lhdriver.LiquidhandlingDriver){
	target.GilsonPipetmaxSubType: func(conn *grpc.ClientConn) lhdriver.LiquidhandlingDriver {
		return client.NewLowLevelClientFromConn(conn)
	},
	target.CyBioSubType: func(conn *grpc.ClientConn) lhdriver.LiquidhandlingDriver {
		return client.NewLowLevelClientFromConn(conn)
	},
	target.LabcyteSubType: func(conn *grpc.ClientConn) lhdriver.LiquidhandlingDriver {
		return client.NewHighLevelClientFromConn(conn)
	},
}

type BaseMixer struct {
	id              workflow.DeviceInstanceID
	connection      workflow.ParsedConnection
	expectedSubType target.MixerDriverSubType

	logger *logger.Logger

	lock        sync.Mutex
	cmd         *exec.Cmd
	cmdFinished chan struct{}
	conn        *grpc.ClientConn
	properties  *lhdriver.LHProperties
}

func NewBaseMixer(logger *logger.Logger, id workflow.DeviceInstanceID, connection workflow.ParsedConnection, subType target.MixerDriverSubType) *BaseMixer {
	return &BaseMixer{
		id:              id,
		connection:      connection,
		expectedSubType: subType,
		logger:          logger.With("instructionPlugin", string(id)),
	}
}

func (bm *BaseMixer) connect(wf *workflow.Workflow, data []byte) error {
	if err := bm.maybeLinkedDriver(wf, data); err != nil {
		bm.Close()
		return err
	} else {
		bm.maybeExec()
		if err := bm.maybeDial(); err != nil {
			bm.Close()
			return err
		} else if err := bm.maybeConfigureConn(wf, data); err != nil {
			bm.Close()
			return err
		}
	}
	if bm.properties == nil {
		return fmt.Errorf("Unable to establish connection to mixer instructionPlugin for %v.", bm.id)
	} else {
		return nil
	}
}

// async. Blocks only until error on exec, or some data received from
// cmd's stdout or stderr, whichever is soonist.
func (bm *BaseMixer) maybeExec() {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	if bm.cmd == nil && bm.connection.ExecFile != "" {
		rng := rand.New(rand.NewSource(time.Now().Unix()))
		port := fmt.Sprint(1024 + rng.Intn(65536-1024))
		bm.cmd = exec.Command(bm.connection.ExecFile, "-port", port)
		bm.cmd.Env = []string{}
		bm.connection.HostPort = "localhost:" + port

		running := make(chan struct{})
		bm.cmdFinished = make(chan struct{})
		// copy it out so we don't have a data race: we won't be holding
		// the lock later in the go-routine when we close the cmdFinished chan.
		cmdFinished := bm.cmdFinished
		go func() {
			// We have to be careful here: we want to wait until we get
			// something out of either stdout or stderr, which of course
			// could be concurrent, hence the locking and careful
			// signalling.
			lock := new(sync.Mutex)
			logFun := func(pairs ...interface{}) error {
				lock.Lock()
				defer lock.Unlock()
				select {
				case <-running:
				default:
					close(running)
				}
				return bm.logger.Log(pairs...)
			}
			err := composer.RunAndLogCommand(bm.cmd, logFun)
			close(cmdFinished)
			if err != nil {
				bm.logger.Log("error", err)
			}
			bm.Close() // this is why Close() must be idempotent and thread safe!
		}()
		<-running
	}
}

func (bm *BaseMixer) maybeDial() error {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	if bm.conn == nil && bm.connection.HostPort != "" {
		bm.logger.Log("dialing", bm.connection.HostPort)
		conn, err := grpc.Dial(bm.connection.HostPort, grpc.WithInsecure())
		if err != nil {
			return err
		}
		bm.conn = conn
		c := driver.NewDriverClient(conn)
		ctx := context.Background()
		if reply, err := c.DriverType(ctx, &driver.TypeRequest{}); err != nil {
			return err
		} else if typ := reply.GetType(); typ != target.DriverSelectorV1Mixer.Value {
			return fmt.Errorf("Expected to find a mixer instructionPlugin at %s but instead found: %s", bm.connection, typ)
		} else if subtypes := reply.GetSubtypes(); len(subtypes) != 1 || subtypes[0] != string(bm.expectedSubType) {
			return fmt.Errorf("Expected to find a [%v] mixer instructionPlugin at %s but instead found: %v", bm.expectedSubType, bm.connection, subtypes)
		}
	}
	return nil
}

func (bm *BaseMixer) maybeConfigureConn(wf *workflow.Workflow, data []byte) error {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	if bm.conn != nil && bm.properties == nil {
		if fun, found := subTypeToConnDriverFun[bm.expectedSubType]; !found {
			return fmt.Errorf("Unable to find connection function for mixer subtype %v", bm.expectedSubType)
		} else {
			driver := fun(bm.conn)
			if props, status := driver.Configure(wf.JobId, wf.Meta.Name, bm.id, data); !status.Ok() {
				return status.GetError()
			} else {
				props.Driver = driver
				bm.properties = props
				return nil
			}
		}
	}
	return nil
}

func (bm *BaseMixer) Close() {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	bm.properties = nil

	if bm.conn != nil {
		bm.conn.Close()
		bm.conn = nil
	}

	if bm.cmd != nil {
		// copy it out to avoid data race
		proc := bm.cmd.Process
		if proc != nil {
			go func() {
				// these signal calls will fail on some OS, and will likely
				// fail if the process has already exited.
				proc.Signal(syscall.SIGTERM)
				// give it 1 second to shut down cleanly, then just kill it
				// hard.
				time.Sleep(time.Second)
				proc.Kill()
			}()
		}
		<-bm.cmdFinished
		bm.cmd = nil
		bm.cmdFinished = nil
	}
}

func (bm *BaseMixer) CanCompile(req effects.Request) bool {
	if bm.properties == nil {
		panic("CanCompile called without an active connection to instructionPlugin")
	}
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Mixer,
		},
	}
	if bm.properties.CanPrompt() {
		can.Selector = append(can.Selector, target.DriverSelectorV1Prompter)
	}
	return can.Contains(req)
}

func mix(labEffects *effects.LaboratoryEffects, instrs []*wtype.LHInstruction, req *liquidhandling.LHRequest, props *lhdriver.LHProperties) ([]effects.Inst, error) {
	if err := addCustomPolicies(instrs, req); err != nil {
		return nil, err
	}

	hasOutputPlate := func(typ wtype.PlateTypeName, id string) bool {
		for _, p := range req.OutputPlatetypes {
			if p.Type == typ && (id == "" || p.ID == id) {
				return true
			}
		}
		return false
	}

	for _, instr := range instrs {
		if instr.OutPlate != nil {
			if p, found := req.OutputPlates[instr.OutPlate.ID]; found && p != instr.OutPlate {
				return nil, fmt.Errorf("Mix setup error: Plate %s already requested in different state for mix.", p.ID)
			} else {
				req.OutputPlates[instr.OutPlate.ID] = instr.OutPlate
			}
		}

		if len(instr.Platetype) != 0 && !hasOutputPlate(instr.Platetype, instr.PlateID) {
			if pt, err := labEffects.Inventory.PlateTypes.NewPlate(instr.Platetype); err != nil {
				return nil, err
			} else {
				pt.ID = instr.PlateID
				req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			}
		}
		req.Add_instruction(instr)
	}

	if err := liquidhandling.Init(props).MakeSolutions(labEffects, req); err != nil {
		return nil, err
	}

	return nil, nil // FIXME
}

func mergePolicies(basePolicy, priorityPolicy wtype.LHPolicy) (newPolicy wtype.LHPolicy) {
	newPolicy = make(wtype.LHPolicy)

	for key, value := range priorityPolicy {
		newPolicy[key] = value
	}

	for key, value := range basePolicy {
		if _, found := priorityPolicy[key]; !found {
			newPolicy[key] = value
		}
	}
	return newPolicy
}

// any customised user policies are added to the LHRequest PolicyManager here
// Any component type names with modified policies are iterated until unique i.e. SmartMix_modified_1
func addCustomPolicies(mixes []*wtype.LHInstruction, lhreq *liquidhandling.LHRequest) error {
	systemPolicyRuleSet := lhreq.GetPolicyManager().Policies()
	systemPolicies := systemPolicyRuleSet.Policies
	var userPolicies = make(map[string]wtype.LHPolicy)
	var allPolicies = make(map[string]wtype.LHPolicy)
	var liquidClassConversionMap = make(map[string]string)

	for key, value := range systemPolicies {
		allPolicies[key] = value
	}

	userPolicyRuleSet := wtype.NewLHPolicyRuleSet()

	for _, mixInstruction := range mixes {
		for _, component := range mixInstruction.Inputs {
			if len(component.Policy) > 0 {
				if matchingSystemPolicy, found := allPolicies[string(component.Type)]; found {
					mergedPolicy := mergePolicies(matchingSystemPolicy, component.Policy)
					if !wtype.EquivalentPolicies(mergedPolicy, matchingSystemPolicy) {
						num := 1
						newPolicyName := makemodifiedTypeName(component.Type, num)
						existingCustomPolicy, found := allPolicies[newPolicyName]
						for found {
							// check if existing policy with modified name is the same
							if !wtype.EquivalentPolicies(mergedPolicy, existingCustomPolicy) {
								// if not increase number and try again
								num++
								newPolicyName = makemodifiedTypeName(component.Type, num)
								existingCustomPolicy, found = allPolicies[newPolicyName]
							} else {
								// otherwise use existing modified policy
								found = false
							}
						}
						allPolicies[newPolicyName] = mergedPolicy
						userPolicies[newPolicyName] = mergedPolicy
						component.Type = wtype.LiquidType(newPolicyName)
						liquidClassConversionMap[newPolicyName] = matchingSystemPolicy.Name()
					}
				} else {
					allPolicies[string(component.Type)] = component.Policy
					userPolicies[string(component.Type)] = component.Policy
				}
			}
		}
	}

	if len(userPolicies) > 0 {
		userPolicyRuleSet, err := wtype.AddUniversalRules(userPolicyRuleSet, userPolicies)
		if err != nil {
			return err
		}
		for newClass, original := range liquidClassConversionMap {
			err := wtype.CopyRulesFromPolicy(userPolicyRuleSet, original, newClass)
			if err != nil {
				return err
			}
		}
		lhreq.AddUserPolicies(userPolicyRuleSet)
	}

	return nil
}

const modifiedPolicySuffix = "_modified_"

func makemodifiedTypeName(componentType wtype.LiquidType, number int) string {
	return string(componentType) + modifiedPolicySuffix + strconv.Itoa(number)
}

// unModifyTypeName will trim a _modified_ suffix from a LiquidType in the CSV file.
// These are added to LiquidType names when a Liquid is modified in an element.
func unModifyTypeName(componentType string) string {
	return strings.Split(componentType, modifiedPolicySuffix)[0]
}

func floatValue(a, b *float64) float64 {
	if a != nil {
		return *a
	} else {
		return *b
	}
}

func checkInstructions(nodes []effects.Node) ([]*wtype.LHInstruction, error) {
	instrs := make([]*wtype.LHInstruction, 0, len(nodes))
	for _, node := range nodes {
		if cmd, ok := node.(*effects.Command); !ok {
			return nil, fmt.Errorf("cannot compile %T", node)
		} else if instr, ok := cmd.Inst.(*wtype.LHInstruction); !ok {
			return nil, fmt.Errorf("cannot compile %T", cmd.Inst)
		} else {
			instrs = append(instrs, instr)
		}
	}
	if len(instrs) == 0 {
		return nil, errors.New("No instructions to mix!")
	} else {
		return instrs, nil
	}
}
