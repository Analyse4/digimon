package panel

type RulingCounterPanel struct {
	counter []uint64
}

func NewRulingCounterPanel(num int32) *RulingCounterPanel {
	return &RulingCounterPanel{counter: make([]uint64, num)}
}

func (rCP *RulingCounterPanel) Update(winID uint64) {
	for i, v := range rCP.counter {
		if v == 0 {
			rCP.counter[i] = winID
		}
	}
}

func (rCP *RulingCounterPanel) IsAttackerWin(attackerID uint64) bool {
	for _, v := range rCP.counter {
		if v != attackerID {
			return false
		}
	}
	return true
}

func (rCP *RulingCounterPanel) IsBattleEnd() bool {
	for _, v := range rCP.counter {
		if v == 0 {
			return false
		}
	}
	return true
}

func (rCP *RulingCounterPanel) Refresh() {
	rCP.counter = nil
}

func (rCP *RulingCounterPanel) GetNum() int {
	return len(rCP.counter)
}
