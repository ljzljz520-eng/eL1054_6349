package validation

import (
	"schoolbusauth/internal/domain"
	"testing"
)

func TestRequestValidationAndNormalization(t *testing.T) {
	req := domain.CreateRequest{Student: domain.Student{Name: "  李  明 ", ClassName: " 三班 ", SchoolNumber: " s-1 "}, Route: domain.Route{Name: " 东线 ", BusNumber: " b1 ", PickupPoint: " 南门 "}, Guardian: domain.Guardian{Name: " 李华 ", Relationship: " 父亲 ", DocumentLastFour: "a123"}, AuthorizedDate: "2026-09-01", TeacherName: " 王老师 "}
	req = Normalize(req)
	if req.Student.Name != "李 明" || req.Student.SchoolNumber != "S-1" || req.Guardian.DocumentLastFour != "A123" {
		t.Fatalf("normalization failed: %#v", req)
	}
	if err := CreateRequest(req); err != nil {
		t.Fatal(err)
	}
	req.Guardian.DocumentLastFour = "12"
	if err := CreateRequest(req); err == nil {
		t.Fatal("invalid document suffix accepted")
	}
}

func TestPrintCopies(t *testing.T) {
	if Copies(1) != nil || Copies(10) != nil {
		t.Fatal("valid copies rejected")
	}
	if Copies(0) == nil || Copies(11) == nil {
		t.Fatal("invalid copies accepted")
	}
}
