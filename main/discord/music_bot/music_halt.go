package musicbot

import (
	"sync/atomic"
)

var halted atomic.Bool

func haltBot() {
	halted.Store(true)
}

func resetHalt() {
	halted.Store(false)
}

func isHalted() bool {
	return halted.Load()
}
