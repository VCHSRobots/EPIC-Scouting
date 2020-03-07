/*
Package config provides tools for managing the configuration file.
*/
package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var defaultConfigString string = `DatabasePath: /home/USERNAME/database/
LogLocation: /home/USERNAME/log/`

var defaultConfig []byte = []byte(defaultConfigString)

func init() {
	_, err := os.Stat("config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			err := ioutil.WriteFile("config.yaml", defaultConfig, 0755)
			if err != nil {
				log.Fatal("Unable to create configuration file: " + err.Error())
			} else {
				log.Fatal("A new configuration file was created. Enter required values.")
			}
		} else {
			log.Fatal("Unable to read configuration file: " + err.Error())
		}
	}
}

// YAML describes the structure of the YAML configuration file. See https://godoc.org/gopkg.in/yaml.v3#Marshal for more information.
type YAML struct {
	DatabaseBackupFrequency string `yaml:"DatabaseBackupFrequency"`
	DatabaseBackupPath      string `yaml:"DatabaseBackupPath"`
	DatabasePath            string `yaml:"DatabasePath"`
	LogPath                 string `yaml:"LogPath"`
	Port                    int    `yaml:"Port"`
	TBAAuthKey              string `yaml:"TBAAuthKey"`
	Verbosity               int    `yaml:"Verbosity"`
}

/*
BuildInformation returns the build name and date.
*/
func BuildInformation() (buildName, buildDate string) {
	buildName = "Prerelease 0.1" // TODO: Update this with each release. "Prerelease" for development versions, "Release" for production versions.
	buildDate = "2020.064.2"     // TODO: Update this with each release. External script updates number and it is sourced from an external file. Format: YYYY.DDD.N
	return
}

/*
Load configuration values from "config.yaml".
*/
func Load() YAML {
	configuration, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Unable to read configuration file: " + err.Error())
	}

	var config YAML
	err2 := yaml.Unmarshal(configuration, &config)
	if err2 != nil {
		log.Fatal("Unable to parse configuration file as YAML: " + err.Error())
	}

	// Parse configuration settings.

	// Check mandatory settings.
	if config.DatabasePath == "" {
		log.Fatal("DatabasePath must be specified in configuration file.")
	}
	if config.LogPath == "" {
		log.Fatal("LogPath must be specified in configuration file.")
	}

	return config
}
