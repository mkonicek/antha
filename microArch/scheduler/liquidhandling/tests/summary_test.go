package tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
)

// AssertLayoutsEquivalent checks that the layouts match, ignoring differences in object IDs
func AssertLayoutsEquivalent(got, expected []byte) error {
	var g, e lh.LayoutSummary
	if err := json.Unmarshal(got, &g); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'got'")
	} else if err := json.Unmarshal(expected, &e); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'expected'")
	}

	NormalizeLayoutIDs(&e)
	NormalizeLayoutIDs(&g)

	if !reflect.DeepEqual(e, g) {
		// serialize to make tracking down the actual differeces easier
		if eJSON, err := json.Marshal(e); err != nil {
			panic(err)
		} else if gJSON, err := json.Marshal(g); err != nil {
			panic(err)
		} else {
			return errors.Errorf("generated layout differs from expected:\n\te: %s\n\tg: %s", string(eJSON), string(gJSON))
		}
	}
	return nil
}

// NormalizeIDs change the IDs of the contained objects to a standard form
func NormalizeLayoutIDs(ls *lh.LayoutSummary) {
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

// AssertActionsEquivalent checks that the actions match, ignoring differences in object IDs
func AssertActionsEquivalent(got, expected []byte) error {
	var g, e lh.ActionsSummary
	if err := json.Unmarshal(got, &g); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'got'")
	} else if err := json.Unmarshal(expected, &e); err != nil {
		return errors.WithMessage(err, "failed to unmarshal 'expected'")
	}

	NormalizeActionsIDs(e)
	NormalizeActionsIDs(g)

	if diffs := deep.Equal(e, g); len(diffs) > 0 {
		return errors.Errorf("generated action summary differs from expected: %s", strings.Join(diffs, "; "))
	}
	return nil
}

func NormalizeActionsIDs(as lh.ActionsSummary) {
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
		if tAction, ok := action.(*lh.TransferAction); ok {
			for _, tc := range tAction.Children {
				switch tChild := tc.(type) {
				case *lh.ParallelTransfer:
					for _, tSummary := range tChild.Channels {
						tSummary.From.Location.DeckItemID = updateID(tSummary.From.Location.DeckItemID)
						for _, update := range tSummary.To {
							update.Location.DeckItemID = updateID(update.Location.DeckItemID)
						}
					}
				case *lh.TipAction:
					for _, channel := range tChild.Channels {
						channel.DeckItemID = updateID(channel.DeckItemID)
					}
				}
			}
		}
	}
}
