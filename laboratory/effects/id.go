package effects

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type IDGenerator struct {
	prefix  string
	counter uint64
}

func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{
		prefix: prefix,
	}
}

func (idg *IDGenerator) NextID() string {
	num := atomic.AddUint64(&idg.counter, 1)
	return fmt.Sprintf("%s-%08d", idg.prefix, num)
}

func (idg *IDGenerator) MarshalJSON() ([]byte, error) {
	return json.Marshal(&idGeneratorJson{
		Prefix:  idg.prefix,
		Counter: atomic.LoadUint64(&idg.counter),
	})
}

func (idg *IDGenerator) UnmarshalJSON(bs []byte) error {
	idgj := &idGeneratorJson{}
	if err := json.UnmarshalJSON(bs, idgj); err != nil {
		return err
	} else {
		idg.prefix = idgj.Prefix
		atomic.StoreUint64(&idg.counter, idgj.Counter)
		return nil
	}
}

type idGeneratorJson struct {
	Prefix  string
	Counter uint64
}
