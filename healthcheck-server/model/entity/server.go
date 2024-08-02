package entity

import "time"

type Server struct {
	ID          uint   `gorm:"primary_key" gorm:"autoIncrement"`
	Name        string `gorm:"column:server_name" gorm:"unique"`
	IPv4        string
	Status      int
	CreatedTime time.Time
	LastUpdated time.Time
}

type ServerDoc struct {
	Server    Server `json:"server"`
	Timestamp int64  `json:"timestamp"`
	Duration  int    `json:"duration"`
}
