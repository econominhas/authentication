package ulid

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

type Ulid struct{}

var entropyPool = sync.Pool{
	New: func() any {
		entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
		return entropy
	},
}

func (adp *Ulid) GenId() (string, error) {
	e := entropyPool.Get().(*ulid.MonotonicEntropy)
	s := ulid.MustNew(ulid.Timestamp(time.Now()), e).String()
	entropyPool.Put(e)

	return s, nil
}

func NewUlid() *Ulid {
	return &Ulid{}
}
