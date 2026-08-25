package document

import (
	"encoding/json"
	"fmt"
)

type MachineDocument struct {
	SchemaVersion   int                 `json:"schema_version"`
	DocumentType    string              `json:"document_type"`
	AuthorizationID string              `json:"authorization_id"`
	Validity        MachineValidity     `json:"validity"`
	Student         MachineStudent      `json:"student"`
	Transport       MachineTransport    `json:"transport"`
	PickupPerson    MachinePickupPerson `json:"pickup_person"`
	Issuer          MachineIssuer       `json:"issuer"`
	Verification    []string            `json:"verification"`
}

type MachineValidity struct {
	Date   string `json:"date"`
	Status string `json:"status"`
	Notice string `json:"notice"`
}
type MachineStudent struct {
	Name         string `json:"name"`
	ClassName    string `json:"class_name"`
	SchoolNumber string `json:"school_number"`
}
type MachineTransport struct {
	Route       string `json:"route"`
	BusNumber   string `json:"bus_number"`
	PickupPoint string `json:"pickup_point"`
	GuardPost   string `json:"guard_post"`
}
type MachinePickupPerson struct {
	Name           string `json:"name"`
	Relationship   string `json:"relationship"`
	MaskedDocument string `json:"masked_document"`
}
type MachineIssuer struct {
	Teacher     string `json:"teacher"`
	GeneratedAt string `json:"generated_at"`
}

func BuildMachineDocument(view View) MachineDocument {
	return MachineDocument{
		SchemaVersion:   1,
		DocumentType:    "school_bus_temporary_pickup_authorization",
		AuthorizationID: view.AuthorizationID,
		Validity:        MachineValidity{Date: view.AuthorizedDate, Status: view.Status, Notice: view.NoticeTitle},
		Student:         MachineStudent{Name: view.StudentName, ClassName: view.ClassName, SchoolNumber: view.SchoolNumber},
		Transport:       MachineTransport{Route: view.RouteName, BusNumber: view.BusNumber, PickupPoint: view.PickupPoint, GuardPost: view.GuardPost},
		PickupPerson:    MachinePickupPerson{Name: view.GuardianName, Relationship: view.Relationship, MaskedDocument: view.MaskedDocument},
		Issuer:          MachineIssuer{Teacher: view.TeacherName, GeneratedAt: view.GeneratedAt},
		Verification:    append([]string(nil), view.VerificationItems...),
	}
}

func RenderJSON(view View) ([]byte, error) {
	if view.AuthorizationID == "" {
		return nil, fmt.Errorf("authorization id is required")
	}
	data, err := json.MarshalIndent(BuildMachineDocument(view), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("render machine document: %w", err)
	}
	return append(data, '\n'), nil
}
