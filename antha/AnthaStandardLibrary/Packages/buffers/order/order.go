// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// pacakge order deals with adding order details to an LHComponent.
package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Option is a type for use as an argument in the SetOrderInfo function.
type Option string

const (
	// Force overwriting of Order Details of a wtype.Liquid in the SetOrderInfo function.
	ForceUpdate Option = "FORCEUPDATE"
)

// Key to look up order details from a wtype.Liquid.
const OrderDetails = "ORDERDETAILS"

var (
	errNotFound = errors.New("no order info found")
)

// Details stores the order details of a wtype.Liquid
type Details struct {

	// Name of Manufacturer.
	Manufacturer string

	// The catalogue number of the item for ordering.
	LotID string

	// The specific batch number of this instance of the item.
	BatchID string

	// The expiry date
	ExpiryDate time.Time

	// Any restrictions for storage of the item.
	StorageRequirements StorageConditions
}

// StorageConditions stores any restrictions.
type StorageConditions struct {
	// Minimum Temperature recommended for storage.
	MinTemp wunit.Temperature
	// Maximum Temperature recommended for storage.
	MaxTemp wunit.Temperature
	// Is the Item sensitive to light.
	LightSensitive bool
	// Is the Item sensitive to moisture.
	MoistureSensistive bool
	// Is the Item sensitive to oxygen.
	OxygenSensistive bool
	// Is the Item sensitive to Freeze thaws.
	FreezeThawSensitive bool
}

// String returns a summary of any storage restrictions.
func (s StorageConditions) String() string {
	var sensitive []string
	var notSensitive []string

	names := map[string]bool{
		"light":       s.LightSensitive,
		"moisture":    s.MoistureSensistive,
		"oxygen":      s.OxygenSensistive,
		"freeze/thaw": s.FreezeThawSensitive,
	}

	for name, s := range names {
		if s {
			sensitive = append(sensitive, name)
		} else {
			notSensitive = append(notSensitive, name)
		}
	}

	sensitivities := ""
	if len(sensitive) > 0 {
		sensitivities = fmt.Sprintf(", sensitive to %s", strings.Join(sensitive, ", "))
	}
	insensitivities := ""
	if len(notSensitive) > 0 {
		insensitivities = fmt.Sprintf(", not sensitive to %s", strings.Join(notSensitive, ", "))
	}

	return fmt.Sprintf("Temperature Range: [%v - %v]%s%s.", s.MinTemp, s.MaxTemp, sensitivities, insensitivities)
}

// GetOrderDetails returns order Details for a component.
func GetOrderDetails(comp *wtype.Liquid) (orderDetails Details, err error) {

	order, found := comp.Extra[OrderDetails]

	if !found {
		return orderDetails, errNotFound
	}

	var bts []byte

	bts, err = json.Marshal(order)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &orderDetails)

	if err != nil {
		err = fmt.Errorf("problem getting %s order details: %s", comp.Name(), err.Error())
	}

	return
}

// SetOrderDetails adds order details to a wtype.Liquid.
// An error will be returned if order details are already encountered unless the ForceUpdate option is used as an Option argument in the function.
// In which case any existing order details will be overwritten.
func SetOrderDetails(comp *wtype.Liquid, orderDetails Details, options ...Option) (*wtype.Liquid, error) {

	// look for existing order details
	existingDetails, err := GetOrderDetails(comp)

	if err == errNotFound {
		comp.Extra[OrderDetails] = orderDetails
		return comp, nil
	} else if err == nil {
		if !inOptions(ForceUpdate, options) {
			return comp, fmt.Errorf("component %s already contains order info %v. Use order ForceUpdate to override.", comp.Name(), existingDetails)
		}
	} else {
		return comp, err
	}

	return comp, nil
}

func inOptions(query Option, options []Option) bool {
	for _, option := range options {
		if strings.EqualFold(string(query), string(option)) {
			return true
		}
	}
	return false
}
