package document

import (
	"fmt"
	"strings"
)

type GateLabel struct {
	AuthorizationID  string
	PrimaryLine      string
	SecondaryLine    string
	VerificationLine string
	Warning          string
}

func BuildGateLabel(view View) GateLabel {
	warning := view.NoticeTitle
	if view.NoticeLevel == NoticeValid {
		warning = "当日有效 · 核验证件"
	}
	return GateLabel{
		AuthorizationID:  view.AuthorizationID,
		PrimaryLine:      view.StudentName + " · " + view.ClassName,
		SecondaryLine:    view.RouteName + " / " + view.BusNumber + " / " + view.PickupPoint,
		VerificationLine: view.GuardianName + " (" + view.Relationship + ") " + view.MaskedDocument,
		Warning:          warning,
	}
}

func RenderGateLabels(views []View, columns int) (string, error) {
	if columns < 1 || columns > 4 {
		return "", fmt.Errorf("label columns must be between 1 and 4")
	}
	var output strings.Builder
	output.WriteString("<!doctype html><html lang=\"zh-CN\"><head><meta charset=\"utf-8\"><style>")
	output.WriteString("body{font-family:sans-serif;margin:8mm}.labels{display:grid;grid-template-columns:repeat(")
	output.WriteString(fmt.Sprintf("%d", columns))
	output.WriteString(",1fr);gap:4mm}.label{border:1px solid #222;padding:4mm;break-inside:avoid}.warning{font-weight:bold;color:#a00}.id{font-size:10px;color:#555}</style></head><body><div class=\"labels\">")
	for _, view := range views {
		label := BuildGateLabel(view)
		output.WriteString("<section class=\"label\">")
		output.WriteString("<div class=\"warning\">" + templateEscape(label.Warning) + "</div>")
		output.WriteString("<h2>" + templateEscape(label.PrimaryLine) + "</h2>")
		output.WriteString("<div>" + templateEscape(label.SecondaryLine) + "</div>")
		output.WriteString("<div>" + templateEscape(label.VerificationLine) + "</div>")
		output.WriteString("<div class=\"id\">" + templateEscape(label.AuthorizationID) + "</div>")
		output.WriteString("</section>")
	}
	output.WriteString("</div></body></html>")
	return output.String(), nil
}

func templateEscape(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;")
	return replacer.Replace(value)
}
