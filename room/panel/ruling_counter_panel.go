package panel

type RulingCounterPanel struct {
	counter []bool
}

func NewRulingCounterPanel(num int32) *RulingCounterPanel {
	return &RulingCounterPanel{counter: make([]bool, num)}
}

func (rCP *RulingCounterPanel) Update() {
	for i, v := range rCP.counter {
		if v == false {
			rCP.counter[i] = true
		}
	}
}

func (rCP *RulingCounterPanel) IsEnd() bool {
	for _, v := range rCP.counter {
		if v == false {
			return false
		}
	}
	return true
}

func (rCP *RulingCounterPanel) Refresh() {
	rCP.counter = nil
}
