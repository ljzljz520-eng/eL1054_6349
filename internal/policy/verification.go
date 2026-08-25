package policy

import (
	"fmt"
	"sort"
	"strings"
)

type Verification struct {
	ItemCode string `json:"item_code"`
	Passed   bool   `json:"passed"`
	Note     string `json:"note"`
}

type VerificationResult struct {
	Complete bool              `json:"complete"`
	Missing  []ChecklistItem   `json:"missing"`
	Warnings []string          `json:"warnings"`
	Evidence map[string]string `json:"evidence"`
}

func Verify(checklist []ChecklistItem, answers []Verification) VerificationResult {
	answerMap := make(map[string]Verification, len(answers))
	warnings := make([]string, 0)
	for _, answer := range answers {
		if _, exists := answerMap[answer.ItemCode]; exists {
			warnings = append(warnings, "核验项重复: "+answer.ItemCode)
		}
		answerMap[answer.ItemCode] = answer
	}
	result := VerificationResult{Complete: true, Missing: make([]ChecklistItem, 0), Warnings: warnings, Evidence: map[string]string{}}
	for _, item := range checklist {
		answer, exists := answerMap[item.Code]
		if !exists || !answer.Passed {
			if item.Required {
				result.Complete = false
				result.Missing = append(result.Missing, item)
			}
			continue
		}
		result.Evidence[item.Code] = strings.TrimSpace(answer.Note)
		if item.Required && strings.TrimSpace(answer.Note) == "" {
			result.Warnings = append(result.Warnings, "核验项未填写备注: "+item.Label)
		}
	}
	sort.Slice(result.Missing, func(i, j int) bool { return result.Missing[i].Code < result.Missing[j].Code })
	sort.Strings(result.Warnings)
	return result
}

func (r VerificationResult) Error() error {
	if r.Complete {
		return nil
	}
	names := make([]string, 0, len(r.Missing))
	for _, item := range r.Missing {
		names = append(names, item.Label)
	}
	return fmt.Errorf("required gate checks incomplete: %s", strings.Join(names, ", "))
}

func DefaultAnswers(checklist []ChecklistItem, note string) []Verification {
	result := make([]Verification, 0, len(checklist))
	for _, item := range checklist {
		if item.Required {
			result = append(result, Verification{ItemCode: item.Code, Passed: true, Note: note})
		}
	}
	return result
}
