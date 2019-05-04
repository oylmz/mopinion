package mopinion

import (
	"context"
	"fmt"
)

// DatasetsInterface contains CRUD methods around Dataset model.
type DatasetsInterface interface {
	Get(ctx context.Context, datasetID int) (*Dataset, *Response, error)
	Add(ctx context.Context, dataset *Dataset) (*Dataset, *Response, error)
	Update(ctx context.Context, dataset *Dataset) (*Dataset, *Response, error)
	Delete(ctx context.Context, datasetID int, dryRun bool) (*DeleteResponse, *Response, error)
}

// DatasetsService implemets DatasetsInterface, holds a reference of service back to mopinion client.
type DatasetsService struct {
	service
}

// Get returns Dataset by given id
func (s *DatasetsService) Get(ctx context.Context, datasetID int) (*Dataset, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("datasets/%d", datasetID), nil)
	if err != nil {
		return nil, nil, err
	}
	dataset := new(Dataset)
	resp, err := s.client.Do(ctx, req, dataset)
	if err != nil {
		return nil, resp, err
	}

	return dataset, resp, nil
}

func (s *DatasetsService) validate(dataset *Dataset) error {
	if dataset.ReportID <= 0 {
		return fmt.Errorf("dataset report id not set")
	}

	if dataset.Name == "" {
		return fmt.Errorf("dataset name cannot be empty")
	}
	return nil
}

// Add creates a new dataset if validation passes.
// Dataset report id or name cannot be empty.
func (s *DatasetsService) Add(ctx context.Context, dataset *Dataset) (*Dataset, *Response, error) {
	if err := s.validate(dataset); err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("POST", "datasets", dataset)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(ctx, req, dataset)
	if err != nil {
		return nil, resp, err
	}

	return dataset, resp, nil
}

// Update simply edits the given dataset.
func (s *DatasetsService) Update(ctx context.Context, dataset *Dataset) (*Dataset, *Response, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("datasets/%d", dataset.ID), dataset)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(ctx, req, dataset)
	if err != nil {
		return nil, resp, err
	}

	return dataset, resp, nil
}

// Delete removes the record with given id. If dry-run is true,
// It rehearses the operation.
func (s *DatasetsService) Delete(ctx context.Context, datasetID int, dryRun bool) (*DeleteResponse, *Response, error) {
	u := fmt.Sprintf("datasets/%d", datasetID)
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
