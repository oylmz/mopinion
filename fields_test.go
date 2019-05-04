package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDatasetFields(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	datasetFieldsStr := `{
		"data": [
			{
				"report_id": 1,
				"dataset_id": 2,
				"label": "label",
				"short_label": "short-label",
				"key": "A-KEY-TO_SHORT-LABEL",
				"type": "thumbs"
			}
		],
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

	mux.HandleFunc("/datasets/2/fields", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, datasetFieldsStr)
	})

	client.Token.Get(context.Background())
	fields, _, err := client.Fields.GetByDataset(context.Background(), 2)
	if err != nil {
		t.Errorf("fields API should not return an error: %s", err)
	}

	var expectedFields = &Fields{}
	err = json.Unmarshal([]byte(datasetFieldsStr), expectedFields)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedFields, fields) {
		t.Errorf("expected datasetFields: %+v but got: %+v", expectedFields, fields)
	}
}

func TestReportFields(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	reportStr := `{
		"data": [
			{
				"report_id": 1,
				"dataset_id": 2,
				"label": "label",
				"short_label": "short-label",
				"key": "A-KEY-TO_SHORT-LABEL",
				"type": "thumbs"
			}
		],
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

	mux.HandleFunc("/reports/1/fields", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reportStr)
	})

	client.Token.Get(context.Background())
	fields, _, err := client.Fields.GetByReport(context.Background(), 1)
	if err != nil {
		t.Errorf("fields API should not return an error: %s", err)
	}

	var expectedFields = &Fields{}
	err = json.Unmarshal([]byte(reportStr), expectedFields)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedFields, fields) {
		t.Errorf("expected datasetFields: %+v but got: %+v", expectedFields, fields)
	}
}
