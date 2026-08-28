package policy

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"sort"
	"strings"
)

type DisclosureChannel string

const (
	ChannelScreen DisclosureChannel = "screen"
	ChannelPrint  DisclosureChannel = "print"
	ChannelExport DisclosureChannel = "export"
)

type DisclosureField struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
	Masked    bool   `json:"masked"`
	Purpose   string `json:"purpose"`
}

type PrivacyReview struct {
	Channel    DisclosureChannel `json:"channel"`
	Allowed    bool              `json:"allowed"`
	Fields     []DisclosureField `json:"fields"`
	Violations []string          `json:"violations"`
}

func ReviewDisclosure(bundle domain.AuthorizationBundle, channel DisclosureChannel) PrivacyReview {
	fields := []DisclosureField{
		{Name: "student_name", Value: bundle.Student.Name, Sensitive: true, Masked: false, Purpose: "识别乘车学生"},
		{Name: "class_name", Value: bundle.Student.ClassName, Sensitive: false, Masked: false, Purpose: "辅助门卫定位学生"},
		{Name: "school_number", Value: bundle.Student.SchoolNumber, Sensitive: true, Masked: channel == ChannelPrint, Purpose: "消除同名歧义"},
		{Name: "route_name", Value: bundle.Route.Name, Sensitive: false, Masked: false, Purpose: "确认校车线路"},
		{Name: "bus_number", Value: bundle.Route.BusNumber, Sensitive: false, Masked: false, Purpose: "确认车辆"},
		{Name: "guardian_name", Value: bundle.Guardian.Name, Sensitive: true, Masked: false, Purpose: "核验临时接送人"},
		{Name: "relationship", Value: bundle.Guardian.Relationship, Sensitive: true, Masked: false, Purpose: "辅助身份核验"},
		{Name: "document_last_four", Value: "****" + bundle.Guardian.DocumentLastFour, Sensitive: true, Masked: true, Purpose: "最小化证件核验"},
		{Name: "authorized_date", Value: bundle.Authorization.AuthorizedDate, Sensitive: false, Masked: false, Purpose: "确认当日有效"},
		{Name: "teacher_name", Value: bundle.Authorization.TeacherName, Sensitive: false, Masked: false, Purpose: "异常联络"},
	}
	review := PrivacyReview{Channel: channel, Allowed: true, Fields: fields, Violations: make([]string, 0)}
	for _, field := range fields {
		if strings.TrimSpace(field.Value) == "" {
			review.Violations = append(review.Violations, "披露字段为空: "+field.Name)
		}
		if field.Name == "document_last_four" && (!field.Masked || len([]rune(field.Value)) != 8) {
			review.Violations = append(review.Violations, "证件信息未按后四位掩码")
		}
		if channel == ChannelExport && field.Sensitive && !field.Masked && field.Name != "student_name" {
			review.Violations = append(review.Violations, "导出包含未掩码敏感字段: "+field.Name)
		}
	}
	if channel != ChannelScreen && channel != ChannelPrint && channel != ChannelExport {
		review.Violations = append(review.Violations, "未知披露渠道")
	}
	review.Allowed = len(review.Violations) == 0
	sort.Strings(review.Violations)
	return review
}

func (r PrivacyReview) Error() error {
	if r.Allowed {
		return nil
	}
	return fmt.Errorf("privacy review rejected disclosure: %s", strings.Join(r.Violations, "; "))
}

func MaskForChannel(fields []DisclosureField, channel DisclosureChannel) map[string]string {
	result := make(map[string]string, len(fields))
	for _, field := range fields {
		value := field.Value
		if channel == ChannelExport && field.Sensitive && !field.Masked {
			runes := []rune(value)
			if len(runes) > 1 {
				value = string(runes[0]) + strings.Repeat("*", len(runes)-1)
			} else {
				value = "*"
			}
		}
		result[field.Name] = value
	}
	return result
}
