package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const datasetFeedbackStr = `{
	"data": [
		{
			"id": 1740961,
			"created": "2019-05-02",
			"report_id": 1,
			"dataset_id": 2,
			"fields": [
				{
					"key": "123.INPUT.x3588gw1",
					"label": "label value",
					"value": "field value"
				},
				{
					"key": "124.THUMBS.prkx951l",
					"label": "label value 2",
					"value": "positive"
				}
			],
			"tags": []
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

var filterCollection = &FilterCollection{
	Filters: []Filter{
		{
			Key:      Date,
			Modifier: Gte,
			Value:    "2019-10-01",
		},
	},
}

var paginationOptions = &PaginationOptions{
	Page:  1,
	Limit: 10,
}

func TestDatasetFeedback(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/datasets/2/feedback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, datasetFeedbackStr)
	})

	client.Token.Get(context.Background())
	feedback, _, err := client.Feedback.GetByDataset(context.Background(), 2, nil, nil)
	if err != nil {
		t.Errorf("feedback API should not return any error: %s", err)
	}

	var expectedFeedback = &Feedback{}
	err = json.Unmarshal([]byte(datasetFeedbackStr), expectedFeedback)
	if err != nil {
		t.Errorf("unmarshaling should not return any error: %s", err)
	}

	if !reflect.DeepEqual(expectedFeedback, feedback) {
		t.Errorf("expected datasetFeedback: %+v but got: %+v", expectedFeedback, feedback)
	}
}

func TestDatasetFeedbackWithFilters(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/datasets/2/feedback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery == "filter[>>date]=2019-10-01" {
			fmt.Fprint(w, datasetFeedbackStr)
		}
	})

	client.Token.Get(context.Background())

	feedback, _, err := client.Feedback.GetByDataset(context.Background(), 2, nil, filterCollection)
	if err != nil {
		t.Errorf("feedback API should not return any error: %s", err)
	}

	var expectedFeedback = &Feedback{}
	err = json.Unmarshal([]byte(datasetFeedbackStr), expectedFeedback)
	if err != nil {
		t.Errorf("unmarshaling should not return any error: %s", err)
	}

	if !reflect.DeepEqual(expectedFeedback, feedback) {
		t.Errorf("expected datasetFeedback: %+v but got: %+v", expectedFeedback, feedback)
	}
}

func TestDatasetFeedbackWithFiltersAndPagination(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/datasets/2/feedback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery == "limit=10&page=1&filter[>>date]=2019-10-01" {
			fmt.Fprint(w, datasetFeedbackStr)
		}
	})

	client.Token.Get(context.Background())
	feedback, _, err := client.Feedback.GetByDataset(context.Background(), 2, paginationOptions, filterCollection)
	if err != nil {
		t.Errorf("feedback API should not return any error: %s", err)
	}

	var expectedFeedback = &Feedback{}
	err = json.Unmarshal([]byte(datasetFeedbackStr), expectedFeedback)
	if err != nil {
		t.Errorf("unmarshaling should not return any error: %s", err)
	}

	if !reflect.DeepEqual(expectedFeedback, feedback) {
		t.Errorf("expected datasetFeedback: %+v but got: %+v", expectedFeedback, feedback)
	}
}
