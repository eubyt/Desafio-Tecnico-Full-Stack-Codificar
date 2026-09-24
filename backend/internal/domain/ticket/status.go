package ticket

import "strings"

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

func ParseStatus(val string) (Status, error) {
	clean := strings.ToLower(strings.TrimSpace(val))
	switch clean {
	case "open":
		return StatusOpen, nil
	case "in_progress":
		return StatusInProgress, nil
	case "resolved":
		return StatusResolved, nil
	case "closed":
		return StatusClosed, nil
	default:
		return "", ErrInvalidStatus
	}
}

func (s Status) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusResolved, StatusClosed:
		return true
	default:
		return false
	}
}

func (s Status) IsOpen() bool {
	return s == StatusOpen || s == StatusInProgress
}

func (s Status) CanTransitionTo(next Status) bool {
	if !next.IsValid() {
		return false
	}
	if s == next {
		return true
	}

	switch s {
	case StatusOpen:
		return next == StatusInProgress || next == StatusResolved || next == StatusClosed
	case StatusInProgress:
		return next == StatusResolved || next == StatusClosed || next == StatusOpen
	case StatusResolved:
		return next == StatusClosed || next == StatusInProgress || next == StatusOpen
	case StatusClosed:
		return next == StatusOpen || next == StatusInProgress
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}
