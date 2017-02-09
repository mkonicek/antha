package org_antha_lang_antha_v1

func (a *Measurement) MeasurementUnit() string {
	return a.Unit
}

func (a *Measurement) Quantity() float64 {
	return a.Value
}
