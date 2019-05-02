package parquetinfo

import (
	"encoding/json"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/schemas"
	"github.com/antha-lang/antha/antha/anthalib/data/parquet"
)

// ParquetInfoKey is used to embed JSON DataInfo in a Parquet file footer as
// `key_value_metadata`.
const ParquetInfoKey = "antha:data-info"

// AddAnthaInfo inserts the JSON DataInfo representation in the footer when
// writing a file. Use it like:
//      meta := parquet.FileKeyValueMetadata{}
//      _ = parquetinfo.AddAnthaInfo(meta, dataInfo)
//      parquet.TableToBytes(table, meta.Write())
func AddAnthaInfo(meta parquet.FileKeyValueMetadata, dataInfo schemas.DataInfo) error {
	if meta == nil {
		return nil
	}
	infoJSON, err := json.Marshal(dataInfo)
	if err != nil {
		return err
	}
	meta[ParquetInfoKey] = string(infoJSON)
	return err
}

// GetAnthaInfo reads the JSON DataInfo representation from the footer.
// Use it together with data tables like:
//      meta := parquet.FileKeyValueMetadata{}
//		table, err = parquet.TableFromReader(reader, meta.Read())
//		dataInfo, err := GetAnthaInfo(meta)
func GetAnthaInfo(meta parquet.FileKeyValueMetadata) (schemas.DataInfo, error) {
	inf := schemas.DataInfo{}
	if meta == nil {
		return inf, nil
	}
	infoJSON := meta[ParquetInfoKey]
	if infoJSON == "" {
		return inf, nil
	}
	err := json.Unmarshal([]byte(infoJSON), &inf)
	return inf, err
}
