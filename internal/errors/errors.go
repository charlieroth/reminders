package errors

import "fmt"

type DuplicateTaskError struct {
	Title string
}

func (e DuplicateTaskError) Error() string {
	return fmt.Sprintf("task with title %s already exists", e.Title)
}
