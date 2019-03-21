package parquet

import (
	"bytes"
	"io"

	"github.com/xitongsys/parquet-go/ParquetFile"
)

// Implementations of in-memory ParquetFile.ParquetFile.
// They are needed because parquet-go library does not contain any valid tools for reading/writing Parquet in memory.

// a read-only ParquetFile.ParquetFile on the top of a memory buffer
type readOnlyMemoryParquetFile struct {
	data   []byte        // array of bytes to read from
	reader *bytes.Reader // bytes.Reader to read this array of bytes
}

func newReadOnlyMemoryParquetFile(data []byte) *readOnlyMemoryParquetFile {
	return &readOnlyMemoryParquetFile{
		data:   data,
		reader: bytes.NewReader(data),
	}
}

func (f *readOnlyMemoryParquetFile) Create(name string) (ParquetFile.ParquetFile, error) {
	panic("SHOULD NOT HAPPEN: readOnlyMemoryParquetFile.Create is not implemented")
}

func (f *readOnlyMemoryParquetFile) Open(name string) (ParquetFile.ParquetFile, error) {
	// parquet-go uses Open("") only - this weird interface means creating (and opening) a copy of the same file
	if name != "" {
		panic("SHOULD NOT HAPPEN: readOnlyMemoryParquetFile.Open(fileName) is not implemented")
	}
	return newReadOnlyMemoryParquetFile(f.data), nil
}

func (f *readOnlyMemoryParquetFile) Seek(offset int64, pos int) (int64, error) {
	return f.reader.Seek(offset, pos)
}

func (f *readOnlyMemoryParquetFile) Read(b []byte) (cnt int, err error) {
	return f.reader.Read(b)
}

func (f *readOnlyMemoryParquetFile) Write(b []byte) (n int, err error) {
	panic("SHOULD NOT HAPPEN: readOnlyMemoryParquetFile.Write is not implemented")
}

func (f *readOnlyMemoryParquetFile) Close() error {
	return nil
}

var _ ParquetFile.ParquetFile = (*readOnlyMemoryParquetFile)(nil)

// a write-only ParquetFile.ParquetFile on the top of io.Writer
type writeOnlyParquetFile struct {
	writer io.Writer
}

func newWriteOnlyParquetFile(writer io.Writer) *writeOnlyParquetFile {
	return &writeOnlyParquetFile{
		writer: writer,
	}
}

func (f *writeOnlyParquetFile) Create(name string) (ParquetFile.ParquetFile, error) {
	panic("SHOULD NOT HAPPEN: writeOnlyParquetFile.Create is not implemented")
}

func (f *writeOnlyParquetFile) Open(name string) (ParquetFile.ParquetFile, error) {
	panic("SHOULD NOT HAPPEN: writeOnlyParquetFile.Open is not implemented")
}

func (f *writeOnlyParquetFile) Seek(offset int64, pos int) (int64, error) {
	panic("SHOULD NOT HAPPEN: writeOnlyParquetFile.Seek is not implemented")
}

func (f *writeOnlyParquetFile) Read(b []byte) (cnt int, err error) {
	panic("SHOULD NOT HAPPEN: writeOnlyParquetFile.Read is not implemented")
}

func (f *writeOnlyParquetFile) Write(b []byte) (n int, err error) {
	return f.writer.Write(b)
}

func (f *writeOnlyParquetFile) Close() error {
	return nil
}

var _ ParquetFile.ParquetFile = (*writeOnlyParquetFile)(nil)
