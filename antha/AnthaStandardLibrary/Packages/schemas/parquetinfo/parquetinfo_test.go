package parquetinfo

import (
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/schemas"
	"github.com/antha-lang/antha/antha/anthalib/data/parquet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAnthaInfo(t *testing.T) {
	meta := (parquet.FileKeyValueMetadata)(nil)
	require.NoError(t, AddAnthaInfo(meta, schemas.DataInfo{}))
	assert.Nil(t, meta, "nil")
	meta = parquet.FileKeyValueMetadata{"x": "1"}
	require.NoError(t, AddAnthaInfo(meta, schemas.DataInfo{}))
	assert.Contains(t, meta, "x")
	assert.Contains(t, meta, ParquetInfoKey)
}

func TestRoundtripAnthaInfo(t *testing.T) {
	meta := parquet.FileKeyValueMetadata{}
	inf := schemas.DataInfo{Tool: &schemas.Tool{ExternalID: "a tool"}}
	require.NoError(t, AddAnthaInfo(meta, inf))
	infUnmarshal, err := GetAnthaInfo(meta)
	require.NoError(t, err)
	assert.EqualValues(t, inf, infUnmarshal)
}
