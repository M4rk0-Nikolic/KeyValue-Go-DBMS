package config

import (
	"os"
	"log"
	"gopkg.in/yaml.v3"
)


file, err := os.ReadFile("config.yaml")
var GenericMap map[string]any

if err != nil{
	log.Println("Failed to open and read config.yaml - using defaults instead.")
}

err := yaml.Unmarshal(yamlFile, &GenericMap)
if err != nil{
	log.Println("Failed to open and read config.yaml - using defaults instead.")
}
