package mopinion

import (
	"context"
	"fmt"
)

// ReportsInterface contains CRUD methods for Report model.
type ReportsInterface interface {
	Get(ctx context.Context, reportID int) (*Report, *Response, error)
	Add(ctx context.Context, report *Report) (*Report, *Response, error)
	Update(ctx context.Context, report *Report) (*Report, *Response, error)
	Delete(ctx context.Context, reportID int, dryRun bool) (*DeleteResponse, *Response, error)
}

// ReportsService implements ReportsInterface.
type ReportsService struct {
	service
}

// Get returns a Report for given report ID.
func (s *ReportsService) Get(ctx context.Context, reportID int) (*Report, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("reports/%d", reportID), nil)
	if err != nil {
		return nil, nil, err
	}
	report := new(Report)
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, resp, err
	}

	return report, resp, nil
}

func (s *ReportsService) validate(report *Report) error {
	if report.Name == "" {
		return fmt.Errorf("report name cannot be empty")
	}
	return nil
}

// Add creates a new Report, if the given report is valid.
func (s *ReportsService) Add(ctx context.Context, report *Report) (*Report, *Response, error) {
	if err := s.validate(report); err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", "reports", report)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, resp, err
	}

	return report, resp, nil
}

// Update edits an existing report model.
func (s *ReportsService) Update(ctx context.Context, report *Report) (*Report, *Response, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("reports/%d", report.ID), report)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, resp, err
	}

	return report, resp, nil
}

// Delete removes the report for given given ID.
func (s *ReportsService) Delete(ctx context.Context, reportID int, dryRun bool) (*DeleteResponse, *Response, error) {
	u := fmt.Sprintf("reports/%d", reportID)
	if dryRun {
		u = fmt.Sprintf("%s?dry-run=true", u)
	}
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}
	deleteRes := new(DeleteResponse)
	resp, err := s.client.Do(ctx, req, deleteRes)
	if err != nil {
		return nil, resp, err
	}

	return deleteRes, resp, nil
}
