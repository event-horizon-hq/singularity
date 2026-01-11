package data

type ServerReport struct {
	OnlineCount int32 `json:"online_count" bson:"online_count"`
	OnlineSince int64 `json:"online_since" bson:"online_since"`
	MemoryUsage int64 `json:"memory_usage" bson:"memory_usage"`
	TotalMemory int64 `json:"total_memory" bson:"total_memory"`
	CpuUsage    int64 `json:"cpu_usage" bson:"cpu_usage"`
}
