package model

import "time"

type WebDavSession struct {
	Username    string    `json:"username"`
	ClientIP    string    `json:"client_ip"`
	OuterIP     string    `json:"outer_ip"`
	UserAgent   string    `json:"user_agent"`
	ConnectTime time.Time `json:"connect_time"`
	SessionID   string    `json:"session_id"`
}
