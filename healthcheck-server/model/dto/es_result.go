package dto

type uptime struct {
	Value float64 `json:"value"`
}
type ServerUptime struct {
	ID     uint   `json:"key"`
	Uptime uptime `json:"uptime_avg"`
}

type Response struct {
	Aggregtions struct {
		Server struct {
			Buckets []ServerUptime `json:"buckets"`
		} `json:"by_server"`
	} `json:"aggregations"`
}
