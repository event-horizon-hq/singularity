package data

import "singularity/internal/enum"

type Server struct {
	Blueprint     Blueprint     `json:"blueprint" bson:"blueprint"`
	Discriminator string        `json:"discriminator" bson:"discriminator"`
	Port          int32         `json:"port" bson:"port"`
	Status        enum.Status   `bson:"status" bson:"status"`
	Report        *ServerReport `json:"report" bson:"report"`
}

func (server Server) Id() string {
	return server.Blueprint.Id + "-" + server.Discriminator
}

func (server Server) Name() string {
	return server.Blueprint.Name + " #" + server.Discriminator
}
