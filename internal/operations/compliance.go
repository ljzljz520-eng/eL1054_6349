package operations

import (
	"schoolbusauth/internal/domain"
	"sort"
	"time"
)

type ComplianceReport struct {
	AuthorizationID string              `json:"authorization_id"`
	Consistent      bool                `json:"consistent"`
	Findings        []ComplianceFinding `json:"findings"`
	Lifecycle       []string            `json:"lifecycle"`
}

type ComplianceFinding struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity int    `json:"severity"`
}

func ReviewCompliance(bundle domain.AuthorizationBundle, audits []domain.AuditEvent, prints []domain.PrintRecord, now time.Time) ComplianceReport {
	result := ComplianceReport{AuthorizationID: bundle.Authorization.ID, Consistent: true, Findings: make([]ComplianceFinding, 0), Lifecycle: make([]string, 0)}
	if bundle.Authorization.StudentID != bundle.Student.ID {
		result.add("student_reference", "授权引用的学生记录不一致", 100)
	}
	if bundle.Authorization.RouteID != bundle.Route.ID {
		result.add("route_reference", "授权引用的线路记录不一致", 100)
	}
	if bundle.Authorization.GuardianID != bundle.Guardian.ID {
		result.add("guardian_reference", "授权引用的接送人记录不一致", 100)
	}
	if bundle.Authorization.Revision < 1 {
		result.add("revision", "授权修订号无效", 90)
	}
	if bundle.Authorization.Status == domain.StatusIssued && bundle.Authorization.IssuedAt == nil {
		result.add("issued_time", "已签发授权缺少签发时间", 90)
	}
	if bundle.Authorization.Status == domain.StatusRevoked && bundle.Authorization.RevokedAt == nil {
		result.add("revoked_time", "已撤销授权缺少撤销时间", 90)
	}
	if bundle.Authorization.Status == domain.StatusDraft && bundle.Authorization.IssuedAt != nil {
		result.add("draft_issued", "草稿不应包含签发时间", 80)
	}
	if bundle.Authorization.IsExpired(now) && len(prints) > 0 {
		for _, print := range prints {
			if print.PrintedAt.After(now) {
				result.add("future_print", "打印记录时间晚于审查时间", 70)
			}
		}
	}
	created := false
	issued := false
	revoked := false
	for _, audit := range audits {
		result.Lifecycle = append(result.Lifecycle, audit.OccurredAt.UTC().Format(time.RFC3339)+" "+audit.Action+" by "+audit.Actor)
		switch audit.Action {
		case "created":
			created = true
		case "issued":
			issued = true
		case "revoked":
			revoked = true
		case "updated":
		default:
			result.add("unknown_audit", "存在未知审计动作: "+audit.Action, 30)
		}
	}
	if !created {
		result.add("created_audit", "缺少创建审计记录", 80)
	}
	if bundle.Authorization.Status == domain.StatusIssued && !issued {
		result.add("issued_audit", "缺少签发审计记录", 80)
	}
	if bundle.Authorization.Status == domain.StatusRevoked && !revoked {
		result.add("revoked_audit", "缺少撤销审计记录", 80)
	}
	for _, print := range prints {
		if print.AuthorizationID != bundle.Authorization.ID {
			result.add("print_reference", "打印记录引用其他授权", 70)
		}
		if print.Copies < 1 || print.Copies > 10 {
			result.add("print_copies", "打印份数超出允许范围", 50)
		}
		if print.Printer == "" {
			result.add("print_actor", "打印记录缺少打印人", 50)
		}
	}
	sort.Slice(result.Findings, func(i, j int) bool {
		if result.Findings[i].Severity == result.Findings[j].Severity {
			return result.Findings[i].Code < result.Findings[j].Code
		}
		return result.Findings[i].Severity > result.Findings[j].Severity
	})
	sort.Strings(result.Lifecycle)
	return result
}

func (r *ComplianceReport) add(code, message string, severity int) {
	r.Consistent = false
	r.Findings = append(r.Findings, ComplianceFinding{Code: code, Message: message, Severity: severity})
}
