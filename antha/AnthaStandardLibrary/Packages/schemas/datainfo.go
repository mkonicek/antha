package schemas

import "encoding/json"

// DataInfo contains metadata that cannot be inferred from a file directly.
// All fields should be considered optional.
type DataInfo struct {
	Tool    *Tool    `json:"tool,omitempty"`
	Origin  *Origin  `json:"origin,omitempty"`
	Columns *Columns `json:"columns,omitempty"`

	// A short name identifying the type of data represented - such as
	// "spectrometer" or "cell count". This is not human-readable and may be
	// used in future as as key in a schema registry or ontology.
	Kind string `json:"kind,omitempty"`

	// A longer comment on the provenance of the data, which might be useful to
	// end-users.
	Comment string `json:"comment,omitempty"`
}

// Origin describes the original source of this data, if any.
type Origin struct {
	// If the data was transformed from a vendor-specific file format, this
	// can be used to identify that format (eg. "xml" or "csv")
	Format string `json:"format,omitempty"`

	// Arbitrary key-value properties from a vendor-specific file.
	Properties map[string]json.RawMessage `json:"properties,omitempty"`
}

// Tool describes the software tool that generated this data.
type Tool struct {
	// An artifact version, in machine-readable form such as a docker image
	// name.
	ExternalID string `json:"external_id,omitempty"`
}

// Columns holds column-level metadata.
type Columns struct {
	// If set, this should contain human-readable column names in the same order
	// as the physical dataset. This may be useful if the file format does not
	// allow arbitrary names in the physical schema.
	Names []string `json:"names,omitempty"`
}

// TODO further usecases
// - data sort order (NB. duplicating Parquet feature?)
// - column statistics
// - Go datatype information
// - biological ontology
// - units of measure
// - timezone of potentially ambiguous timestamp values
// - time source used (eg based on what device's clock)
// - operator ID
// - legal - copyright info, information barriers
