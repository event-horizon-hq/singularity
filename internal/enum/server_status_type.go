package enum

type Status string

const (
	StatusInactive   Status = "INACTIVE"
	StatusError      Status = "ERROR"
	StatusCreating   Status = "CREATING"
	StatusActive     Status = "ACTIVE"
	StatusRestarting Status = "RESTARTING"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusInactive, StatusError, StatusCreating, StatusActive, StatusRestarting:
		return true
	default:
		return false
	}
}
