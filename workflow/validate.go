package workflow

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/antha-lang/antha/utils"
)

func (wf *Workflow) validate() error {
	if wf.JobId == "" {
		return errors.New("Validation error: Workflow has empty JobId")
	} else {
		return utils.ErrorSlice{
			wf.Repositories.validate(),
			wf.ElementTypes.validate(wf),
			wf.ElementInstances.validate(wf),
			wf.ElementInstancesParameters.validate(wf),
			wf.ElementInstancesConnections.validate(wf),
			wf.Inventory.validate(),
			wf.Config.validate(),
		}.Pack()
	}
}

func (rs Repositories) validate() error {
	if len(rs) == 0 {
		return errors.New("Validation error: Workflow has no Repositories")
	} else {
		// Until we switch to go modules, we have to enforce that all
		// repositories are not only unique, but that no one repository
		// is a prefix of another. To enforce this, we sort the prefixes
		// (so shortest will come first) and then we need to only test
		// against the tail of the list.
		prefixes := make([]string, 0, len(rs))
		for prefix := range rs {
			prefixes = append(prefixes, string(prefix))
		}
		sort.Strings(prefixes)
		// Yes there's probably some algorithm to make this even more
		// efficient, but for now we're only dealing with a very small
		// number of repos, so less code and simpler code wins.
		for idx, prefix := range prefixes {
			for _, later := range prefixes[idx+1:] {
				if strings.HasPrefix(later, prefix) {
					return fmt.Errorf("Validation error: Two repositories found where one is a prefix of the other. This is not allowed, sorry. '%s' is a prefix of '%s'", prefix, later)
				}
			}
		}

		for _, repo := range rs {
			if err := repo.validate(); err != nil {
				return err
			}
		}
		return nil
	}
}

func (r Repository) validate() error {
	if info, err := os.Stat(filepath.FromSlash(r.Directory)); err != nil {
		return err
	} else if !info.Mode().IsDir() {
		return fmt.Errorf("Validation error: Repository Directory is not a directory: '%s'", r.Directory)
	} else if bEmpty, cEmpty := r.Branch == "", r.Commit == ""; !bEmpty && !cEmpty {
		return fmt.Errorf("Validation error: Repository cannot have both Branch and Commit specified. At most one. ('%s', '%s')", r.Branch, r.Commit)
	} else {
		return nil
	}
}

func (ets ElementTypes) validate(wf *Workflow) error {
	// we don't support import aliasing for elements. This means that
	// we require that every element type has a unique type name.
	namesToPath := make(map[ElementTypeName]ElementPath, len(ets))
	for _, et := range ets {
		if err := et.validate(wf); err != nil {
			return err
		} else if ep, found := namesToPath[et.Name()]; found {
			return fmt.Errorf("Validation error: ElementType '%v' is ambiguous (ElementPaths '%v' and '%v')", et.Name(), et.ElementPath, ep)
		} else {
			namesToPath[et.Name()] = et.ElementPath
		}
	}
	return nil
}

func (et ElementType) validate(wf *Workflow) error {
	if _, found := wf.Repositories[et.RepositoryPrefix]; !found {
		return fmt.Errorf("Validation error: ElementType uses unknown RepositoryPrefix: '%s'", et.RepositoryPrefix)
	} else {
		return nil
	}
}

func (eis ElementInstances) validate(wf *Workflow) error {
	for name, ei := range eis {
		if name == "" {
			return errors.New("Validation error: ElementInstance cannot have an empty name")
		} else if err := ei.validate(wf); err != nil {
			return err
		}
	}
	return nil
}

func (ei ElementInstance) validate(wf *Workflow) error {
	if _, found := wf.TypeNames()[ei.ElementTypeName]; !found {
		return fmt.Errorf("Validation error: ElementInstance with ElementTypeName '%v' is unknown", ei.ElementTypeName)
	} else {
		return nil
	}
}

func (eps ElementInstancesParameters) validate(wf *Workflow) error {
	for name, _ := range eps {
		if _, found := wf.ElementInstances[name]; !found {
			return fmt.Errorf("Validation error: ElementInstancesParameters provided for unknown ElementInstance '%v'", name)
		}
	}
	return nil
}

func (conns ElementInstancesConnections) validate(wf *Workflow) error {
	for _, conn := range conns {
		if err := conn.Source.validate(wf); err != nil {
			return err
		} else if err := conn.Target.validate(wf); err != nil {
			return err
		}
	}
	return nil
}

func (soc ElementSocket) validate(wf *Workflow) error {
	if _, found := wf.ElementInstances[soc.ElementInstance]; !found {
		return fmt.Errorf("Validation error: ElementConnection uses ElementInstance '%v' which does not exist.", soc.ElementInstance)
	} else if soc.ParameterName == "" {
		return fmt.Errorf("Validation error: ElementConnection using ElementInstance '%v' must specify a ParameterName.", soc.ElementInstance)
	} else {
		return nil
	}
}

func (inv Inventory) validate() error {
	return inv.PlateTypes.Validate()
}

func (cfg Config) validate() error {
	// NB the validation here is purely static - i.e. we're not
	// attempting to connect to any device plugins at this stage.
	return utils.ErrorSlice{
		cfg.GlobalMixer.validate(),
		cfg.GilsonPipetMax.validate(),
		cfg.assertOnlyOneMixer(),
	}.Pack()
}

func (cfg Config) assertOnlyOneMixer() error {
	// remove / revise when we get better
	if len(cfg.GilsonPipetMax.Devices) > 1 {
		return fmt.Errorf("Currently a maximum of one mixer can be used per workflow. You have %d configured.", len(cfg.GilsonPipetMax.Devices))
	}
	return nil
}

func (global GlobalMixerConfig) validate() error {
	for idx, p := range global.InputPlates {
		if p == nil {
			return fmt.Errorf("GlobalMixer contains illegal nil input plate at index %d", idx)
		}
	}
	for idx, p := range global.OutputPlates {
		if p == nil {
			return fmt.Errorf("GlobalMixer contains illegal nil Output plate at index %d", idx)
		}
	}
	// We cannot validate plates and plate types until we have a
	// working inventory system.
	return nil
}

func (gilsons GilsonPipetMaxConfig) validate() error {
	if err := gilsons.Defaults.validate("Defaults"); err != nil {
		return err
	}
	for id, inst := range gilsons.Devices {
		if err := inst.validate(id); err != nil {
			return err
		}
	}
	return nil
}

func (inst *GilsonPipetMaxInstanceConfig) validate(id DeviceInstanceID) error {
	if len(id) == 0 {
		return errors.New("A device may not have an empty name")
	}
	if inst.Connection != "" {
		if _, _, err := net.SplitHostPort(inst.Connection); err != nil {
			return fmt.Errorf("Cannot parse connection string in device config for %v - '%s': %v", id, inst.Connection, err)
		}
	}
	// We cannot validate plates or tipes at this point because the
	// inventory may not be loaded. So those get validated later on.
	return inst.LayoutPreferences.validate()
}

func (lo *LayoutOpt) validate() error {
	if lo == nil {
		return nil
	}
	return utils.ErrorSlice{
		lo.Tipboxes.validate("Tipboxes"),
		lo.Inputs.validate("Inputs"),
		lo.Outputs.validate("Outputs"),
		lo.Tipwastes.validate("Tipwastes"),
		lo.Wastes.validate("Wastes"),
		lo.Washes.validate("Washes"),
	}.Pack()
}

func (a Addresses) validate(layoutOptionName string) error {
	if len(a.Map()) != len(a) {
		return fmt.Errorf("Layout option field %s has duplicate addresses: %v", layoutOptionName, a)
	}
	return nil
}
