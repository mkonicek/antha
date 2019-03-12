package data

import (
	"github.com/apache/arrow/go/arrow"
)

// TODO: determine which time/timestamp data types we need (or we even don't need them at all and can store them in int64?)

// TimestampMillis is timestamp measured in ms
type TimestampMillis arrow.Timestamp

// TimestampMicros is timestamp measured in us
type TimestampMicros arrow.Timestamp
