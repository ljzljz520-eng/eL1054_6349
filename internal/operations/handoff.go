package operations

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"sort"
	"strings"
	"time"
)

type HandoffSheet struct {
	Date        string       `json:"date"`
	GuardPost   string       `json:"guard_post"`
	Rows        []HandoffRow `json:"rows"`
	GeneratedAt string       `json:"generated_at"`
}

type HandoffRow struct {
	Sequence        int    `json:"sequence"`
	AuthorizationID string `json:"authorization_id"`
	Student         string `json:"student"`
	ClassName       string `json:"class_name"`
	Route           string `json:"route"`
	BusNumber       string `json:"bus_number"`
	Guardian        string `json:"guardian"`
	DocumentSuffix  string `json:"document_suffix"`
	Teacher         string `json:"teacher"`
	ReleaseTime     string `json:"release_time"`
	GuardSignature  string `json:"guard_signature"`
}

func BuildHandoffSheet(bundles []domain.AuthorizationBundle, guardPost string, now time.Time) HandoffSheet {
	rows := make([]HandoffRow, 0)
	for _, bundle := range bundles {
		if bundle.Route.GuardPost != guardPost {
			continue
		}
		if !bundle.Authorization.IsUsable(now) {
			continue
		}
		rows = append(rows, HandoffRow{
			AuthorizationID: bundle.Authorization.ID,
			Student:         bundle.Student.Name,
			ClassName:       bundle.Student.ClassName,
			Route:           bundle.Route.Name,
			BusNumber:       bundle.Route.BusNumber,
			Guardian:        bundle.Guardian.Name,
			DocumentSuffix:  "****" + bundle.Guardian.DocumentLastFour,
			Teacher:         bundle.Authorization.TeacherName,
			ReleaseTime:     "____________",
			GuardSignature:  "____________",
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Route == rows[j].Route {
			if rows[i].ClassName == rows[j].ClassName {
				return rows[i].Student < rows[j].Student
			}
			return rows[i].ClassName < rows[j].ClassName
		}
		return rows[i].Route < rows[j].Route
	})
	for index := range rows {
		rows[index].Sequence = index + 1
	}
	return HandoffSheet{Date: now.Format("2006-01-02"), GuardPost: guardPost, Rows: rows, GeneratedAt: now.UTC().Format(time.RFC3339)}
}

func (s HandoffSheet) PlainText() string {
	var output strings.Builder
	output.WriteString("校车临时接送门岗交接表\n")
	output.WriteString("日期: " + s.Date + "  门岗: " + s.GuardPost + "\n")
	output.WriteString(strings.Repeat("=", 96) + "\n")
	for _, row := range s.Rows {
		output.WriteString(fmt.Sprintf("%02d  %-10s %-10s %-12s %-10s %-10s %s\n", row.Sequence, row.Student, row.ClassName, row.Route, row.Guardian, row.DocumentSuffix, row.ReleaseTime))
	}
	if len(s.Rows) == 0 {
		output.WriteString("今日暂无可用授权\n")
	}
	output.WriteString(strings.Repeat("-", 96) + "\n")
	output.WriteString("交班门卫签字: ______________  接班门卫签字: ______________\n")
	return output.String()
}
