package document

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/validation"
	"time"
)

type NoticeLevel string

const (
	NoticeValid   NoticeLevel = "valid"
	NoticeExpired NoticeLevel = "expired"
	NoticeRevoked NoticeLevel = "revoked"
	NoticeDraft   NoticeLevel = "draft"
)

type View struct {
	Title             string
	AuthorizationID   string
	StudentName       string
	ClassName         string
	SchoolNumber      string
	RouteName         string
	BusNumber         string
	PickupPoint       string
	GuardPost         string
	GuardianName      string
	Relationship      string
	MaskedDocument    string
	AuthorizedDate    string
	TeacherName       string
	Reason            string
	Status            string
	NoticeLevel       NoticeLevel
	NoticeTitle       string
	NoticeMessage     string
	GeneratedAt       string
	VerificationItems []string
}

func BuildView(bundle domain.AuthorizationBundle, today time.Time) View {
	view := View{
		Title:             "校车临时接送授权单",
		AuthorizationID:   bundle.Authorization.ID,
		StudentName:       bundle.Student.Name,
		ClassName:         bundle.Student.ClassName,
		SchoolNumber:      bundle.Student.SchoolNumber,
		RouteName:         bundle.Route.Name,
		BusNumber:         bundle.Route.BusNumber,
		PickupPoint:       bundle.Route.PickupPoint,
		GuardPost:         bundle.Route.GuardPost,
		GuardianName:      bundle.Guardian.Name,
		Relationship:      bundle.Guardian.Relationship,
		MaskedDocument:    validation.MaskDocument(bundle.Guardian.DocumentLastFour),
		AuthorizedDate:    bundle.Authorization.AuthorizedDate,
		TeacherName:       bundle.Authorization.TeacherName,
		Reason:            bundle.Authorization.Reason,
		Status:            string(bundle.Authorization.Status),
		GeneratedAt:       today.UTC().Format("2006-01-02 15:04 UTC"),
		VerificationItems: []string{"核对学生姓名与班级", "核对接送人姓名与关系", "核对证件后四位", "确认授权日期为当天", "登记放行时间"},
	}
	setNotice(&view, bundle.Authorization, today)
	return view
}

func setNotice(view *View, auth domain.Authorization, today time.Time) {
	switch {
	case auth.Status == domain.StatusRevoked:
		view.NoticeLevel = NoticeRevoked
		view.NoticeTitle = "授权已撤销"
		view.NoticeMessage = "本授权不可用于放行，请联系班主任核实。"
	case auth.IsExpired(today):
		view.NoticeLevel = NoticeExpired
		view.NoticeTitle = "授权已过期"
		view.NoticeMessage = fmt.Sprintf("授权日期为 %s，不得用于今日放行。", auth.AuthorizedDate)
	case auth.Status == domain.StatusDraft:
		view.NoticeLevel = NoticeDraft
		view.NoticeTitle = "尚未签发"
		view.NoticeMessage = "本预览仅供班主任核对，签发后方可打印。"
	default:
		view.NoticeLevel = NoticeValid
		view.NoticeTitle = "仅限当日使用"
		view.NoticeMessage = "门卫核对证件后四位并登记后放行。"
	}
}
