package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

// DeploymentsInterface contains CRUD methods for Deployments model.
type DeploymentsInterface interface {
	Get(ctx context.Context) (*Deployments, *Response, error)
	Add(ctx context.Context, deployment *Deployment) (*Deployments, *Response, error)
	Delete(ctx context.Context, deploymentKey string, dryRun bool) (*DeleteResponse, *Response, error)
}

// DeploymentsService implements DeploymentsInterface.
type DeploymentsService struct {
	service
}

type jdata map[string]*json.RawMessage

// UnmarshalJSON implements Unmarshaler interface.
func (d *Deployments) UnmarshalJSON(buf []byte) (err error) {
	var (
		meta *json.RawMessage
		ok   bool
	)
	jdata := new(jdata)

	if err = json.Unmarshal(buf, jdata); err != nil {
		return
	}
	if meta, ok = (*jdata)["_meta"]; ok {
		if err = json.Unmarshal(*meta, &d.Meta); err != nil {
			return
		}
	}
	for k, v := range *jdata {
		if k == "_meta" {
			continue
		}
		deployment := &Deployment{}
		if err = json.Unmarshal(*v, deployment); err != nil {
			return err
		}
		d.Deployments = append(d.Deployments, *deployment)
		sort.SliceStable(d.Deployments, func(i, j int) bool {
			return d.Deployments[i].Key < d.Deployments[j].Key
		})
	}
	return nil
}

// Get returns Deployments
func (s *DeploymentsService) Get(ctx context.Context) (*Deployments, *Response, error) {
	req, err := s.client.NewRequest("GET", "deployments", nil)
	if err != nil {
		return nil, nil, err
	}
	deployments := new(Deployments)
	resp, err := s.client.Do(ctx, req, deployments)
	if err != nil {
		return nil, resp, err
	}

	return deployments, resp, nil
}

func (s *DeploymentsService) validate(deployment *Deployment) error {
	if deployment.Name == "" {
		return fmt.Errorf("deployment name cannot be empty")
	}
	return nil
}

// Add creates a new Deployment, if the given deployment is valid.
func (s *DeploymentsService) Add(ctx context.Context, deployment *Deployment) (*Deployments, *Response, error) {
	if err := s.validate(deployment); err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("POST", "deployments", deployment)
	if err != nil {
		return nil, nil, err
	}
	deployments := new(Deployments)
	resp, err := s.client.Do(ctx, req, deployments)
	if err != nil {
		return nil, resp, err
	}

	return deployments, resp, nil
}

// Delete removes the deployment for given given ID.
func (s *DeploymentsService) Delete(ctx context.Context, deploymentKey string, dryRun bool) (*DeleteResponse, *Response, error) {
	u := fmt.Sprintf("deployments/%s", deploymentKey)
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
