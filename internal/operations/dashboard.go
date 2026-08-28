package operations

import (
	"schoolbusauth/internal/domain"
	"sort"
	"time"
)

type Dashboard struct {
	Date      string          `json:"date"`
	Total     int             `json:"total"`
	Usable    int             `json:"usable"`
	Expired   int             `json:"expired"`
	Draft     int             `json:"draft"`
	Revoked   int             `json:"revoked"`
	ByRoute   []RouteVolume   `json:"by_route"`
	Attention []AttentionItem `json:"attention"`
}

type RouteVolume struct {
	Route      string   `json:"route"`
	Count      int      `json:"count"`
	BusNumbers []string `json:"bus_numbers"`
}
type AttentionItem struct {
	AuthorizationID string `json:"authorization_id"`
	Student         string `json:"student"`
	Reason          string `json:"reason"`
	Priority        int    `json:"priority"`
}

func BuildDashboard(bundles []domain.AuthorizationBundle, today time.Time) Dashboard {
	result := Dashboard{Date: today.Format("2006-01-02"), Total: len(bundles), ByRoute: make([]RouteVolume, 0), Attention: make([]AttentionItem, 0)}
	routes := map[string]*RouteVolume{}
	for _, bundle := range bundles {
		switch {
		case bundle.Authorization.Status == domain.StatusRevoked:
			result.Revoked++
			result.Attention = append(result.Attention, attention(bundle, "授权已撤销", 90))
		case bundle.Authorization.IsExpired(today):
			result.Expired++
			result.Attention = append(result.Attention, attention(bundle, "授权已过期", 80))
		case bundle.Authorization.Status == domain.StatusDraft:
			result.Draft++
			result.Attention = append(result.Attention, attention(bundle, "等待班主任签发", 50))
		case bundle.Authorization.IsUsable(today):
			result.Usable++
		}
		volume := routes[bundle.Route.Name]
		if volume == nil {
			volume = &RouteVolume{Route: bundle.Route.Name}
			routes[bundle.Route.Name] = volume
		}
		volume.Count++
		if !contains(volume.BusNumbers, bundle.Route.BusNumber) {
			volume.BusNumbers = append(volume.BusNumbers, bundle.Route.BusNumber)
		}
	}
	for _, volume := range routes {
		sort.Strings(volume.BusNumbers)
		result.ByRoute = append(result.ByRoute, *volume)
	}
	sort.Slice(result.ByRoute, func(i, j int) bool {
		if result.ByRoute[i].Count == result.ByRoute[j].Count {
			return result.ByRoute[i].Route < result.ByRoute[j].Route
		}
		return result.ByRoute[i].Count > result.ByRoute[j].Count
	})
	sort.Slice(result.Attention, func(i, j int) bool {
		if result.Attention[i].Priority == result.Attention[j].Priority {
			return result.Attention[i].Student < result.Attention[j].Student
		}
		return result.Attention[i].Priority > result.Attention[j].Priority
	})
	return result
}

func attention(bundle domain.AuthorizationBundle, reason string, priority int) AttentionItem {
	return AttentionItem{AuthorizationID: bundle.Authorization.ID, Student: bundle.Student.Name, Reason: reason, Priority: priority}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
