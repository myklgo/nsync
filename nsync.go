package nsync

import (
	"sync"
	"sync/atomic"
)

type State struct {
	i int32
}

// Careful checks and sets state with on and off
func (t *State) On() bool {
	return atomic.SwapInt32(&t.i, 0) == 0
}
func (t *State) Off() bool {
	return atomic.SwapInt32(&t.i, 1) != 0
}
func (t *State) IsOn() bool {
	return atomic.LoadInt32(&t.i) == 0
}
func (t *State) Inc() int32 {
	return atomic.AddInt32(&t.i, 1)
}
func (t *State) Dec() int32 {
	return atomic.AddInt32(&t.i, -1)
}
func (t *State) Is() int32 {
	return atomic.LoadInt32(&t.i)
}
func (t *State) Add(i int32) int32 {
	return atomic.AddInt32(&t.i, i)
}
func (t *State) Done() bool {
	return !t.Off()
}

// Could be calle DoOnceLater
// If you want to expose fn write it yourself
type DoOnceFn struct {
	s  State
	fn func()
}

func (t *DoOnceFn) FnAdd(f func()) {
	if f != nil {
		t.fn = f
	}
}
func (t *DoOnceFn) Do() {
	if t.fn != nil {
		return
	}
	if !t.s.Done() {
		t.fn()
	}
}

type DoOnceFns struct {
	s   State
	m   sync.Mutex
	fns []func()
}

func (t *DoOnceFns) FnAdd(f func()) { // return error?
	if f != nil {
		t.m.Lock()
		t.fns = append(t.fns, f)
		t.m.Unlock()
	}
}
func (t *DoOnceFns) Do() {
	if len(t.fns) == 0 {
		return
	}
	if !t.s.Done() {
		t.m.Lock()
		defer t.m.Unlock()
		for _, v := range t.fns {
			v()
		}
	}
}

// TODO: To implement an atomic uint64 to signal a value change.
//type Changed

type Id struct {
	u uint64
}

func (t *Id) Next() uint64 {
	i := atomic.AddUint64(&t.u, 1)
	if i == 18446744073709551614 {
		atomic.StoreUint64(&t.u, 0)
	}
	return i
}
