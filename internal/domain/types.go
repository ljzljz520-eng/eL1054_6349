package domain

import "time"

type AuthorizationStatus string

const (
	StatusDraft   AuthorizationStatus = "draft"
	StatusIssued  AuthorizationStatus = "issued"
	StatusRevoked AuthorizationStatus = "revoked"
)

type Student struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ClassName    string `json:"class_name"`
	SchoolNumber string `json:"school_number"`
}

type Route struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BusNumber   string `json:"bus_number"`
	PickupPoint string `json:"pickup_point"`
	GuardPost   string `json:"guard_post"`
}

type Guardian struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Relationship     string `json:"relationship"`
	DocumentLastFour string `json:"document_last_four"`
}

type Authorization struct {
	ID             string              `json:"id"`
	StudentID      string              `json:"student_id"`
	RouteID        string              `json:"route_id"`
	GuardianID     string              `json:"guardian_id"`
	AuthorizedDate string              `json:"authorized_date"`
	TeacherName    string              `json:"teacher_name"`
	Reason         string              `json:"reason"`
	Status         AuthorizationStatus `json:"status"`
	CreatedAt      time.Time           `json:"created_at"`
	IssuedAt       *time.Time          `json:"issued_at,omitempty"`
	RevokedAt      *time.Time          `json:"revoked_at,omitempty"`
	Revision       int                 `json:"revision"`
}

type AuthorizationBundle struct {
	Authorization Authorization
	Student       Student
	Route         Route
	Guardian      Guardian
}

type AuditEvent struct {
	ID              string    `json:"id"`
	AuthorizationID string    `json:"authorization_id"`
	Action          string    `json:"action"`
	Actor           string    `json:"actor"`
	OccurredAt      time.Time `json:"occurred_at"`
	Details         string    `json:"details"`
}

type PrintRecord struct {
	ID              string    `json:"id"`
	AuthorizationID string    `json:"authorization_id"`
	Printer         string    `json:"printer"`
	Copies          int       `json:"copies"`
	PrintedAt       time.Time `json:"printed_at"`
}

type CreateRequest struct {
	Student        Student
	Route          Route
	Guardian       Guardian
	AuthorizedDate string
	TeacherName    string
	Reason         string
}

type SearchFilter struct {
	Date        string
	StudentName string
	RouteName   string
	Status      AuthorizationStatus
}
