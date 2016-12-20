package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"fmt"
)

func ConfigReader() *App {
	return &App{}
}

func (config *App) Read(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithField("file", path).Fatal("Configuratoin file not found")
	}

	log.WithField("file", path).Info("Read configuration file")

	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.WithError(err).Fatal("Error while reading configuration file")
	}

	err = yaml.Unmarshal(data, config)

	fmt.Printf("Read configuration: %+v\n", config)

	if err != nil {
		log.WithError(err).Fatal("Error while parsing configuration file")
	}
}

func (config *App) Apply(app *App) error {
	if err := mergo.Map(app, config); err != nil {
		log.WithError(err).Fatal("Error while applying configuration")
	}

	if config.Download == false {
		app.Download = false
	}

	if config.Upload == false {
		app.Upload = false
	}

	if config.Verbose == true {
		app.Verbose = true
	}

	if config.Debug == true {
		app.Debug = true
	}

	if app.Pattern != CONFIG_PATTERN {
		app.Pattern = config.Pattern
	}

	/*if app.Include != []string{CONFIG_INCLUDE} {
		app.Include = config.Include
	}

	if app.Exclude != []string{CONFIG_EXCLUDE} {
		app.Exclude = config.Exclude
	}*/

	fmt.Printf("Try to apply: %+v\n", app)

	return nil
}
