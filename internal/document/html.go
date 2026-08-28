package document

import (
	"bytes"
	"fmt"
	"html/template"
)

type HTMLRenderer struct{ template *template.Template }

func NewHTMLRenderer() (*HTMLRenderer, error) {
	tmpl, err := template.New("authorization").Parse(pageTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse document template: %w", err)
	}
	return &HTMLRenderer{template: tmpl}, nil
}

func (r *HTMLRenderer) Render(view View) ([]byte, error) {
	var output bytes.Buffer
	if err := r.template.Execute(&output, view); err != nil {
		return nil, fmt.Errorf("render document: %w", err)
	}
	return output.Bytes(), nil
}

const pageTemplate = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}} - {{.AuthorizationID}}</title>
<style>
body { margin: 0; color: #17202a; background: #eef1f3; font-family: sans-serif; }
.sheet { width: 180mm; min-height: 250mm; margin: 16px auto; padding: 14mm; background: white; box-sizing: border-box; }
h1 { margin: 0 0 4px; text-align: center; font-size: 26px; }
.serial { text-align: center; color: #566573; margin-bottom: 18px; }
.notice { border: 3px solid #2471a3; padding: 12px; margin: 12px 0 20px; }
.notice.expired, .notice.revoked { border-color: #c0392b; background: #fdecea; color: #922b21; }
.notice.draft { border-color: #b9770e; background: #fef5e7; }
.notice h2 { margin: 0 0 6px; font-size: 20px; }
.grid { display: grid; grid-template-columns: 1fr 1fr; border-top: 1px solid #85929e; border-left: 1px solid #85929e; }
.field { padding: 10px; border-right: 1px solid #85929e; border-bottom: 1px solid #85929e; min-height: 42px; }
.field.wide { grid-column: 1 / -1; }
.label { display: block; color: #566573; font-size: 12px; margin-bottom: 5px; }
.value { font-size: 17px; font-weight: 600; }
.checklist { margin-top: 22px; }
.checklist li { margin: 9px 0; }
.signatures { display: grid; grid-template-columns: 1fr 1fr; gap: 28px; margin-top: 34px; }
.signature { border-bottom: 1px solid #17202a; height: 34px; }
.footer { margin-top: 28px; font-size: 11px; color: #566573; text-align: center; }
.actions { text-align: center; margin: 12px; }
button { padding: 9px 18px; }
@media print { body { background: white; } .sheet { margin: 0; width: auto; min-height: auto; } .actions { display: none; } }
</style>
</head>
<body>
<div class="actions"><button onclick="window.print()">打印授权单</button></div>
<main class="sheet">
<h1>{{.Title}}</h1>
<div class="serial">授权编号：{{.AuthorizationID}}</div>
<section class="notice {{.NoticeLevel}}"><h2>{{.NoticeTitle}}</h2><div>{{.NoticeMessage}}</div></section>
<section class="grid">
<div class="field"><span class="label">学生姓名</span><span class="value">{{.StudentName}}</span></div>
<div class="field"><span class="label">班级</span><span class="value">{{.ClassName}}</span></div>
<div class="field"><span class="label">学号</span><span class="value">{{.SchoolNumber}}</span></div>
<div class="field"><span class="label">授权日期</span><span class="value">{{.AuthorizedDate}}</span></div>
<div class="field"><span class="label">校车线路</span><span class="value">{{.RouteName}}</span></div>
<div class="field"><span class="label">车辆编号</span><span class="value">{{.BusNumber}}</span></div>
<div class="field"><span class="label">接送点</span><span class="value">{{.PickupPoint}}</span></div>
<div class="field"><span class="label">核验门岗</span><span class="value">{{.GuardPost}}</span></div>
<div class="field"><span class="label">临时接送人</span><span class="value">{{.GuardianName}}</span></div>
<div class="field"><span class="label">与学生关系</span><span class="value">{{.Relationship}}</span></div>
<div class="field wide"><span class="label">证件号码（仅显示后四位）</span><span class="value">{{.MaskedDocument}}</span></div>
<div class="field"><span class="label">班主任</span><span class="value">{{.TeacherName}}</span></div>
<div class="field"><span class="label">当前状态</span><span class="value">{{.Status}}</span></div>
<div class="field wide"><span class="label">临时变更原因</span><span class="value">{{.Reason}}</span></div>
</section>
<section class="checklist"><h2>门卫核验清单</h2><ol>{{range .VerificationItems}}<li>□ {{.}}</li>{{end}}</ol></section>
<section class="signatures"><div><div class="signature"></div><div>班主任签字</div></div><div><div class="signature"></div><div>门卫核验签字</div></div></section>
<div class="footer">生成时间：{{.GeneratedAt}} · 本文书只保存必要的证件后四位</div>
</main>
</body>
</html>`
