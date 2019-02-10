package mixer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	driver "github.com/antha-lang/antha/driver/antha_driver_v1"
	"github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"google.golang.org/grpc"
)

type BaseMixer struct {
	connection       string
	expectedSubTypes []string
}

func NewBaseMixer(connection string, subTypes ...string) *BaseMixer {
	return &BaseMixer{
		connection:       connection,
		expectedSubTypes: subTypes,
	}
}

func (bm *BaseMixer) ConnectInit() (*grpc.ClientConn, error) {
	if bm.connection == "" {
		return nil, nil

	} else {
		conn, err := grpc.Dial(bm.connection, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		c := driver.NewDriverClient(conn)
		ctx := context.Background()
		if reply, err := c.DriverType(ctx, &driver.TypeRequest{}); err != nil {
			return nil, err
		} else if typ := reply.GetType(); typ != "antha.mixer.v1.Mixer" {
			return nil, fmt.Errorf("Expected to find a mixer driver at %s but instead found: %s", bm.connection, typ)
		} else if subtypes := reply.GetSubtypes(); len(subtypes) != len(bm.expectedSubTypes) {
			return nil, fmt.Errorf("Expected to find a %v mixer driver at %s but instead found: %v", bm.expectedSubTypes, bm.connection, subtypes)
		} else {
			for idx, est := range bm.expectedSubTypes {
				if subtypes[idx] != est {
					return nil, fmt.Errorf("Expected to find a %v mixer driver at %s but instead found: %v", bm.expectedSubTypes, bm.connection, subtypes)
				}
			}
			return conn, nil
		}
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
