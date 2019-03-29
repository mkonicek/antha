/*
Package parquet provides tools for data tables seriaization to/from Parquet files.

Current implementation

Serialization to/from Parquet currently uses https://github.com/xitongsys/parquet-go library which supports serialization
to/from a slice of structs only. Therefore currently data tables are converted to (reflectively created) structs, which is pretty slow.

Future development

To make serialization faster, it would be beneficial to write Parquet files directly without using reflection. If some day
a Go library for serializing Arrow to/from Parquet is written, we should use it instead of the current implementation.

*/
package parquet
