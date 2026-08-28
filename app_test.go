package schoolbusauth

import (
	"path/filepath"
	"testing"
	"time"
)

func TestBusinessChain12(t *testing.T) {
	date := time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)
	application, err := Open(Config{DatabasePath: filepath.Join(t.TempDir(), "chain.db"), BusinessDate: date})
	if err != nil {
		t.Fatal(err)
	}
	defer application.Close()
	draft, err := application.Service.CreateDraft(DemoRequest("2026-09-01"))
	if err != nil {
		t.Fatal(err)
	}
	issued, err := application.Service.Issue(draft.Authorization.ID, "王老师")
	if err != nil {
		t.Fatal(err)
	}
	preview, err := application.Service.Preview(issued.Authorization.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.HTML) == 0 {
		t.Fatal("empty preview")
	}
	if err := application.Service.Complete(issued.Authorization.ID); err != nil {
		t.Fatal(err)
	}
}
