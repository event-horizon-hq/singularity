package dto

import "singularity/gen/blueprint"

type BlueprintResponse struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Type        string               `json:"type"`
	Author      *string              `json:"author,omitempty"`
	Image       string               `json:"image"`
	Resources   ResourcesResponse    `json:"resources"`
	Environment map[string]string    `json:"environment,omitempty"`
	Volumes     []VolumeResponse     `json:"volumes,omitempty"`
	NetworkMode string               `json:"networkMode"`
	HealthCheck *HealthCheckResponse `json:"healthCheck,omitempty"`
	Entrypoint  *[]string            `json:"entrypoint,omitempty"`
	Cmd         *[]string            `json:"cmd,omitempty"`
	WorkingDir  *string              `json:"workingDir,omitempty"`
	UserId      *uint                `json:"userId,omitempty"`
	GroupId     *uint                `json:"groupId,omitempty"`
	Logging     *LoggingResponse     `json:"logging,omitempty"`
}

type ResourcesResponse struct {
	Memory uint    `json:"memory"`
	Cpu    float64 `json:"cpu"`
}

type VolumeResponse struct {
	Id           string `json:"id"`
	TargetFolder string `json:"targetFolder"`
	ReadOnly     bool   `json:"readOnly"`
}

type HealthCheckResponse struct {
	Test        []string `json:"test"`
	Interval    uint     `json:"interval"`
	Timeout     uint     `json:"timeout"`
	Retries     uint     `json:"retries"`
	StartPeriod uint     `json:"startPeriod"`
}

type LoggingResponse struct {
	Driver  string            `json:"driver"`
	Options map[string]string `json:"options,omitempty"`
}

func NewBlueprintResponse(bp *blueprint.Blueprint) BlueprintResponse {
	var volumes []VolumeResponse
	if len(bp.Volumes) > 0 {
		volumes = make([]VolumeResponse, len(bp.Volumes))
		for i, v := range bp.Volumes {
			volumes[i] = VolumeResponse{
				Id:           v.Id,
				TargetFolder: v.TargetFolder,
				ReadOnly:     v.ReadOnly,
			}
		}
	}

	var healthCheck *HealthCheckResponse
	if bp.HealthCheck != nil {
		healthCheck = &HealthCheckResponse{
			Test:        bp.HealthCheck.Test,
			Interval:    bp.HealthCheck.Interval,
			Timeout:     bp.HealthCheck.Timeout,
			Retries:     bp.HealthCheck.Retries,
			StartPeriod: bp.HealthCheck.StartPeriod,
		}
	}

	var logging *LoggingResponse
	if bp.Logging != nil {
		logging = &LoggingResponse{
			Driver:  bp.Logging.Driver,
			Options: bp.Logging.Options,
		}
	}

	var environment map[string]string
	if len(bp.Environment) > 0 {
		environment = bp.Environment
	}

	return BlueprintResponse{
		Id:     bp.Id,
		Name:   bp.Name,
		Type:   bp.Type,
		Author: bp.Author,
		Image:  bp.Image,
		Resources: ResourcesResponse{
			Memory: bp.Resources.Memory,
			Cpu:    bp.Resources.Cpu,
		},
		Environment: environment,
		Volumes:     volumes,
		NetworkMode: bp.NetworkMode,
		HealthCheck: healthCheck,
		Entrypoint:  bp.Entrypoint,
		Cmd:         bp.Cmd,
		WorkingDir:  bp.WorkingDir,
		UserId:      bp.UserId,
		GroupId:     bp.GroupId,
		Logging:     logging,
	}
}
