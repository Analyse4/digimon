package panel

type Panel interface {
	Update(uint64, int32)
	IsEnd() bool
	Refresh()
	GetRPC(id uint64) int32
	RPCCompute() (uint64, error)
}
