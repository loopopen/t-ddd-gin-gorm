package conf

import (
	"gorm.io/gorm/logger"
)

type MySQL struct {
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type ORM struct {
	AutoMigrate bool   `yaml:"auto_migrate"`
	MySQL       MySQL  `yaml:"mysql"`
	Level       string `yaml:"log_level"`
}

func (c ORM) GORMLogLevel() logger.LogLevel {
	switch c.Level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "off":
		return logger.Silent
	}
	return logger.Info
}
