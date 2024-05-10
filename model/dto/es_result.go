package dto

type uptime struct {
	Value float64 `json:"value"`
}
type ServerUptime struct {
	ID     uint   `json:"key"`
	Uptime uptime `json:"uptime"`
}

type Response struct {
	Aggregtions struct {
		Server struct {
			Buckets []ServerUptime `json:"buckets"`
		} `json:"server"`
	} `json:"aggregations"`
}
