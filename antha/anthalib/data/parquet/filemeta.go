package parquet

import (
	"github.com/xitongsys/parquet-go/parquet"
)

// FileKeyValueMetadata represents keys in file-level Parquet
// metadata, as defined in:
// https://github.com/apache/parquet-format/blob/master/src/main/thrift/parquet.thrift#L924 .
// In the presence of duplicate keys ion the file, behavior is undefined.
type FileKeyValueMetadata map[string]string

// Write gives a WriteOpt that writes all the given keyvalues into file metadata.
func (m FileKeyValueMetadata) Write() WriteOpt {
	return func(w *writeState) error {

		for k, v := range m {
			kv := parquet.NewKeyValue()
			kv.Key = k
			// we need a copy, not &v
			metaValue := v
			kv.Value = &metaValue
			w.writer.Footer.KeyValueMetadata = append(w.writer.Footer.KeyValueMetadata, kv)
		}

		return nil
	}
}

// Read gives a ReadOpt that populates this map as a side effect, when the file is read.
func (m FileKeyValueMetadata) Read() ReadOpt {
	return func(r *readState) error {
		kv := r.reader.Footer.GetKeyValueMetadata()
		for _, v := range kv {
			m[v.GetKey()] = v.GetValue()
		}

		return nil
	}
}
