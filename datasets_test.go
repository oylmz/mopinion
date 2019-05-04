package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDatasetGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	datasetStr := `{
		"name": "dataset name",
		"report_id": 1,
		"description": "dataset description",
		"id": 1,
		"data_source": "form",
		"_meta": {
			"code": 200,
			"message": "OK",
			"has_more": false,
			"previous": false,
			"next": false,
			"count": 1,
			"total": 1
		}
	}`

	mux.HandleFunc("/datasets/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, datasetStr)
	})

	client.Token.Get(context.Background())
	dataset, _, err := client.Datasets.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("datasets API should not return an error: %s", err)
	}

	var expectedDataset = &Dataset{}
	err = json.Unmarshal([]byte(datasetStr), expectedDataset)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDataset, dataset) {
		t.Errorf("expected dataset: %+v but got: %+v", expectedDataset, dataset)
	}
}

func TestDatasetAdd(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	datasetStr := `{
		"name": "dataset name",
		"report_id": 1,
		"description": "description",
		"id": 1,
		"data_source": "import"
	}`

	mux.HandleFunc("/datasets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, datasetStr)
	})

	client.Token.Get(context.Background())
	newDataset := &Dataset{
		Name:        "dataset name",
		Description: "description",
		ReportID:    1,
	}
	dataset, _, err := client.Datasets.Add(context.Background(), newDataset)
	if err != nil {
		t.Errorf("datasets API should not return an error: %s", err)
	}

	var expectedDataset = &Dataset{}
	err = json.Unmarshal([]byte(datasetStr), expectedDataset)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDataset, dataset) {
		t.Errorf("expected dataset: %+v but got: %+v", expectedDataset, dataset)
	}
}

func TestDatasetUpdate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	datasetStr := `{
		"name": "dataset name",
		"report_id": 1,
		"description": "description",
		"id": 1,
		"data_source": "import"
	}`

	mux.HandleFunc("/datasets/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, datasetStr)
	})

	client.Token.Get(context.Background())
	newDataset := &Dataset{
		ID:          1,
		Name:        "dataset name",
		Description: "description",
	}
	dataset, _, err := client.Datasets.Update(context.Background(), newDataset)
	if err != nil {
		t.Errorf("datasets API should not return an error: %s", err)
	}

	var expectedDataset = &Dataset{}
	err = json.Unmarshal([]byte(datasetStr), expectedDataset)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDataset, dataset) {
		t.Errorf("expected dataset: %+v but got: %+v", expectedDataset, dataset)
	}
}

func TestDatasetDelete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	datasetStr := `{
		"executed": false,
		"resources_affected": {}
	  }`

	mux.HandleFunc("/datasets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, datasetStr)
	})

	client.Token.Get(context.Background())
	deleteRes, _, err := client.Datasets.Delete(context.Background(), 1, true)
	if err != nil {
		t.Errorf("datasets API should not return an error: %s", err)
	}

	var expectedDeleteResponse = &DeleteResponse{}
	err = json.Unmarshal([]byte(datasetStr), expectedDeleteResponse)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDeleteResponse, deleteRes) {
		t.Errorf("expected deleteResponse: %+v but got: %+v", expectedDeleteResponse, deleteRes)
	}
}
