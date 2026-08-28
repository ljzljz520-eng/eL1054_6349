package document

import (
	"schoolbusauth/internal/domain"
	"strings"
	"testing"
	"time"
)

func documentBundle(date string, status domain.AuthorizationStatus) domain.AuthorizationBundle {
	return domain.AuthorizationBundle{Authorization: domain.Authorization{ID: "a1", AuthorizedDate: date, TeacherName: "王老师", Status: status}, Student: domain.Student{Name: "李明", ClassName: "三班", SchoolNumber: "S1"}, Route: domain.Route{Name: "东线", BusNumber: "B1", PickupPoint: "南门", GuardPost: "东门"}, Guardian: domain.Guardian{Name: "李华", Relationship: "父亲", DocumentLastFour: "4832"}}
}

func TestExpiredDocumentIsProminent(t *testing.T) {
	view := BuildView(documentBundle("2026-08-31", domain.StatusIssued), time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC))
	if view.NoticeLevel != NoticeExpired {
		t.Fatal("expected expired notice")
	}
	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal(err)
	}
	html, err := renderer.Render(view)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), "notice expired") || !strings.Contains(string(html), "授权已过期") {
		t.Fatal("expired styling missing")
	}
}

func TestDocumentFormats(t *testing.T) {
	view := BuildView(documentBundle("2026-09-01", domain.StatusIssued), time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC))
	text, err := (TextRenderer{}).Render(view)
	if err != nil || !strings.Contains(string(text), "4832") {
		t.Fatal("text render failed")
	}
	jsonData, err := RenderJSON(view)
	if err != nil || !strings.Contains(string(jsonData), "school_bus_temporary") {
		t.Fatal("json render failed")
	}
	labels, err := RenderGateLabels([]View{view}, 2)
	if err != nil || !strings.Contains(labels, "李明") {
		t.Fatal("label render failed")
	}
}
