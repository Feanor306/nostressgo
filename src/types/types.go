package types

import "sync"

const (
	EVENT string = "EVENT"
	REQ   string = "REQ"
	CLOSE string = "CLOSE"
)

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
