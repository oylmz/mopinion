package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDeploymentsGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	deploymentsStr := `{
		"0": {
			"key": "ab25of859d3",
			"name": "deployment 1"
		},
		"1": {
			"key": "dpg93g038fm",
			"name": "deployment 2"
		},
		"_meta": {
			"code": 200,
			"message": "OK",
			"has_more": false,
			"previous": false,
			"next": false,
			"count": 2,
			"total": 2
		}
	}`

	mux.HandleFunc("/deployments", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, deploymentsStr)
	})

	client.Token.Get(context.Background())
	deployments, _, err := client.Deployments.Get(context.Background())
	if err != nil {
		t.Errorf("deployments API should not return an error: %s", err)
	}

	var expectedDeployments = &Deployments{}
	err = json.Unmarshal([]byte(deploymentsStr), expectedDeployments)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDeployments, deployments) {
		t.Errorf("expected deployments: %+v but got: %+v", expectedDeployments, deployments)
	}
}

func TestDeploymentsAdd(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	deploymentsStr := `{
			"_meta": {
				"code": 200,
				"message": "OK",
				"has_more": false,
				"previous": false,
				"next": false,
				"count": 1,
				"total": 1
			},
			"0": {
				"key": "76pg3seur7occo1hogv88eltdtmxoxxl81vj",
				"name": "Default implementation"
			}
		}`

	mux.HandleFunc("/deployments", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, deploymentsStr)
	})

	client.Token.Get(context.Background())
	newDeployment := &Deployment{
		Key:  "76pg3seur7occo1hogv88eltdtmxoxxl81vj",
		Name: "Default implementation",
	}
	deployments, _, err := client.Deployments.Add(context.Background(), newDeployment)
	if err != nil {
		t.Errorf("deployments API should not return an error: %s", err)
	}

	var expectedDeployments = &Deployments{}
	err = json.Unmarshal([]byte(deploymentsStr), expectedDeployments)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDeployments, deployments) {
		t.Errorf("expected deployments: %+v but got: %+v", expectedDeployments, deployments)
	}
}

func TestDeploymentsDelete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	deploymentsStr := `{
		"executed": false,
		"resources_affected": {}
	  }`

	mux.HandleFunc("/deployments/76pg3seur7occo1hogv88eltdtmxoxxl81vj", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, deploymentsStr)
	})

	client.Token.Get(context.Background())
	deleteRes, _, err := client.Deployments.Delete(context.Background(), "76pg3seur7occo1hogv88eltdtmxoxxl81vj", true)
	if err != nil {
		t.Errorf("deployments API should not return an error: %s", err)
	}

	var expectedDeleteResponse = &DeleteResponse{}
	err = json.Unmarshal([]byte(deploymentsStr), expectedDeleteResponse)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedDeleteResponse, deleteRes) {
		t.Errorf("expected deleteResponse: %+v but got: %+v", expectedDeleteResponse, deleteRes)
	}
}
