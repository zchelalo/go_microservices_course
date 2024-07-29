package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("name is required")
var ErrStartDateRequired = errors.New("start date is required")
var ErrEndDateRequired = errors.New("end date is required")

type ErrNotFound struct {
	CourseId string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("course '%s' doesn't exist", e.CourseId)
}
