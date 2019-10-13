package panel

type Panel interface {
	Update(uint64)
	IsEnd() bool
	Refresh()
}
