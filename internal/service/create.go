package service

import (
	"fmt"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/validation"
)

func (s *Service) CreateDraft(req domain.CreateRequest) (domain.AuthorizationBundle, error) {
	req = validation.Normalize(req)
	if err := validation.CreateRequest(req); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	now := s.clock.Now().UTC()
	studentID := entityID("student", req.Student.SchoolNumber)
	routeID := entityID("route", req.Route.Name+req.Route.BusNumber)
	guardianID := entityID("guardian", req.Student.SchoolNumber+req.Guardian.Name+req.Guardian.DocumentLastFour)
	authID := entityID("auth", req.Student.SchoolNumber+req.AuthorizedDate+req.Route.BusNumber)
	req.Student.ID = studentID
	req.Route.ID = routeID
	req.Guardian.ID = guardianID
	bundle := domain.AuthorizationBundle{
		Student:  req.Student,
		Route:    req.Route,
		Guardian: req.Guardian,
		Authorization: domain.Authorization{
			ID:             authID,
			StudentID:      studentID,
			RouteID:        routeID,
			GuardianID:     guardianID,
			AuthorizedDate: req.AuthorizedDate,
			TeacherName:    req.TeacherName,
			Reason:         req.Reason,
			Status:         domain.StatusDraft,
			CreatedAt:      now,
			Revision:       1,
		},
	}
	audit := domain.AuditEvent{ID: eventID(authID, "created", 1), AuthorizationID: authID, Action: "created", Actor: req.TeacherName, OccurredAt: now, Details: bundle.Summary()}
	if err := s.store.CreateBundle(bundle, audit); err != nil {
		return domain.AuthorizationBundle{}, fmt.Errorf("save draft: %w", err)
	}
	return bundle, nil
}

func (s *Service) UpdateDraft(id string, req domain.CreateRequest) (domain.AuthorizationBundle, error) {
	req = validation.Normalize(req)
	if err := validation.CreateRequest(req); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	bundle, err := s.store.Bundle(id)
	if err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if bundle.Authorization.Status != domain.StatusDraft {
		return domain.AuthorizationBundle{}, fmt.Errorf("only a draft can be updated")
	}
	bundle.Student.Name = req.Student.Name
	bundle.Student.ClassName = req.Student.ClassName
	bundle.Student.SchoolNumber = req.Student.SchoolNumber
	bundle.Route.Name = req.Route.Name
	bundle.Route.BusNumber = req.Route.BusNumber
	bundle.Route.PickupPoint = req.Route.PickupPoint
	bundle.Route.GuardPost = req.Route.GuardPost
	bundle.Guardian.Name = req.Guardian.Name
	bundle.Guardian.Relationship = req.Guardian.Relationship
	bundle.Guardian.DocumentLastFour = req.Guardian.DocumentLastFour
	bundle.Authorization.AuthorizedDate = req.AuthorizedDate
	bundle.Authorization.TeacherName = req.TeacherName
	bundle.Authorization.Reason = req.Reason
	bundle.Authorization.Revision++
	if err := s.store.SaveStudent(bundle.Student); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if err := s.store.SaveRoute(bundle.Route); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	if err := s.store.SaveGuardian(bundle.Guardian); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	audit := domain.AuditEvent{ID: eventID(id, "updated", bundle.Authorization.Revision), AuthorizationID: id, Action: "updated", Actor: req.TeacherName, OccurredAt: s.clock.Now().UTC(), Details: bundle.Summary()}
	if err := s.store.UpdateAuthorization(bundle.Authorization, audit); err != nil {
		return domain.AuthorizationBundle{}, err
	}
	return bundle, nil
}
