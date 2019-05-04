package mopinion

import (
	"context"
	"fmt"
)

// FieldsInterface has two methods returning fields for a specific dataset or report.
type FieldsInterface interface {
	GetByDataset(ctx context.Context, datasetID int) (*Fields, *Response, error)
	GetByReport(ctx context.Context, reportID int) (*Fields, *Response, error)
}

// FieldsService implements FieldsInterface
type FieldsService struct {
	service
}

func (s *FieldsService) get(ctx context.Context, url string) (*Fields, *Response, error) {
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	fields := new(Fields)
	resp, err := s.client.Do(ctx, req, fields)
	if err != nil {
		return nil, resp, err
	}

	return fields, resp, nil
}

// GetByDataset returns the fields of given dataset id
func (s *FieldsService) GetByDataset(ctx context.Context, datasetID int) (*Fields, *Response, error) {
	return s.get(ctx, fmt.Sprintf("datasets/%d/fields", datasetID))
}

// GetByReport returns the fields of given report id
func (s *FieldsService) GetByReport(ctx context.Context, reportID int) (*Fields, *Response, error) {
	return s.get(ctx, fmt.Sprintf("reports/%d/fields", reportID))
}
