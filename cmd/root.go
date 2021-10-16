package main

import (
	"github.com/mehditeymorian/etefagh/internal/cmd/server"
	"github.com/mehditeymorian/etefagh/internal/config"
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
)

func main() {
	var config config.Config = config.Config{
		Api: handler.Config{
			Port: "3000",
		},
		DB: db.Config{
			Uri:  "mongodb://localhost:27017",
			Name: "events",
		},
	}
	server.Main(config)
}
