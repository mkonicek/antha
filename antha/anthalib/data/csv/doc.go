/*
Package csv provides tools for data tables seriaization to/from CSV files.

The following CSV files format is currently supported:
	"column1name,column1type","column2name,column2type", ... ,"columnNname,columnNtype"
	value1_1,value1_2, ... , value1_N
	value2_1,value2_2, ... , value2_N

The following column types are currently supported:
	bool
	int64
	float64
	string
	TimestampMillis
	TimestampMicros

*/
package csv
