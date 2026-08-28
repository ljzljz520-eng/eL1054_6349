package domain

import (
	"fmt"
	"strings"
	"time"
)

func (a Authorization) IsExpired(today time.Time) bool {
	date, err := time.Parse("2006-01-02", a.AuthorizedDate)
	if err != nil {
		return true
	}
	return date.Before(dayOnly(today))
}

func (a Authorization) IsUsable(today time.Time) bool {
	if a.Status != StatusIssued {
		return false
	}
	date, err := time.Parse("2006-01-02", a.AuthorizedDate)
	if err != nil {
		return false
	}
	return date.Equal(dayOnly(today))
}

func (a Authorization) CanIssue() error {
	if a.Status != StatusDraft {
		return fmt.Errorf("authorization must be draft")
	}
	if a.Revision < 1 {
		return fmt.Errorf("authorization revision is invalid")
	}
	return nil
}

func (a Authorization) CanRevoke() error {
	if a.Status == StatusRevoked {
		return fmt.Errorf("authorization already revoked")
	}
	if a.Status != StatusIssued {
		return fmt.Errorf("only issued authorization can be revoked")
	}
	return nil
}

func (b AuthorizationBundle) Summary() string {
	parts := []string{b.Student.Name, b.Student.ClassName, b.Route.Name, b.Guardian.Name, b.Authorization.AuthorizedDate}
	return strings.Join(parts, " | ")
}

func (b AuthorizationBundle) Matches(filter SearchFilter) bool {
	if filter.Date != "" && b.Authorization.AuthorizedDate != filter.Date {
		return false
	}
	if filter.Status != "" && b.Authorization.Status != filter.Status {
		return false
	}
	if filter.StudentName != "" && !strings.Contains(strings.ToLower(b.Student.Name), strings.ToLower(filter.StudentName)) {
		return false
	}
	if filter.RouteName != "" && !strings.Contains(strings.ToLower(b.Route.Name), strings.ToLower(filter.RouteName)) {
		return false
	}
	return true
}

func dayOnly(value time.Time) time.Time {
	y, m, d := value.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
