package models

import "time"

type Notice struct {
	ID        string
	Title     string
	Content   string
	URL       string
	PostDate  time.Time
	CreatedAt time.Time
}
