package policy

import (
	"schoolbusauth/internal/domain"
	"testing"
	"time"
)

func policyBundle() domain.AuthorizationBundle {
	return domain.AuthorizationBundle{Authorization: domain.Authorization{ID: "a1", AuthorizedDate: "2026-09-01", TeacherName: "王老师", Reason: "临时安排", Status: domain.StatusDraft, Revision: 1}, Student: domain.Student{ID: "s1", Name: "李明", ClassName: "三班", SchoolNumber: "S1"}, Route: domain.Route{ID: "r1", Name: "东线", BusNumber: "B1", PickupPoint: "南门", GuardPost: "东门"}, Guardian: domain.Guardian{ID: "g1", Name: "李华", Relationship: "父亲", DocumentLastFour: "4832"}}
}

func TestIssuePolicyAndGateVerification(t *testing.T) {
	decision := New("测试学校").EvaluateIssue(policyBundle(), time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC))
	if !decision.Allowed || decision.BlockingError() != nil {
		t.Fatal("valid bundle rejected")
	}
	answers := DefaultAnswers(decision.Checklist, "已现场核对")
	verified := Verify(decision.Checklist, answers)
	if !verified.Complete || verified.Error() != nil {
		t.Fatal("complete checks rejected")
	}
	verified = Verify(decision.Checklist, answers[:1])
	if verified.Complete || verified.Error() == nil {
		t.Fatal("incomplete checks accepted")
	}
}

func TestPrivacyReviewMasksDocument(t *testing.T) {
	review := ReviewDisclosure(policyBundle(), ChannelPrint)
	if !review.Allowed {
		t.Fatal(review.Error())
	}
	values := MaskForChannel(review.Fields, ChannelExport)
	if values["document_last_four"] != "****4832" {
		t.Fatal("document suffix changed")
	}
}
