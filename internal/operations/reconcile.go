package operations

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"sort"
	"strings"
	"time"
)

type GateEntry struct {
	AuthorizationID string    `json:"authorization_id"`
	StudentName     string    `json:"student_name"`
	GuardianName    string    `json:"guardian_name"`
	DocumentSuffix  string    `json:"document_suffix"`
	GuardPost       string    `json:"guard_post"`
	ReleasedAt      time.Time `json:"released_at"`
	GuardName       string    `json:"guard_name"`
}

type ReconciliationStatus string

const (
	Reconciled       ReconciliationStatus = "reconciled"
	MissingEntry     ReconciliationStatus = "missing_entry"
	UnexpectedEntry  ReconciliationStatus = "unexpected_entry"
	IdentityMismatch ReconciliationStatus = "identity_mismatch"
	TimingMismatch   ReconciliationStatus = "timing_mismatch"
)

type ReconciliationRow struct {
	AuthorizationID string               `json:"authorization_id"`
	StudentName     string               `json:"student_name"`
	Status          ReconciliationStatus `json:"status"`
	Messages        []string             `json:"messages"`
	ReleasedAt      *time.Time           `json:"released_at,omitempty"`
}

type ReconciliationReport struct {
	Date       string              `json:"date"`
	Matched    int                 `json:"matched"`
	Exceptions int                 `json:"exceptions"`
	Rows       []ReconciliationRow `json:"rows"`
}

func ReconcileGateEntries(bundles []domain.AuthorizationBundle, entries []GateEntry, date time.Time) ReconciliationReport {
	authMap := make(map[string]domain.AuthorizationBundle, len(bundles))
	for _, bundle := range bundles {
		if bundle.Authorization.AuthorizedDate == date.Format("2006-01-02") {
			authMap[bundle.Authorization.ID] = bundle
		}
	}
	entryMap := make(map[string][]GateEntry)
	for _, entry := range entries {
		entryMap[entry.AuthorizationID] = append(entryMap[entry.AuthorizationID], entry)
	}
	report := ReconciliationReport{Date: date.Format("2006-01-02"), Rows: make([]ReconciliationRow, 0)}
	for id, bundle := range authMap {
		matching := entryMap[id]
		row := ReconciliationRow{AuthorizationID: id, StudentName: bundle.Student.Name, Status: Reconciled, Messages: make([]string, 0)}
		if len(matching) == 0 {
			if bundle.Authorization.Status == domain.StatusIssued {
				row.Status = MissingEntry
				row.Messages = append(row.Messages, "已签发授权没有门岗放行记录")
			}
		} else {
			entry := matching[0]
			row.ReleasedAt = &entry.ReleasedAt
			if entry.StudentName != bundle.Student.Name {
				row.Status = IdentityMismatch
				row.Messages = append(row.Messages, "学生姓名与授权不一致")
			}
			if entry.GuardianName != bundle.Guardian.Name {
				row.Status = IdentityMismatch
				row.Messages = append(row.Messages, "接送人姓名与授权不一致")
			}
			if entry.DocumentSuffix != bundle.Guardian.DocumentLastFour {
				row.Status = IdentityMismatch
				row.Messages = append(row.Messages, "证件后四位与授权不一致")
			}
			if entry.GuardPost != bundle.Route.GuardPost {
				row.Status = IdentityMismatch
				row.Messages = append(row.Messages, "实际放行门岗与授权不一致")
			}
			if entry.GuardName == "" {
				row.Status = IdentityMismatch
				row.Messages = append(row.Messages, "放行记录缺少门卫姓名")
			}
			if entry.ReleasedAt.Format("2006-01-02") != report.Date {
				row.Status = TimingMismatch
				row.Messages = append(row.Messages, "放行时间不在授权当天")
			}
			if len(matching) > 1 {
				row.Status = IdentityMismatch
				row.Messages = append(row.Messages, fmt.Sprintf("同一授权存在 %d 条放行记录", len(matching)))
			}
		}
		if row.Status == Reconciled {
			report.Matched++
		} else {
			report.Exceptions++
		}
		report.Rows = append(report.Rows, row)
	}
	for id, matching := range entryMap {
		if _, exists := authMap[id]; exists {
			continue
		}
		for _, entry := range matching {
			released := entry.ReleasedAt
			report.Rows = append(report.Rows, ReconciliationRow{AuthorizationID: id, StudentName: entry.StudentName, Status: UnexpectedEntry, Messages: []string{"门岗记录找不到当日授权"}, ReleasedAt: &released})
			report.Exceptions++
		}
	}
	sort.Slice(report.Rows, func(i, j int) bool {
		if report.Rows[i].Status == report.Rows[j].Status {
			return report.Rows[i].StudentName < report.Rows[j].StudentName
		}
		return report.Rows[i].Status < report.Rows[j].Status
	})
	return report
}

func (r ReconciliationReport) Error() error {
	if r.Exceptions == 0 {
		return nil
	}
	items := make([]string, 0)
	for _, row := range r.Rows {
		if row.Status != Reconciled {
			items = append(items, row.AuthorizationID+": "+strings.Join(row.Messages, ", "))
		}
	}
	return fmt.Errorf("gate reconciliation has %d exceptions: %s", r.Exceptions, strings.Join(items, "; "))
}
