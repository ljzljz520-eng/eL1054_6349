package policy

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"strings"
	"time"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityBlocking Severity = "blocking"
)

type Finding struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Field    string   `json:"field"`
}

type Decision struct {
	Allowed   bool            `json:"allowed"`
	Findings  []Finding       `json:"findings"`
	Checklist []ChecklistItem `json:"checklist"`
}

type ChecklistItem struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Evidence string `json:"evidence"`
}

type Engine struct {
	schoolName      string
	maxReasonLength int
}

func New(schoolName string) Engine {
	if strings.TrimSpace(schoolName) == "" {
		schoolName = "本校"
	}
	return Engine{schoolName: schoolName, maxReasonLength: 200}
}

func (e Engine) EvaluateIssue(bundle domain.AuthorizationBundle, today time.Time) Decision {
	findings := make([]Finding, 0)
	findings = append(findings, e.identityFindings(bundle)...)
	findings = append(findings, e.routeFindings(bundle)...)
	findings = append(findings, e.authorizationFindings(bundle, today)...)
	allowed := true
	for _, finding := range findings {
		if finding.Severity == SeverityBlocking {
			allowed = false
		}
	}
	return Decision{Allowed: allowed, Findings: findings, Checklist: e.GateChecklist(bundle)}
}

func (e Engine) identityFindings(bundle domain.AuthorizationBundle) []Finding {
	result := make([]Finding, 0)
	if strings.TrimSpace(bundle.Student.Name) == "" {
		result = append(result, Finding{"student_name_missing", SeverityBlocking, "学生姓名缺失", "student.name"})
	}
	if strings.TrimSpace(bundle.Student.ClassName) == "" {
		result = append(result, Finding{"class_missing", SeverityBlocking, "学生班级缺失", "student.class_name"})
	}
	if strings.TrimSpace(bundle.Student.SchoolNumber) == "" {
		result = append(result, Finding{"school_number_missing", SeverityBlocking, "学生学号缺失", "student.school_number"})
	}
	if strings.TrimSpace(bundle.Guardian.Name) == "" {
		result = append(result, Finding{"guardian_name_missing", SeverityBlocking, "接送人姓名缺失", "guardian.name"})
	}
	if len(bundle.Guardian.DocumentLastFour) != 4 {
		result = append(result, Finding{"document_suffix_invalid", SeverityBlocking, "证件后四位格式不正确", "guardian.document_last_four"})
	}
	if strings.EqualFold(bundle.Guardian.Relationship, "本人") {
		result = append(result, Finding{"relationship_ambiguous", SeverityWarning, "接送人与学生关系需要进一步核对", "guardian.relationship"})
	}
	return result
}

func (e Engine) routeFindings(bundle domain.AuthorizationBundle) []Finding {
	result := make([]Finding, 0)
	if strings.TrimSpace(bundle.Route.Name) == "" {
		result = append(result, Finding{"route_missing", SeverityBlocking, "校车线路缺失", "route.name"})
	}
	if strings.TrimSpace(bundle.Route.BusNumber) == "" {
		result = append(result, Finding{"bus_missing", SeverityBlocking, "车辆编号缺失", "route.bus_number"})
	}
	if strings.TrimSpace(bundle.Route.PickupPoint) == "" {
		result = append(result, Finding{"pickup_missing", SeverityBlocking, "临时接送点缺失", "route.pickup_point"})
	}
	if strings.TrimSpace(bundle.Route.GuardPost) == "" {
		result = append(result, Finding{"guard_post_missing", SeverityWarning, "未指定核验门岗，将由现场安排", "route.guard_post"})
	}
	return result
}

func (e Engine) authorizationFindings(bundle domain.AuthorizationBundle, today time.Time) []Finding {
	result := make([]Finding, 0)
	date, err := time.Parse("2006-01-02", bundle.Authorization.AuthorizedDate)
	if err != nil {
		return append(result, Finding{"date_invalid", SeverityBlocking, "授权日期格式无效", "authorization.authorized_date"})
	}
	y, m, d := today.Date()
	today = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	if date.Before(today) {
		result = append(result, Finding{"date_expired", SeverityBlocking, "过期授权不能签发", "authorization.authorized_date"})
	}
	if date.After(today) {
		result = append(result, Finding{"date_future", SeverityInfo, "授权将在指定日期生效", "authorization.authorized_date"})
	}
	if bundle.Authorization.Status != domain.StatusDraft {
		result = append(result, Finding{"status_not_draft", SeverityBlocking, "只有草稿可以签发", "authorization.status"})
	}
	if strings.TrimSpace(bundle.Authorization.TeacherName) == "" {
		result = append(result, Finding{"teacher_missing", SeverityBlocking, "班主任姓名缺失", "authorization.teacher_name"})
	}
	if len([]rune(bundle.Authorization.Reason)) > e.maxReasonLength {
		result = append(result, Finding{"reason_too_long", SeverityBlocking, "变更原因过长", "authorization.reason"})
	}
	if strings.TrimSpace(bundle.Authorization.Reason) == "" {
		result = append(result, Finding{"reason_empty", SeverityWarning, "建议填写临时变更原因", "authorization.reason"})
	}
	return result
}

func (e Engine) GateChecklist(bundle domain.AuthorizationBundle) []ChecklistItem {
	return []ChecklistItem{
		{"student", "核对学生姓名、班级与学号", true, bundle.Student.Name + " / " + bundle.Student.ClassName + " / " + bundle.Student.SchoolNumber},
		{"date", "确认授权日期是放行当天", true, bundle.Authorization.AuthorizedDate},
		{"guardian", "核对临时接送人姓名", true, bundle.Guardian.Name},
		{"relationship", "口头确认与学生关系", true, bundle.Guardian.Relationship},
		{"document", "核对证件后四位", true, "****" + bundle.Guardian.DocumentLastFour},
		{"route", "确认线路与车辆", true, bundle.Route.Name + " / " + bundle.Route.BusNumber},
		{"pickup", "告知临时接送点", true, bundle.Route.PickupPoint},
		{"signature", "请接送人签字确认", true, "现场签字"},
		{"time", "登记实际放行时间", true, "门卫填写"},
		{"exception", "异常情况联系班主任", false, bundle.Authorization.TeacherName},
	}
}

func (d Decision) BlockingError() error {
	if d.Allowed {
		return nil
	}
	messages := make([]string, 0)
	for _, finding := range d.Findings {
		if finding.Severity == SeverityBlocking {
			messages = append(messages, finding.Message)
		}
	}
	return fmt.Errorf("policy rejected authorization: %s", strings.Join(messages, "; "))
}
