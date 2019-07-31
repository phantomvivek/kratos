package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kratos/models"
)

//Config holds the config for this load test
var Config models.Configuration

func init() {
	Config = models.Configuration{}
}

//SetConfig setting the config on startup
func SetConfig() {

	//Set configuration
	configPath := ""

	configs := os.Args[1:]
	for _, config := range configs {
		vals := strings.Split(config, "=")
		if len(vals) > 1 {
			if vals[0] == "--config" {
				configPath = vals[1]
			}
		}
	}

	if configPath != "" {
		//Read file
		jsonConfig, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(jsonConfig), &Config)
		if err != nil {
			panic(err)
		}
	}
}
