package parquet

import (
	"github.com/xitongsys/parquet-go/parquet"
)

// FileKeyValueMetadata represents keys in file-level Parquet
// metadata, as defined in:
// https://github.com/apache/parquet-format/blob/master/src/main/thrift/parquet.thrift#L924 .
// In the presence of duplicate keys in the file, behavior is undefined.
type FileKeyValueMetadata map[string]string

// Write gives a WriteOpt that writes all the given keyvalues into file metadata.
func (m FileKeyValueMetadata) Write() WriteOpt {
	return func(w *writeState) error {

		for k, v := range m {
			// we need a copy, not &v
			metaValue := v
			kv := parquet.NewKeyValue()
			kv.Key = k
			kv.Value = &metaValue
			w.writer.Footer.KeyValueMetadata = append(w.writer.Footer.KeyValueMetadata, kv)
		}

		return nil
	}
}

// Read gives a ReadOpt that populates this map as a side effect, when the file is read.
func (m FileKeyValueMetadata) Read() ReadOpt {
	return func(r *readState) error {
		kvs := r.reader.Footer.GetKeyValueMetadata()
		for _, kv := range kvs {
			m[kv.GetKey()] = kv.GetValue()
		}

		return nil
	}
}
