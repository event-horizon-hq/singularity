package data

import (
	"singularity/gen/blueprint"
	"singularity/internal/enum"
)

type Server struct {
	Blueprint     blueprint.Blueprint `json:"blueprint" bson:"blueprint"`
	Discriminator string              `json:"discriminator" bson:"discriminator"`
	Port          int                 `json:"port" bson:"port"`
	MetricsPort   *int                `json:"metrics_port" bson:"metrics_port"`
	Status        enum.Status         `json:"status" bson:"status"`
	Report        *ServerReport       `json:"report" bson:"report"`
}

func (server Server) Id() string {
	return server.Blueprint.Id + "-" + server.Discriminator
}

func (server Server) Name() string {
	return server.Blueprint.Name + " #" + server.Discriminator
}
