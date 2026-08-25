package validation

import (
	"schoolbusauth/internal/domain"
	"strings"
)

func Normalize(req domain.CreateRequest) domain.CreateRequest {
	req.Student.Name = clean(req.Student.Name)
	req.Student.ClassName = clean(req.Student.ClassName)
	req.Student.SchoolNumber = strings.ToUpper(clean(req.Student.SchoolNumber))
	req.Route.Name = clean(req.Route.Name)
	req.Route.BusNumber = strings.ToUpper(clean(req.Route.BusNumber))
	req.Route.PickupPoint = clean(req.Route.PickupPoint)
	req.Route.GuardPost = clean(req.Route.GuardPost)
	req.Guardian.Name = clean(req.Guardian.Name)
	req.Guardian.Relationship = clean(req.Guardian.Relationship)
	req.Guardian.DocumentLastFour = strings.ToUpper(clean(req.Guardian.DocumentLastFour))
	req.TeacherName = clean(req.TeacherName)
	req.Reason = clean(req.Reason)
	return req
}

func MaskDocument(lastFour string) string {
	if len(lastFour) != 4 {
		return "****"
	}
	return "**************" + strings.ToUpper(lastFour)
}

func clean(value string) string { return strings.Join(strings.Fields(value), " ") }
