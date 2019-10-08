package prometheus

import (
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var playerGauge stdprometheus.Gauge
var roomTotalGauge stdprometheus.Gauge
var inGameRoomGauge stdprometheus.Gauge
var handlerLatencySummary *stdprometheus.SummaryVec

func Init() {
	playerGauge = promauto.NewGauge(stdprometheus.GaugeOpts{
		Namespace: "player",
		Name:      "number",
		Help:      "total number of players in game server",
	})
	roomTotalGauge = promauto.NewGauge(stdprometheus.GaugeOpts{
		Namespace: "room",
		Name:      "number",
		Help:      "total rooms of rooms in game server",
	})
	inGameRoomGauge = promauto.NewGauge(stdprometheus.GaugeOpts{
		Namespace: "room",
		Subsystem: "in_game",
		Name:      "number",
		Help:      "total number of in-game players in game server",
	})
	//handlerLatencySummary = promauto.NewSummary(stdprometheus.SummaryOpts{
	//	Namespace:   "handler",
	//	Name:        "",
	//	Help:        "",
	//	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	//})
	handlerLatencySummary = promauto.NewSummaryVec(stdprometheus.SummaryOpts{
		Namespace:  "handler",
		Name:       "latency",
		Help:       "handler latency",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"router", "result"})
}

func GetPlayerGauge() stdprometheus.Gauge {
	return playerGauge
}

func GetRoomGauge() stdprometheus.Gauge {
	return roomTotalGauge
}

func GetInGameRoomGauge() stdprometheus.Gauge {
	return inGameRoomGauge
}

func GethandlerLatencySummary() *stdprometheus.SummaryVec {
	return handlerLatencySummary
}
