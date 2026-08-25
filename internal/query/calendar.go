package query

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"sort"
	"time"
)

type CalendarDay struct {
	Date     string   `json:"date"`
	Weekday  string   `json:"weekday"`
	Total    int      `json:"total"`
	Issued   int      `json:"issued"`
	Draft    int      `json:"draft"`
	Revoked  int      `json:"revoked"`
	Students []string `json:"students"`
	Routes   []string `json:"routes"`
}

type Calendar struct {
	Start string        `json:"start"`
	End   string        `json:"end"`
	Days  []CalendarDay `json:"days"`
}

func (s *Service) Calendar(start, end time.Time) (Calendar, error) {
	if end.Before(start) {
		return Calendar{}, fmt.Errorf("calendar end must not precede start")
	}
	if end.Sub(start) > 31*24*time.Hour {
		return Calendar{}, fmt.Errorf("calendar range must not exceed 31 days")
	}
	values, err := s.store.Bundles()
	if err != nil {
		return Calendar{}, err
	}
	byDate := make(map[string][]domain.AuthorizationBundle)
	for _, value := range values {
		byDate[value.Authorization.AuthorizedDate] = append(byDate[value.Authorization.AuthorizedDate], value)
	}
	result := Calendar{Start: start.Format("2006-01-02"), End: end.Format("2006-01-02"), Days: make([]CalendarDay, 0)}
	for cursor := day(start); !cursor.After(day(end)); cursor = cursor.AddDate(0, 0, 1) {
		key := cursor.Format("2006-01-02")
		calendarDay := CalendarDay{Date: key, Weekday: cursor.Weekday().String(), Students: make([]string, 0), Routes: make([]string, 0)}
		routeSet := map[string]struct{}{}
		for _, value := range byDate[key] {
			calendarDay.Total++
			calendarDay.Students = append(calendarDay.Students, value.Student.Name)
			routeSet[value.Route.Name] = struct{}{}
			switch value.Authorization.Status {
			case domain.StatusIssued:
				calendarDay.Issued++
			case domain.StatusDraft:
				calendarDay.Draft++
			case domain.StatusRevoked:
				calendarDay.Revoked++
			}
		}
		for route := range routeSet {
			calendarDay.Routes = append(calendarDay.Routes, route)
		}
		sort.Strings(calendarDay.Students)
		sort.Strings(calendarDay.Routes)
		result.Days = append(result.Days, calendarDay)
	}
	return result, nil
}

func day(value time.Time) time.Time {
	y, m, d := value.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (c Calendar) BusyDays(threshold int) []CalendarDay {
	result := make([]CalendarDay, 0)
	for _, value := range c.Days {
		if value.Total >= threshold {
			result = append(result, value)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Total == result[j].Total {
			return result[i].Date < result[j].Date
		}
		return result[i].Total > result[j].Total
	})
	return result
}
