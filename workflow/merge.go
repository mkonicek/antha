package workflow

import (
	"errors"
	"fmt"
	"sort"

	"github.com/antha-lang/antha/utils"
)

func (a *Workflow) merge(b *Workflow) error {
	if a == nil || b == nil {
		return errors.New("Attempt to merge nil workflow")
	}

	if a.JobId == "" {
		a.JobId = b.JobId
	} else if b.JobId != "" && a.JobId != b.JobId {
		return fmt.Errorf("Cannot merge: different JobIds: %v vs %v", a.JobId, b.JobId)
	}

	return utils.ErrorSlice{
		a.Repositories.merge(b.Repositories),
		a.ElementTypes.merge(b.ElementTypes),
		a.ElementInstances.merge(b.ElementInstances),
		a.ElementInstancesParameters.merge(b.ElementInstancesParameters),
		a.ElementInstancesConnections.merge(b.ElementInstancesConnections),
		a.Inventory.merge(b.Inventory),
		a.Config.merge(b.Config),
	}.Pack()
}

func (a Repositories) merge(b Repositories) error {
	// It's an error iff a and b contain the same prefix with different definitions.
	for prefix, repoB := range b {
		if repoA, found := a[prefix]; found && repoA != repoB {
			return fmt.Errorf("Cannot merge: repository with prefix '%v' redefined.", prefix)
		} else if !found {
			a[prefix] = repoB
		}
	}
	return nil
}

func (ets ElementTypes) sort() {
	sort.Slice(ets, func(i, j int) bool {
		return ets[i].lessThan(ets[j])
	})
}

func (a *ElementType) lessThan(b *ElementType) bool {
	return a.RepositoryPrefix < b.RepositoryPrefix ||
		(a.RepositoryPrefix == b.RepositoryPrefix && a.ElementPath < b.ElementPath)
}

func (a *ElementType) equals(b *ElementType) bool {
	return a.RepositoryPrefix == b.RepositoryPrefix && a.ElementPath == b.ElementPath
}

func (a *ElementTypes) merge(b ElementTypes) error {
	all := make(ElementTypes, 0, len(*a)+len(b))
	all = append(all, *a...)
	all = append(all, b...)
	all.sort()

	result := make(ElementTypes, 0, len(all))
	old := ElementType{}
	for _, cur := range all {
		if old.equals(cur) {
			continue
		} else {
			result = append(result, cur)
			old = *cur
		}
	}
	*a = result
	return nil
}

func (a ElementInstances) merge(b ElementInstances) error {
	// Element instances from different workflows must be entirely distinct
	for name, elemB := range b {
		if _, found := a[name]; found {
			return fmt.Errorf("Cannot merge: element instance '%v' exists in both workflows", name)
		} else {
			a[name] = elemB
		}
	}
	return nil
}

func (a ElementInstancesParameters) merge(b ElementInstancesParameters) error {
	// Just like element instances, these should be completely distinct
	for name, paramSetB := range b {
		if _, found := a[name]; found {
			return fmt.Errorf("Cannot merge: element parameters '%v' exists in both workflows", name)
		} else {
			a[name] = paramSetB
		}
	}
	return nil
}

func (conns ElementInstancesConnections) sort() {
	sort.Slice(conns, func(i, j int) bool {
		return conns[i].lessThan(conns[j])
	})
}

func (a ElementConnection) lessThan(b ElementConnection) bool {
	return a.Source.lessThan(b.Source) ||
		(!b.Source.lessThan(a.Source) && a.Target.lessThan(b.Target))
}

func (a ElementSocket) lessThan(b ElementSocket) bool {
	return a.ElementInstance < b.ElementInstance ||
		(a.ElementInstance == b.ElementInstance && a.ParameterName < b.ParameterName)
}

func (a *ElementInstancesConnections) merge(b ElementInstancesConnections) error {
	all := make(ElementInstancesConnections, 0, len(*a)+len(b))
	all = append(all, *a...)
	all = append(all, b...)
	all.sort()

	result := make(ElementInstancesConnections, 0, len(all))
	old := ElementConnection{}
	for _, cur := range all {
		if old == cur { // structural equality
			return fmt.Errorf("Cannot merge: element connection '%v' exists in both workflows", cur)
		} else {
			result = append(result, cur)
			old = cur
		}
	}
	*a = result
	return nil
}

func (a *Inventory) merge(b Inventory) error {
	return a.PlateTypes.Merge(b.PlateTypes)
}

func (a *Config) merge(b Config) error {
	return utils.ErrorSlice{
		a.GilsonPipetMax.merge(b.GilsonPipetMax),
		a.GlobalMixer.merge(b.GlobalMixer),
	}.Pack()
}

func (a *GilsonPipetMaxConfig) merge(b GilsonPipetMaxConfig) error {
	// simplest: we merge iff device ids are distinct
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: GilsonPipetMax device '%s' redefined", id)
		}
		a.Devices[id] = cfg
	}
	// for defaults, we just hope one of them is nil
	switch {
	case a.Defaults == nil:
		a.Defaults = b.Defaults
	case b.Defaults != nil:
		return errors.New("Cannot merge: GilsonPipetMax defaults redefined")
	}
	return nil
}

func (a *GlobalMixerConfig) merge(b GlobalMixerConfig) error {
	// disjunction of bools - this seems sensible
	a.PrintInstructions = a.PrintInstructions || b.PrintInstructions
	a.UseDriverTipTracking = a.UseDriverTipTracking || b.UseDriverTipTracking
	a.IgnorePhysicalSimulation = a.IgnorePhysicalSimulation || b.IgnorePhysicalSimulation

	// but we can't allow input or output plates to be speficied multiple times:
	switch inAEmpty, inBEmpty := len(a.InputPlates) == 0, len(b.InputPlates) == 0; {
	case inAEmpty:
		a.InputPlates = b.InputPlates
	case !inBEmpty: // by implication, also !inAEmpty
		return errors.New("Cannot merge: Config.GlobalMixer.InputPlates specified in multiple workflows. This is illegal")
	}

	switch outAEmpty, outBEmpty := len(a.OutputPlates) == 0, len(b.OutputPlates) == 0; {
	case outAEmpty:
		a.OutputPlates = b.OutputPlates
	case !outBEmpty: // by implication, also !outAEmpty
		return errors.New("Cannot merge: Config.GlobalMixer.OutputPlates specified in multiple workflows. This is illegal")
	}

	// for LHPolicyRuleSets, there's already a merge function!
	switch {
	case a.CustomPolicyRuleSet == nil:
		a.CustomPolicyRuleSet = b.CustomPolicyRuleSet
	case b.CustomPolicyRuleSet != nil:
		a.CustomPolicyRuleSet.MergeWith(b.CustomPolicyRuleSet)
	}
	return nil
}
