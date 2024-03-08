package types

import "sync"

const (
	EVENT string = "EVENT"
	REQ   string = "REQ"
	CLOSE string = "CLOSE"
)

// ChanGroup is a neat combination of WaitGroup and channel
// Usage pattern:
// cg := NewChanGroup()
// cg.Add(1) // inside constructor or before goroutine
// do work in goroutine using chan, always call defer cg.Done()
// go func() { defer cg.Done() ... }()
// go func() { cg.WaitClose() }() // close chan after work is .Done()
// for item := range cg.Chan { } // process results in calling func
type ChanGroup struct {
	sync.WaitGroup
	Chan chan *EnvelopeWrapper
}

func NewChanGroup() *ChanGroup {
	cg := ChanGroup{
		Chan: make(chan *EnvelopeWrapper),
	}
	cg.Add(1)
	return &cg
}

func (cg *ChanGroup) WaitClose() {
	cg.Wait()
	close(cg.Chan)
}
