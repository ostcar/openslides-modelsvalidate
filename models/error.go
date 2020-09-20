package models

import (
	"fmt"
	"strings"
)

type errorList struct {
	name   string
	intent int
	errs   []error
}

func (e *errorList) append(err error) {
	if err == nil {
		return
	}

	e.errs = append(e.errs, err)
}

func (e *errorList) Error() string {
	intent := strings.Repeat(" ", e.intent)
	var msgs []string
	for _, err := range e.errs {
		msgs = append(msgs, fmt.Sprintf("%s* %v", intent, err))
	}
	msg := strings.Join(msgs, "\n")
	if e.name != "" {
		return fmt.Sprintf("%s:\n%s", e.name, msg)
	}
	return msg
}

func (e *errorList) empty() bool {
	return len(e.errs) == 0
}
