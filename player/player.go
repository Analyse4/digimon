package player

type Player struct {
	PlayerId int64
	NickName string
}

func New() (*Player, error) {
	p := new(Player)
	p.NickName = "Joker"
	p.PlayerId = 4
	//TODO: add player manager
	return p, nil
}
