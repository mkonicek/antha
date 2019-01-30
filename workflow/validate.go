package workflow

import (
	"errors"
	"fmt"
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
			wf.Config.validate(wf),
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

func (cfg Config) validate(wf *Workflow) error {
	return utils.ErrorSlice{
		cfg.GilsonPipetMax.validate(wf),
		cfg.GlobalMixer.validate(wf),
	}.Pack()
}

func (gilson GilsonPipetMaxConfig) validate(wf *Workflow) error {
	if err := gilson.Defaults.validate("Defaults", wf); err != nil {
		return err
	}
	for name, cfg := range gilson.Devices {
		if err := cfg.validate(string(name), wf); err != nil {
			return err
		}
	}
	return nil
}

func (gilson *GilsonPipetMaxInstanceConfig) validate(name string, wf *Workflow) error {
	switch {
	case gilson == nil: // should only happen when name == Defaults
		return nil
	case gilson.MaxPlates != nil && *gilson.MaxPlates <= 0:
		return fmt.Errorf("Validation error: GilsonPipetMax '%s': MaxPlates must be > 0", name)
	case gilson.MaxWells != nil && *gilson.MaxWells <= 0:
		return fmt.Errorf("Validation error: GilsonPipetMax '%s': MaxWells must be > 0", name)
	case gilson.ResidualVolumeWeight != nil && *gilson.ResidualVolumeWeight < 0:
		return fmt.Errorf("Validation error: GilsonPipetMax '%s': ResidualVolumeWeight must be >= 0", name)
	}

	// We cannot validate plates at this point because the real plate
	// type inventory may not be loaded. So that gets validated later
	// on.
	return nil
}

func (global GlobalMixerConfig) validate(wf *Workflow) error {
	// Again, we cannot validate plates and plate types until we have a
	// working inventory system.
	return nil
}
