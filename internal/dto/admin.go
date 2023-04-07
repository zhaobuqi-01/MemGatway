package dto

import "time"

type AminInfoOutput struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	LoginTime     time.Time `json:"login_time"`
	Avatar        string    `json:"avatar"`
	Introduceions string    `json:"introduceion"`
	Roles         []string  `json:"roles"`
}
