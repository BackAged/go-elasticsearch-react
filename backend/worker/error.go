package worker

import (
	"fmt"
)

// NewErrAlreadyRegisteredTask returns error
func NewErrAlreadyRegisteredTask(task string) error {
	return fmt.Errorf("task-%s already registered", task)
}
