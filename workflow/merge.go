package workflow

import (
	"errors"
	"fmt"
	"sort"

	"github.com/antha-lang/antha/utils"
)

func (a *Workflow) Merge(b *Workflow) error {
	switch {
	case a == nil:
		return errors.New("Cannot merge into a nil Workflow")
	case b == nil:
		return nil
	case a.JobId == "":
		a.JobId = b.JobId
	case a.JobId != b.JobId && b.JobId != "":
		return fmt.Errorf("Cannot merge: different JobIds: %v vs %v", a.JobId, b.JobId)
	}

	return utils.ErrorSlice{
		b.SchemaVersion.Validate(), // every snippet must have a valid SchemaVersion
		a.Repositories.Merge(b.Repositories),
		a.Elements.merge(b.Elements),
		a.Inventory.merge(b.Inventory),
		a.Config.merge(b.Config),
	}.Pack()
}

func (a Repositories) Merge(b Repositories) error {
	for repoName, repoB := range b {
		if repoA, found := a[repoName]; found {
			// We're ok if:
			// .Directory in both are equal (includes both being "")
			// repoA.Directory is "" (we just copy in from repoB.Directory)
			// If both != "" and repoA.Directory != repoB.Directory then error
			if dir, ok := tryMergeStrings(repoA.Directory, repoB.Directory); !ok {
				return fmt.Errorf("Cannot merge: repository with name '%v' redefined illegally (Directory fields not empty and not equal: %s vs %s).",
					repoName, repoA.Directory, repoB.Directory)
			} else {
				repoA.Directory = dir
			}

			// Merge needs to work even in the absence of the Directory
			// field, which means we really can't start using git at this
			// point.  Consequently, all we can do safely here is basic
			// string equality and allow empty fields to be set. This is
			// slightly more restrictive than we might want, for example,
			// it could be valid to have two different branch names that
			// happen to resolve to the same commit, but we can't test
			// for that here so we play it safe and enforce strict
			// equality.

			if commit, ok := tryMergeStrings(repoA.Commit, repoB.Commit); !ok {
				return fmt.Errorf("Cannot merge: repository with name '%v' redefined illegally (Commit fields not empty and not equal; %s vs %s).",
					repoName, repoA.Commit, repoB.Commit)
			} else {
				repoA.Commit = commit
			}

			if branch, ok := tryMergeStrings(repoA.Branch, repoB.Branch); !ok {
				return fmt.Errorf("Cannot merge: repository with name '%v' redefined illegally (Branch fields not empty and not equal; %s vs %s).",
					repoName, repoA.Branch, repoB.Branch)
			} else {
				repoA.Branch = branch
			}

		} else if !found {
			a[repoName] = repoB
		}
	}
	return nil
}

func tryMergeStrings(a, b string) (string, bool) {
	if a != "" && b != "" && a != b {
		return "", false
	} else if a == "" {
		return b, true
	} else {
		return a, true
	}
}

func (a *Elements) merge(b Elements) error {
	return utils.ErrorSlice{
		a.Types.merge(b.Types),
		a.Instances.merge(b.Instances),
		a.InstancesConnections.merge(b.InstancesConnections),
	}.Pack()
}

func (ets ElementTypes) sort() {
	sort.Slice(ets, func(i, j int) bool {
		return ets[i].lessThan(ets[j])
	})
}

func (a *ElementType) lessThan(b *ElementType) bool {
	return a.RepositoryName < b.RepositoryName ||
		(a.RepositoryName == b.RepositoryName && a.ElementPath < b.ElementPath)
}

func (a *ElementType) equals(b *ElementType) bool {
	return a.RepositoryName == b.RepositoryName && a.ElementPath == b.ElementPath
}

func (a *ElementTypes) merge(b ElementTypes) error {
	// for convenience, it's perfectly reasonable to have the same
	// element types defined in multiple places, and we just need to
	// check that duplicates are truly equal
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
		a.GlobalMixer.Merge(b.GlobalMixer),
		a.GilsonPipetMax.Merge(b.GilsonPipetMax),
		a.Tecan.Merge(b.Tecan),
		a.CyBio.Merge(b.CyBio),
		a.Labcyte.Merge(b.Labcyte),
		a.Hamilton.Merge(b.Hamilton),
		a.QPCR.Merge(b.QPCR),
		a.ShakerIncubator.Merge(b.ShakerIncubator),
		a.PlateReader.Merge(b.PlateReader),
	}.Pack()
}

func (a *GilsonPipetMaxConfig) Merge(b GilsonPipetMaxConfig) error {
	// simplest: we merge iff device ids are distinct
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: GilsonPipetMax device '%v' redefined", id)
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

func (a *TecanConfig) Merge(b TecanConfig) error {
	// simplest: we merge iff device ids are distinct
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: Tecan device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}

	// for defaults, we just hope one of them is nil
	switch {
	case a.Defaults == nil:
		a.Defaults = b.Defaults
	case b.Defaults != nil:
		return errors.New("Cannot merge: Tecan defaults redefined")
	}
	return nil
}

func (a *CyBioConfig) Merge(b CyBioConfig) error {
	// simplest: we merge iff device ids are distinct
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: CyBio device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}

	// for defaults, we just hope one of them is nil
	switch {
	case a.Defaults == nil:
		a.Defaults = b.Defaults
	case b.Defaults != nil:
		return errors.New("Cannot merge: CyBio defaults redefined")
	}
	return nil
}

func (a *LabcyteConfig) Merge(b LabcyteConfig) error {
	// simplest: we merge iff device ids are distinct
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: Labcyte device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}

	// for defaults, we just hope one of them is nil
	switch {
	case a.Defaults == nil:
		a.Defaults = b.Defaults
	case b.Defaults != nil:
		return errors.New("Cannot merge: Labcyte defaults redefined")
	}
	return nil
}

func (a *HamiltonConfig) Merge(b HamiltonConfig) error {
	// simplest: we merge iff device ids are distinct
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("cannot merge: Hamilton device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}

	// for defaults, we just hope one of them is nil
	switch {
	case a.Defaults == nil:
		a.Defaults = b.Defaults
	case b.Defaults != nil:
		return errors.New("cannot merge: Hamilton defaults redefined")
	}
	return nil
}

func (a *GlobalMixerConfig) Merge(b GlobalMixerConfig) error {
	// disjunction of bools - this seems sensible
	a.PrintInstructions = a.PrintInstructions || b.PrintInstructions
	a.UseDriverTipTracking = a.UseDriverTipTracking || b.UseDriverTipTracking
	a.IgnorePhysicalSimulation = a.IgnorePhysicalSimulation || b.IgnorePhysicalSimulation

	// but we can't allow input plates to be speficied multiple times:
	switch inAEmpty, inBEmpty := len(a.InputPlates) == 0, len(b.InputPlates) == 0; {
	case inAEmpty:
		a.InputPlates = b.InputPlates
	case !inBEmpty: // by implication, also !inAEmpty
		return errors.New("Cannot merge: Config.GlobalMixer.InputPlates specified in multiple workflows. This is illegal")
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

func (a *QPCRConfig) Merge(b QPCRConfig) error {
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: QPCR device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}
	return nil
}

func (a *ShakerIncubatorConfig) Merge(b ShakerIncubatorConfig) error {
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: ShakerIncubator device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}
	return nil
}

func (a *PlateReaderConfig) Merge(b PlateReaderConfig) error {
	for id, cfg := range b.Devices {
		if _, found := a.Devices[id]; found {
			return fmt.Errorf("Cannot merge: PlateReader device '%v' redefined", id)
		}
		a.Devices[id] = cfg
	}
	return nil
}
