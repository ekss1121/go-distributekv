package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Partition struct {
	Name  string
	Index int
	Host  string
}

type Partitions struct {
	Partitions []Partition
}

func ParseCofig(configPath string) Partitions {
	// read the partition config
	var partitions Partitions
	if _, err := toml.DecodeFile(configPath, &partitions); err != nil {
		log.Fatalf("toml.Decode(%s): %v", configPath, err)
	}

	return partitions
}
