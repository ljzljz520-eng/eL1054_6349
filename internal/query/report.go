package query

import (
	"encoding/csv"
	"fmt"
	"io"
	"schoolbusauth/internal/domain"
	"sort"
	"strconv"
)

type DailySummary struct {
	Date        string         `json:"date"`
	Total       int            `json:"total"`
	Drafts      int            `json:"drafts"`
	Issued      int            `json:"issued"`
	Revoked     int            `json:"revoked"`
	PrintCopies int            `json:"print_copies"`
	RouteCounts map[string]int `json:"route_counts"`
}

func (s *Service) Summary(date string) (DailySummary, error) {
	values, err := s.Search(domain.SearchFilter{Date: date})
	if err != nil {
		return DailySummary{}, err
	}
	result := DailySummary{Date: date, Total: len(values), RouteCounts: map[string]int{}}
	for _, value := range values {
		switch value.Authorization.Status {
		case domain.StatusDraft:
			result.Drafts++
		case domain.StatusIssued:
			result.Issued++
		case domain.StatusRevoked:
			result.Revoked++
		}
		result.RouteCounts[value.Route.Name]++
		prints, err := s.store.PrintRecords(value.Authorization.ID)
		if err != nil {
			return DailySummary{}, err
		}
		for _, print := range prints {
			result.PrintCopies += print.Copies
		}
	}
	return result, nil
}

func (s *Service) ExportCSV(writer io.Writer, filter domain.SearchFilter) error {
	values, err := s.Search(filter)
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(writer)
	rows := [][]string{{"授权编号", "日期", "学生", "班级", "线路", "车辆", "接送人", "关系", "证件后四位", "状态", "班主任"}}
	for _, value := range values {
		rows = append(rows, []string{value.Authorization.ID, value.Authorization.AuthorizedDate, value.Student.Name, value.Student.ClassName, value.Route.Name, value.Route.BusNumber, value.Guardian.Name, value.Guardian.Relationship, value.Guardian.DocumentLastFour, string(value.Authorization.Status), value.Authorization.TeacherName})
	}
	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("write report row: %w", err)
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

func RouteRanking(summary DailySummary) []string {
	type item struct {
		name  string
		count int
	}
	items := make([]item, 0, len(summary.RouteCounts))
	for name, count := range summary.RouteCounts {
		items = append(items, item{name, count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].count == items[j].count {
			return items[i].name < items[j].name
		}
		return items[i].count > items[j].count
	})
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.name+":"+strconv.Itoa(item.count))
	}
	return result
}
