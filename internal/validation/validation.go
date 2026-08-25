package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"schoolbusauth/internal/domain"
)

type FieldError struct{ Field, Message string }
type Errors []FieldError

func (e Errors) Error() string {
	parts := make([]string, 0, len(e))
	for _, item := range e {
		parts = append(parts, item.Field+": "+item.Message)
	}
	return strings.Join(parts, "; ")
}

var lastFourPattern = regexp.MustCompile(`^[0-9A-Za-z]{4}$`)

func CreateRequest(req domain.CreateRequest) error {
	var errs Errors
	student(&errs, req.Student)
	route(&errs, req.Route)
	guardian(&errs, req.Guardian)
	if strings.TrimSpace(req.TeacherName) == "" {
		errs = append(errs, FieldError{"teacher_name", "is required"})
	}
	if len([]rune(strings.TrimSpace(req.TeacherName))) > 40 {
		errs = append(errs, FieldError{"teacher_name", "is too long"})
	}
	if _, err := time.Parse("2006-01-02", req.AuthorizedDate); err != nil {
		errs = append(errs, FieldError{"authorized_date", "must use YYYY-MM-DD"})
	}
	if len([]rune(req.Reason)) > 200 {
		errs = append(errs, FieldError{"reason", "is too long"})
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func student(errs *Errors, value domain.Student) {
	if strings.TrimSpace(value.Name) == "" {
		*errs = append(*errs, FieldError{"student.name", "is required"})
	}
	if strings.TrimSpace(value.ClassName) == "" {
		*errs = append(*errs, FieldError{"student.class_name", "is required"})
	}
	if strings.TrimSpace(value.SchoolNumber) == "" {
		*errs = append(*errs, FieldError{"student.school_number", "is required"})
	}
}

func route(errs *Errors, value domain.Route) {
	if strings.TrimSpace(value.Name) == "" {
		*errs = append(*errs, FieldError{"route.name", "is required"})
	}
	if strings.TrimSpace(value.BusNumber) == "" {
		*errs = append(*errs, FieldError{"route.bus_number", "is required"})
	}
	if strings.TrimSpace(value.PickupPoint) == "" {
		*errs = append(*errs, FieldError{"route.pickup_point", "is required"})
	}
}

func guardian(errs *Errors, value domain.Guardian) {
	if strings.TrimSpace(value.Name) == "" {
		*errs = append(*errs, FieldError{"guardian.name", "is required"})
	}
	if strings.TrimSpace(value.Relationship) == "" {
		*errs = append(*errs, FieldError{"guardian.relationship", "is required"})
	}
	if !lastFourPattern.MatchString(value.DocumentLastFour) {
		*errs = append(*errs, FieldError{"guardian.document_last_four", "must contain four letters or digits"})
	}
}

func Copies(value int) error {
	if value < 1 {
		return fmt.Errorf("copies must be positive")
	}
	if value > 10 {
		return fmt.Errorf("copies must not exceed 10")
	}
	return nil
}
