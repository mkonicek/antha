package liquidhandling

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-test/deep"
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// AssertLayoutsEquivalent checks that the layouts match, ignoring differences in object IDs
func AssertLayoutsEquivalent(got, expected []byte) error {
	var g, e layoutSummary
	if err := json.Unmarshal(got, &g); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'got'")
	} else if err := json.Unmarshal(expected, &e); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'expected'")
	}

	e.NormalizeIDs()
	g.NormalizeIDs()

	if diffs := deep.Equal(e, g); len(diffs) > 0 {
		return errors.Errorf("generated layout summary doesn't match expected: \n%s", strings.Join(diffs, "\n"))

	}
	return nil
}

// NormalizeIDs change the IDs of the contained objects to a standard form
func (ls *layoutSummary) NormalizeIDs() {
	beforeMap := make(map[string]string, len(ls.Before.Positions))
	for posName, pos := range ls.Before.Positions {
		if pos.Item != nil {
			newID := fmt.Sprintf("item_at_%s_before", posName)
			beforeMap[pos.Item.ID] = newID
			pos.Item.ID = newID
			pos.Item.Name = fmt.Sprintf("%s at %s", pos.Item.Kind, posName)
		}
	}

	afterMap := make(map[string]string, len(ls.After.Positions))
	for posName, pos := range ls.After.Positions {
		if pos.Item != nil {
			newID := fmt.Sprintf("item_at_%s_after", posName)
			beforeMap[pos.Item.ID] = newID
			pos.Item.ID = newID
			pos.Item.Name = fmt.Sprintf("%s at %s", pos.Item.Kind, posName)
		}
	}

	newIDMap := make(map[string]string, len(ls.IDMap))
	for before, after := range ls.IDMap {
		newIDMap[beforeMap[before]] = afterMap[after]
	}
	ls.IDMap = newIDMap
}

func (ls *layoutSummary) UnmarshalJSON(bs []byte) error {
	type LayoutSummaryAlias layoutSummary
	var a struct {
		LayoutSummaryAlias
		Version string `json:"version"`
	}
	if err := json.Unmarshal(bs, &a); err != nil {
		return errors.WithMessage(err, "unmarshalling layout")
	}

	if a.Version != LayoutSummaryVersion {
		return errors.Errorf("layout version mismatch: expected %s, got %s", LayoutSummaryVersion, a.Version)
	}

	*ls = layoutSummary(a.LayoutSummaryAlias)
	return nil
}

// AssertActionsEquivalent checks that the actions match, ignoring differences in object IDs
func AssertActionsEquivalent(got, expected []byte) error {
	var g, e actionsSummary
	if err := json.Unmarshal(got, &g); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'got'")
	} else if err := json.Unmarshal(expected, &e); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'expected'")
	}

	e.NormalizeIDs()
	g.NormalizeIDs()

	if diffs := deep.Equal(e, g); len(diffs) > 0 {
		return errors.Errorf("generated action summary differs from expected: \n%s", strings.Join(diffs, "\n"))
	}
	return nil
}

func (as actionsSummary) NormalizeIDs() {

	idUpdates := make(map[string]string)
	updateID := func(oldID string) string {
		if newID, ok := idUpdates[oldID]; ok {
			return newID
		} else {
			newID := fmt.Sprintf("ID_%d", len(idUpdates))
			idUpdates[oldID] = newID
			return newID
		}
	}

	// look through the transfers and update the destination as seen
	for _, action := range as {
		if tAction, ok := action.(*transferAction); ok {
			for _, tc := range tAction.Children {
				switch tChild := tc.(type) {
				case *parallelTransfer:
					for _, tSummary := range tChild.Channels {
						tSummary.From.Location.DeckItemID = updateID(tSummary.From.Location.DeckItemID)
						for _, update := range tSummary.To {
							update.Location.DeckItemID = updateID(update.Location.DeckItemID)
						}
					}
				case *tipAction:
					for _, channel := range tChild.Channels {
						channel.DeckItemID = updateID(channel.DeckItemID)
					}
				}
			}
		}
	}
}

func (as *actionsSummary) UnmarshalJSON(bs []byte) error {
	type partialAction struct {
		Kind string `json:"kind"`
	}

	var actions struct {
		Actions []*json.RawMessage `json:"actions"`
		Version string             `json:"version"`
	}
	if err := json.Unmarshal(bs, &actions); err != nil {
		return err
	}

	if actions.Version != ActionsSummaryVersion {
		return errors.Errorf("actions version mismatch: expected %s, got %s", ActionsSummaryVersion, actions.Version)
	}

	newActions := make(actionsSummary, 0, len(actions.Actions))
	for _, rawAction := range actions.Actions {
		var pa partialAction
		if err := json.Unmarshal(*rawAction, &pa); err != nil {
			return err
		}

		var a action
		switch pa.Kind {
		case "prompt":
			a = &promptAction{}
		case "transfer":
			a = &transferAction{}
		default:
			panic(fmt.Sprintf("unknown action type '%s'", pa.Kind))
		}
		if err := json.Unmarshal(*rawAction, a); err != nil {
			return err
		}
		newActions = append(newActions, a)
	}

	*as = newActions
	return nil
}

func (ta *transferAction) UnmarshalJSON(bs []byte) error {
	type partialChild struct {
		Kind string `json:"kind"`
	}

	type partialTA struct {
		Children []*json.RawMessage `json:"children"`
	}
	var pta partialTA
	if err := json.Unmarshal(bs, &pta); err != nil {
		return err
	}

	newChildren := make([]transferChild, 0, len(pta.Children))
	for _, rawChild := range pta.Children {
		var pc partialChild
		if err := json.Unmarshal(*rawChild, &pc); err != nil {
			return err
		}

		var tc transferChild
		switch pc.Kind {
		case string(loadTipAction):
			tc = &tipAction{}
		case string(unloadTipAction):
			tc = &tipAction{}
		case "parallel_transfer":
			tc = &parallelTransfer{}
		default:
			panic(fmt.Sprintf("unknown child kind '%s'", pc.Kind))
		}
		if err := json.Unmarshal(*rawChild, tc); err != nil {
			return err
		}
		newChildren = append(newChildren, tc)
	}

	ta.Children = newChildren
	return nil
}

func (h height) UnmarshalJSON(bs []byte) error {

	type partialHeight struct {
		measurementSummary
		Reference string `json:"reference"`
	}

	var ph partialHeight
	if err := json.Unmarshal(bs, &ph); err != nil {
		return err
	}

	h.measurementSummary = ph.measurementSummary

	if ref, err := wtype.NewWellReference(ph.Reference); err != nil {
		return err
	} else {
		h.Reference = ref
	}

	return nil
}
