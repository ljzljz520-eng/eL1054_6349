package operations

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"sort"
	"strings"
	"time"
)

type GateShift struct {
	GuardPost string `json:"guard_post"`
	Start     string `json:"start"`
	End       string `json:"end"`
	GuardName string `json:"guard_name"`
}

type ReleaseWindow struct {
	GuardPost     string   `json:"guard_post"`
	Route         string   `json:"route"`
	Students      int      `json:"students"`
	WindowStart   string   `json:"window_start"`
	WindowEnd     string   `json:"window_end"`
	AssignedGuard string   `json:"assigned_guard"`
	Warnings      []string `json:"warnings"`
}

type GateSchedule struct {
	Date       string          `json:"date"`
	Windows    []ReleaseWindow `json:"windows"`
	Unassigned []string        `json:"unassigned"`
}

func BuildGateSchedule(bundles []domain.AuthorizationBundle, shifts []GateShift, date time.Time) GateSchedule {
	groups := map[string][]domain.AuthorizationBundle{}
	for _, bundle := range bundles {
		if !bundle.Authorization.IsUsable(date) {
			continue
		}
		key := bundle.Route.GuardPost + "\x00" + bundle.Route.Name
		groups[key] = append(groups[key], bundle)
	}
	result := GateSchedule{Date: date.Format("2006-01-02"), Windows: make([]ReleaseWindow, 0), Unassigned: make([]string, 0)}
	for key, values := range groups {
		parts := strings.Split(key, "\x00")
		window := ReleaseWindow{GuardPost: parts[0], Route: parts[1], Students: len(values), WindowStart: "15:30", WindowEnd: "17:30", Warnings: make([]string, 0)}
		shift, ok := coveringShift(shifts, parts[0], window.WindowStart, window.WindowEnd)
		if ok {
			window.AssignedGuard = shift.GuardName
		} else {
			window.Warnings = append(window.Warnings, "放行时段没有完整覆盖的门卫班次")
			result.Unassigned = append(result.Unassigned, parts[0]+" / "+parts[1])
		}
		if len(values) > 20 {
			window.Warnings = append(window.Warnings, "临时授权人数较多，建议增派核验人员")
		}
		if parts[0] == "" {
			window.Warnings = append(window.Warnings, "线路未指定门岗")
		}
		result.Windows = append(result.Windows, window)
	}
	sort.Slice(result.Windows, func(i, j int) bool {
		if result.Windows[i].GuardPost == result.Windows[j].GuardPost {
			return result.Windows[i].Route < result.Windows[j].Route
		}
		return result.Windows[i].GuardPost < result.Windows[j].GuardPost
	})
	sort.Strings(result.Unassigned)
	return result
}

func coveringShift(shifts []GateShift, post, start, end string) (GateShift, bool) {
	for _, shift := range shifts {
		if shift.GuardPost != post {
			continue
		}
		if shift.Start <= start && shift.End >= end && strings.TrimSpace(shift.GuardName) != "" {
			return shift, true
		}
	}
	return GateShift{}, false
}

func (s GateSchedule) Validate() error {
	if s.Date == "" {
		return fmt.Errorf("schedule date is required")
	}
	if len(s.Unassigned) > 0 {
		return fmt.Errorf("unassigned release windows: %s", strings.Join(s.Unassigned, ", "))
	}
	for _, window := range s.Windows {
		if window.Students < 1 {
			return fmt.Errorf("release window has no students")
		}
		if window.AssignedGuard == "" {
			return fmt.Errorf("release window has no guard")
		}
	}
	return nil
}
