package enum

type Status int8

const (
	StatusInactive Status = iota
	StatusError
	StatusCreating
	StatusActive
	StatusRestarting
)

func (s Status) IsValid() bool {
	switch s {
	case StatusInactive, StatusError, StatusCreating, StatusActive, StatusRestarting:
		return true
	default:
		return false
	}
}
