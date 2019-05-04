package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestReportsGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	reportsStr := `{
		"name": "report name",
		"description": "report description",
		"language": "en_US",
		"id": 1,
		"dataSets": [
			{
				"name": "dataset name",
				"report_id": 1,
				"description": "",
				"id": 1,
				"data_source": "form"
			},
			{
				"name": "dataset name 2",
				"report_id": 1,
				"description": "",
				"id": 2,
				"data_source": "form"
			}
		],
		"created": "2019-05-02",
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

	mux.HandleFunc("/reports/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reportsStr)
	})

	client.Token.Get(context.Background())
	report, _, err := client.Reports.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("reports API should not return an error: %s", err)
	}

	var expectedReport = &Report{}
	err = json.Unmarshal([]byte(reportsStr), expectedReport)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedReport, report) {
		t.Errorf("expected report: %+v but got: %+v", expectedReport, report)
	}
}

func TestReportsAdd(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	reportsStr := `{
		"id": 1,
		"name": "report name",
		"description": "report description",
		"language": "en_US",
		"created": "2019-05-02"
	}`

	mux.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, reportsStr)
	})

	client.Token.Get(context.Background())
	newReport := &Report{
		Name:        "report name",
		Description: "report description",
		Language:    "en_US",
	}
	report, _, err := client.Reports.Add(context.Background(), newReport)
	if err != nil {
		t.Errorf("reports API should not return an error: %s", err)
	}

	var expectedReport = &Report{}
	err = json.Unmarshal([]byte(reportsStr), expectedReport)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedReport, report) {
		t.Errorf("expected report: %+v but got: %+v", expectedReport, report)
	}
}

func TestReportsUpdate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	reportsStr := `{
		"id": 1,
		"name": "report name",
		"description": "report description",
		"language": "en_US",
		"created": "2017-05-02"
	}`

	mux.HandleFunc("/reports/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, reportsStr)
	})

	client.Token.Get(context.Background())
	newReport := &Report{
		ID:          1,
		Name:        "report name",
		Description: "report description",
		Language:    "en_US",
	}
	report, _, err := client.Reports.Update(context.Background(), newReport)
	if err != nil {
		t.Errorf("reports API should not return an error: %s", err)
	}

	var expectedReport = &Report{}
	err = json.Unmarshal([]byte(reportsStr), expectedReport)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedReport, report) {
		t.Errorf("expected report: %+v but got: %+v", expectedReport, report)
	}
}

func TestReportsDelete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	reportsStr := `{
		"executed": false,
		"resources_affected": {}
	  }`

	mux.HandleFunc("/reports/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, reportsStr)
	})

	client.Token.Get(context.Background())
	deleteRes, _, err := client.Reports.Delete(context.Background(), 1, true)
	if err != nil {
		t.Errorf("reports API should not return an error: %s", err)
	}

	var expectedDeleteResponse = &DeleteResponse{}
	err = json.Unmarshal([]byte(reportsStr), expectedDeleteResponse)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDeleteResponse, deleteRes) {
		t.Errorf("expected deleteResponse: %+v but got: %+v", expectedDeleteResponse, deleteRes)
	}
}
