package domain

import (
	"testing"
	"time"
)

func TestAuthorizationDateRules(t *testing.T) {
	today := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	issued := Authorization{AuthorizedDate: "2026-09-02", Status: StatusIssued, Revision: 1}
	if !issued.IsUsable(today) || issued.IsExpired(today) {
		t.Fatal("today authorization should be usable")
	}
	yesterday := issued
	yesterday.AuthorizedDate = "2026-09-01"
	if !yesterday.IsExpired(today) || yesterday.IsUsable(today) {
		t.Fatal("past authorization should be expired")
	}
	if (Authorization{Status: StatusIssued, Revision: 1}).CanIssue() == nil {
		t.Fatal("issued authorization cannot issue again")
	}
}

func TestBundleSearchMatching(t *testing.T) {
	bundle := AuthorizationBundle{Authorization: Authorization{AuthorizedDate: "2026-09-01", Status: StatusDraft}, Student: Student{Name: "Li Ming"}, Route: Route{Name: "East Lake"}}
	if !bundle.Matches(SearchFilter{StudentName: "ming", RouteName: "east", Status: StatusDraft}) {
		t.Fatal("expected match")
	}
	if bundle.Matches(SearchFilter{Date: "2026-09-02"}) {
		t.Fatal("unexpected date match")
	}
}
