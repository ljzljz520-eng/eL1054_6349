package schoolbusauth

import (
	"fmt"
	"net/http"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/httpapi"
	"schoolbusauth/internal/service"
	"schoolbusauth/internal/store"
	"time"
)

type App struct {
	Store   *store.Store
	Service *service.Service
	Handler http.Handler
}

type Config struct {
	DatabasePath string
	BusinessDate time.Time
	OnComplete   service.CompletionCallback
}

func Open(config Config) (*App, error) {
	if config.DatabasePath == "" {
		return nil, fmt.Errorf("database path is required")
	}
	if config.BusinessDate.IsZero() {
		return nil, fmt.Errorf("business date is required")
	}
	repository, err := store.Open(config.DatabasePath)
	if err != nil {
		return nil, err
	}
	serviceLayer, err := service.New(repository, service.FixedClock{Value: config.BusinessDate}, config.OnComplete)
	if err != nil {
		repository.Close()
		return nil, err
	}
	return &App{Store: repository, Service: serviceLayer, Handler: httpapi.New(serviceLayer).Handler()}, nil
}

func (a *App) Close() error {
	if a == nil || a.Store == nil {
		return nil
	}
	return a.Store.Close()
}

func DemoRequest(date string) domain.CreateRequest {
	return domain.CreateRequest{
		Student:        domain.Student{Name: "李明", ClassName: "三年级二班", SchoolNumber: "S-2026-018"},
		Route:          domain.Route{Name: "东湖一号线", BusNumber: "BUS-08", PickupPoint: "东湖社区南门", GuardPost: "学校东门"},
		Guardian:       domain.Guardian{Name: "李华", Relationship: "父亲", DocumentLastFour: "4832"},
		AuthorizedDate: date,
		TeacherName:    "王老师",
		Reason:         "家长临时委托接送",
	}
}
