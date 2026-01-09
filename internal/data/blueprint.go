package data

import "singularity/internal/enum"

type Blueprint struct {
	Id          string             `json:"id" bson:"id" toml:"id"`
	Name        string             `json:"name" bson:"name" toml:"name"`
	Type        enum.BlueprintType `json:"blueprint_type" bson:"blueprint_type" toml:"blueprint_type"`
	Volumes     []Volume           `json:"volumes" bson:"volumes" toml:"volumes"`
	Environment map[string]string  `json:"environment" bson:"environment" toml:"environment"`
}

type Volume struct {
	Id           string `toml:"id"`
	TargetFolder string `toml:"target_folder"`
}
