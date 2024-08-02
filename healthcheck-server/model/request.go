package model

type HealthcheckRequest struct {
	Payload struct {
		AgentIP   string `json:"ip"`
		Timestamp int64  `json:"timestamp"`
		Duration  int    `json:"duration"`
	} `json:"payload"`
}

type HealthcheckResponse struct {
	Payload struct {
		StatusCode        int    `json:"status_code"`
		ReceivedTimestamp string `json:"received_timestamp"`
		ResponseTimestamp string `json:"response_timestamp"`
		Message           string `json:"message"`
	} `json:"payload"`
}
