package config

import (
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
	"github.com/mehditeymorian/etefagh/internal/logger"
)

// Config Struct of all the configuration in the app
type Config struct {
	Api    handler.Config
	DB     db.Config
	Logger logger.Config
}
