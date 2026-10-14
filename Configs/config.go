package Configs

import (
	"time"
)

type DB struct {
	Path string `json:"path"`
}

type Server struct {
	Port string `json:"port"`
}

type Task struct {
	Info       string `json:"info"`
	ReceivedAt time.Duration
	HandledAt  time.Time
	Handle     bool
}
