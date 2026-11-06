package main

import (
	"time"
)

type Move struct {
	Id          int64 //could be string if UUID
	Name        string
	DanceStyle string
	Difficulty  int
	Description string
	Image       string
	CreatedAt time.Time
}