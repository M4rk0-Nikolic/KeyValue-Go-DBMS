package config

import (
	"os"
	"log"
	"gopkg.in/yaml.v3"
)

type Config struct {
    BlockManager BlockManagerConfig `yaml:"BlockManager"`
	WAL WALConfig `yaml:"WAL"`
}

type BlockManagerConfig struct {
    BlockSizeInKB uint64 `yaml:"BlockSizeInKB"`
	BlockCacheSize uint64 `yaml:"BlockCacheSize"`
}

type WALConfig struct {
    LengthOfSegInBlocks uint64 `yaml:"LengthOfSegInBlocks"`
}
func LoadConfig(){
	file, err := os.ReadFile("config.yaml")
	var Global Config

	if err != nil{
		log.Println("Failed to open and read config.yaml - using defaults instead.")
	}

	err2 := yaml.Unmarshal(file, &Global)
	if err2 != nil{
		log.Println("Failed to open and read config.yaml - using defaults instead.")
	}
}