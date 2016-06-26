package models

import "time"

type Lead struct {
	Id        int
	Email     string
	CreatedAt time.Time
}
