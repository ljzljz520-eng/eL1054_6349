package httpapi

import (
	"html/template"
	"net/http"
	"schoolbusauth/internal/domain"
)

var homeTemplate = template.Must(template.New("home").Parse(`<!doctype html>
<html lang="zh-CN"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>校车临时接送授权单</title><style>
body{margin:0;font-family:sans-serif;color:#17202a;background:#f4f6f6}header{background:#154360;color:white;padding:18px 5vw}main{max-width:1100px;margin:22px auto;padding:0 18px}form{background:white;border:1px solid #ccd1d1;padding:20px}.grid{display:grid;grid-template-columns:repeat(3,1fr);gap:14px}label{display:flex;flex-direction:column;gap:6px;font-weight:600}input,textarea{padding:9px;border:1px solid #85929e}textarea{min-height:70px}.wide{grid-column:1/-1}button{padding:10px 20px;background:#117864;color:white;border:0}.list{margin-top:24px;background:white;border-collapse:collapse;width:100%}th,td{padding:10px;border:1px solid #ccd1d1;text-align:left}@media(max-width:700px){.grid{grid-template-columns:1fr}}
</style></head><body><header><h1>校车临时接送授权单</h1></header><main>
<form action="/authorizations" method="post"><div class="grid">
<label>学生姓名<input name="student_name" required></label><label>班级<input name="class_name" required></label><label>学号<input name="school_number" required></label>
<label>线路名称<input name="route_name" required></label><label>车辆编号<input name="bus_number" required></label><label>临时接送点<input name="pickup_point" required></label>
<label>核验门岗<input name="guard_post" required></label><label>临时接送人<input name="guardian_name" required></label><label>关系<input name="relationship" required></label>
<label>证件后四位<input name="document_last_four" minlength="4" maxlength="4" required></label><label>授权日期<input type="date" name="authorized_date" required></label><label>班主任<input name="teacher_name" required></label>
<label class="wide">临时变更原因<textarea name="reason"></textarea></label><div class="wide"><button type="submit">生成授权草稿</button></div>
</div></form>
<table class="list"><thead><tr><th>日期</th><th>学生</th><th>线路</th><th>接送人</th><th>状态</th><th>操作</th></tr></thead><tbody>{{range .}}<tr><td>{{.Authorization.AuthorizedDate}}</td><td>{{.Student.Name}}</td><td>{{.Route.Name}}</td><td>{{.Guardian.Name}}</td><td>{{.Authorization.Status}}</td><td><a href="/authorizations/{{.Authorization.ID}}">查看</a></td></tr>{{else}}<tr><td colspan="6">暂无授权记录</td></tr>{{end}}</tbody></table>
</main></body></html>`))

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	values, err := s.query.Search(domain.SearchFilter{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := homeTemplate.Execute(w, values); err != nil {
		return
	}
}

func (s *Server) detail(w http.ResponseWriter, r *http.Request) {
	value, err := s.query.Detail(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
