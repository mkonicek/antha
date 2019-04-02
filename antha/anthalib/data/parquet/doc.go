/*
Package parquet provides tools for data tables serialization to and from Parquet files - in the form of files on disk, memory buffer or io.Reader/io.Writer.
Now read methods work pretty slowly with files having hundreds of columns.
As a workaround for now, Read methods support specifying a subset of columns to read.

Current implementation

Serialization to/from Parquet currently uses https://github.com/xitongsys/parquet-go library which supports serialization
to/from a slice of structs only. Therefore currently data tables are converted to (reflectively created) structs, which is pretty slow.

Future development

To make serialization faster, it would be beneficial to write Parquet files directly without using reflection. If some day
a Go library for serializing Arrow to/from Parquet is written, we should use it instead of the current implementation.

*/
package parquet
