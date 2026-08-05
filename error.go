package mailofly

import "fmt"

// Error is returned when the API responds with a non-2xx status.
type Error struct {
	Status         int
	Err            string
	DetailMessage  string
	Body           any
}

func (e *Error) Error() string {
	if e.DetailMessage != "" {
		return fmt.Sprintf("%s: %s", e.Err, e.DetailMessage)
	}
	return e.Err
}
