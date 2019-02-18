package workflow

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
			wf.Elements.validate(wf),
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

func (es Elements) validate(wf *Workflow) error {
	return utils.ErrorSlice{
		es.Types.validate(wf),
		es.Instances.validate(wf),
		es.InstancesParameters.validate(wf),
		es.InstancesConnections.validate(wf),
	}.Pack()
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
		if _, found := wf.Elements.Instances[name]; !found {
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
	if _, found := wf.Elements.Instances[soc.ElementInstance]; !found {
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
		cfg.Tecan.validate(),
		cfg.CyBio.validate(),
		cfg.Labcyte.validate(),
		cfg.assertOnlyOneMixer(),
	}.Pack()
}

func (cfg Config) assertOnlyOneMixer() error {
	// remove / revise when we get better. NB: because of this test, we
	// don't need to check that all devices have unique IDs. Once we
	// relax this, we may need to do that.
	if count := len(cfg.GilsonPipetMax.Devices) + len(cfg.Tecan.Devices) + len(cfg.CyBio.Devices) + len(cfg.Labcyte.Devices); count > 1 {
		return fmt.Errorf("Currently a maximum of one mixer can be used per workflow. You have %d configured.", count)
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

// Gilson
func (gilsons GilsonPipetMaxConfig) validate() error {
	if err := gilsons.Defaults.validate("Defaults", true); err != nil {
		return err
	}
	for id, inst := range gilsons.Devices {
		if err := inst.validate(id, false); err != nil {
			return err
		}
	}
	return nil
}

func (inst *GilsonPipetMaxInstanceConfig) validate(id DeviceInstanceID, isDefaults bool) error {
	if len(id) == 0 {
		return errors.New("GilsonPipetMax: A device may not have an empty name.")

	} else if inst == nil {
		if isDefaults {
			return nil
		} else {
			return fmt.Errorf("GilsonPipetMax device '%s' has no configuration!", id)
		}

	} else if !isDefaults && strings.ToLower(string(id)) == "defaults" {
		return fmt.Errorf("Confusion: GilsonPipetMax device '%s' exists. Did you mean to set GilsonPipetMax.Defaults instead?")

	}
	return inst.commonMixerInstanceConfig.validate(id)
}

// Tecan
func (tecans TecanConfig) validate() error {
	if err := tecans.Defaults.validate("Defaults", true); err != nil {
		return err
	}
	for id, inst := range tecans.Devices {
		if err := inst.validate(id, false); err != nil {
			return err
		}
	}
	return nil
}

func (inst *TecanInstanceConfig) validate(id DeviceInstanceID, isDefaults bool) error {
	if len(id) == 0 {
		return errors.New("Tecan: A device may not have an empty name.")

	} else if inst == nil {
		if isDefaults {
			return nil
		} else {
			return fmt.Errorf("Tecan device '%s' has no configuration!", id)
		}

	} else if !isDefaults && strings.ToLower(string(id)) == "defaults" {
		return fmt.Errorf("Confusion: Tecan device '%s' exists. Did you mean to set Tecan.Defaults instead?")

	}
	return inst.commonMixerInstanceConfig.validate(id)
}

// CyBio
func (cybios CyBioConfig) validate() error {
	if err := cybios.Defaults.validate("Defaults", true); err != nil {
		return err
	}
	for id, inst := range cybios.Devices {
		if err := inst.validate(id, false); err != nil {
			return err
		}
	}
	return nil
}

func (inst *CyBioInstanceConfig) validate(id DeviceInstanceID, isDefaults bool) error {
	if len(id) == 0 {
		return errors.New("CyBio: A device may not have an empty name.")

	} else if inst == nil {
		if isDefaults {
			return nil
		} else {
			return fmt.Errorf("CyBio device '%s' has no configuration!", id)
		}

	} else if !isDefaults && strings.ToLower(string(id)) == "defaults" {
		return fmt.Errorf("Confusion: CyBio device '%s' exists. Did you mean to set CyBio.Defaults instead?")

	}
	return inst.commonMixerInstanceConfig.validate(id)
}

// Labcyte
func (labcytes LabcyteConfig) validate() error {
	if err := labcytes.Defaults.validate("Defaults", true); err != nil {
		return err
	}
	for id, inst := range labcytes.Devices {
		if err := inst.validate(id, false); err != nil {
			return err
		}
	}
	return nil
}

func (inst *LabcyteInstanceConfig) validate(id DeviceInstanceID, isDefaults bool) error {
	if len(id) == 0 {
		return errors.New("Labcyte: A device may not have an empty name.")

	} else if inst == nil {
		if isDefaults {
			return nil
		} else {
			return fmt.Errorf("Labcyte device '%s' has no configuration!", id)
		}

	} else if !isDefaults && strings.ToLower(string(id)) == "defaults" {
		return fmt.Errorf("Confusion: Labcyte device '%s' exists. Did you mean to set Labcyte.Defaults instead?")

	}
	// NB because the instruction plugin itself does validation of the model, we don't do that here!
	return inst.commonMixerInstanceConfig.validate(id)
}

func (inst *commonMixerInstanceConfig) validate(id DeviceInstanceID) error {
	if inst.ExecFile != "" {
		if abs, err := exec.LookPath(inst.ExecFile); err != nil {
			return fmt.Errorf("Error when trying to locate executable at %v for %v: %v", inst.ExecFile, id, err)
		} else {
			inst.ExecFile = abs
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
