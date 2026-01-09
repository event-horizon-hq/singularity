package data

type ServerReport struct {
	OnlineCount int32 `json:"onlineCount" bson:"onlineCount"`
	OnlineSince int64 `json:"onlineSince" bson:"onlineSince"`
	MemoryUsage int64 `json:"memoryUsage" bson:"memoryUsage"`
	TotalMemory int64 `json:"totalMemory" bson:"totalMemory"`
	CpuUsage    int64 `json:"cpuUsage" bson:"cpuUsage"`
}
