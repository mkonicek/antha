package mixer

import (
	"context"
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
	"github.com/antha-lang/antha/logger"
	lhdriver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/workflow"
	"google.golang.org/grpc"
)

type MixerDriverSubType string

const (
	GilsonPipetmaxSubType MixerDriverSubType = "GilsonPipetmax"
)

var subTypeToConnDriverFun = map[MixerDriverSubType](func(*grpc.ClientConn) lhdriver.LiquidhandlingDriver){
	GilsonPipetmaxSubType: func(conn *grpc.ClientConn) lhdriver.LiquidhandlingDriver {
		return client.NewLowLevelClientFromConn(conn)
	},
}

type BaseMixer struct {
	id              workflow.DeviceInstanceID
	connection      workflow.ParsedConnection
	expectedSubType MixerDriverSubType

	logger *logger.Logger

	lock        sync.Mutex
	cmd         *exec.Cmd
	cmdFinished chan struct{}
	conn        *grpc.ClientConn
	properties  *lhdriver.LHProperties
}

func NewBaseMixer(logger *logger.Logger, id workflow.DeviceInstanceID, connection workflow.ParsedConnection, subType MixerDriverSubType) *BaseMixer {
	return &BaseMixer{
		id:              id,
		connection:      connection,
		expectedSubType: subType,
		logger:          logger.With("instructionPlugin", string(id)),
	}
}

func (bm *BaseMixer) Connect(wf *workflow.Workflow) (*lhdriver.LHProperties, error) {
	if err := bm.maybeLinkedDriver(wf); err != nil {
		bm.Close()
		return nil, err
	} else {
		bm.maybeExec()
		if err := bm.maybeDial(); err != nil {
			bm.Close()
			return nil, err
		} else if err := bm.maybeConfigure(wf); err != nil {
			bm.Close()
			return nil, err
		}
	}
	if bm.properties == nil {
		return nil, fmt.Errorf("Unable to establish connection to mixer instructionPlugin for %v.", bm.id)
	} else {
		return bm.properties, nil
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
		} else if typ := reply.GetType(); typ != "antha.mixer.v1.Mixer" {
			return fmt.Errorf("Expected to find a mixer instructionPlugin at %s but instead found: %s", bm.connection, typ)
		} else if subtypes := reply.GetSubtypes(); len(subtypes) != 1 || subtypes[0] != string(bm.expectedSubType) {
			return fmt.Errorf("Expected to find a [%v] mixer instructionPlugin at %s but instead found: %v", bm.expectedSubType, bm.connection, subtypes)
		}
	}
	return nil
}

func (bm *BaseMixer) maybeConfigure(wf *workflow.Workflow) error {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	if bm.conn != nil && bm.properties == nil {
		if fun, found := subTypeToConnDriverFun[bm.expectedSubType]; !found {
			return fmt.Errorf("Unable to find connection function for mixer subtype %v", bm.expectedSubType)
		} else {
			driver := fun(bm.conn)
			if props, status := driver.Configure(wf.JobId, wf.Meta.Name, bm.id); !status.Ok() {
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
