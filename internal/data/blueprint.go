package data

type Blueprint struct {
	Id          string            `json:"id" bson:"id" toml:"id"`
	Name        string            `json:"name" bson:"name" toml:"name"`
	Type        string            `json:"type" toml:"type"`
	Volumes     []Volume          `json:"volumes" bson:"volumes" toml:"volumes"`
	Environment map[string]string `json:"environment" bson:"environment" toml:"environment"`
}

type Volume struct {
	Id           string `json:"id" toml:"id"`
	TargetFolder string `json:"target_folder" toml:"target_folder"`
}
