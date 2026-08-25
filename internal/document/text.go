package document

import (
	"bytes"
	"fmt"
	"strings"
)

type TextRenderer struct{}

func (TextRenderer) Render(view View) ([]byte, error) {
	if view.AuthorizationID == "" {
		return nil, fmt.Errorf("authorization id is required")
	}
	var output bytes.Buffer
	lines := []string{
		center(view.Title, 52),
		center("授权编号: "+view.AuthorizationID, 52),
		strings.Repeat("=", 52),
		"[" + view.NoticeTitle + "] " + view.NoticeMessage,
		strings.Repeat("-", 52),
		pair("学生姓名", view.StudentName),
		pair("班级", view.ClassName),
		pair("学号", view.SchoolNumber),
		pair("授权日期", view.AuthorizedDate),
		pair("线路", view.RouteName),
		pair("车辆", view.BusNumber),
		pair("接送点", view.PickupPoint),
		pair("核验门岗", view.GuardPost),
		pair("临时接送人", view.GuardianName),
		pair("关系", view.Relationship),
		pair("证件", view.MaskedDocument),
		pair("班主任", view.TeacherName),
		pair("原因", view.Reason),
		strings.Repeat("-", 52),
		"门卫核验清单:",
	}
	for index, item := range view.VerificationItems {
		lines = append(lines, fmt.Sprintf("%d. [ ] %s", index+1, item))
	}
	lines = append(lines, "", "班主任签字: ______________", "门卫签字: ______________", "", "生成时间: "+view.GeneratedAt)
	for _, line := range lines {
		output.WriteString(line)
		output.WriteByte('\n')
	}
	return output.Bytes(), nil
}

func pair(label, value string) string { return fmt.Sprintf("%-12s %s", label+":", value) }

func center(value string, width int) string {
	spaces := width - len([]rune(value))
	if spaces <= 0 {
		return value
	}
	return strings.Repeat(" ", spaces/2) + value
}
