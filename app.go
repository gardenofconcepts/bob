package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func (app *App) configure() {
	app.configureLog()

	if len(app.Config) > 0 {
		config := ConfigReader()
		config.Read(app.Config)
		config.Apply(app)
	} else if _, err := os.Stat(CONFIG_FILE); os.IsExist(err) {
		config := ConfigReader()
		config.Read(CONFIG_FILE)
		config.Apply(app)
	}

	app.configureLog()
}

func (app *App) configureLog() {
	level := log.WarnLevel

	if app.Verbose {
		level = log.InfoLevel
	}

	if app.Debug {
		level = log.DebugLevel
	}

	log.SetLevel(level)
}
