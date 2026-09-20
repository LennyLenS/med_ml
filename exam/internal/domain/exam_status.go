package domain

import "fmt"

type ExamStatus string

const (
	// узи создано
	ExamStatusNew ExamStatus = "new"
	// узи обрабатывается
	ExamStatusPending ExamStatus = "pending"
	// узи обработано
	ExamStatusCompleted ExamStatus = "completed"
)

func (s ExamStatus) String() string {
	return string(s)
}

func (s ExamStatus) Parse(status string) (ExamStatus, error) {
	switch status {
	case "new":
		return ExamStatusNew, nil
	case "pending":
		return ExamStatusPending, nil
	case "completed":
		return ExamStatusCompleted, nil
	default:
		return "", fmt.Errorf("invalid status: %s", status)
	}
}
