package smpp

type DLRTracker struct {
	expected map[string]bool
	done     chan string
}

func NewDLRTracker() *DLRTracker {
	return &DLRTracker{
		expected: make(map[string]bool),
		done:     make(chan string, 10),
	}
}

func (t *DLRTracker) Expect(id string) {
	t.expected[id] = true
}

func (t *DLRTracker) Receive(dlr *DLR) {
	_, ok := t.expected[dlr.MessageID]
	if ok {
		t.done <- dlr.MessageID
	}
}

func (t *DLRTracker) Done() <-chan string {
	return t.done
}
