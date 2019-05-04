package mopinion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccount(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	accountStr := `{
		"_meta":{
			"code" : 200,
			"count": 1,
			"hasMore":false,
			"message": "OK",
			"next": false,
			"previous": false,
			"total":1
		},
		"name": "account name",
		"package": "package",
		"endDate": "2019-12-31 23:00:00",
		"number_users" : 205,
		"number_charts": 12,
		"number_forms": 50,
		"number_reports": 26,
		"reports": [
			{
				"id": 1,
				"name": "report",
				"description": "report description",
				"language": "en_US",
				"created": "2019-05-02",
				"datasets": [
					{
						"id": 1,
						"name": "dataset name",
						"report_id": 1,
						"description": "dataset description",
						"data_source": "form"
					}
				]
			}
		]
	}`

	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, accountStr)
	})

	client.Token.Get(context.Background())
	account, _, err := client.Account.Get(context.Background())
	if err != nil {
		t.Errorf("account API should not return an error: %s", err)
	}

	var expectedAccount = &Account{}
	err = json.Unmarshal([]byte(accountStr), expectedAccount)
	if err != nil {
		t.Errorf("unmarshaling should not return an error: %s", err)
	}

	if !reflect.DeepEqual(expectedAccount, account) {
		t.Errorf("expected account: %+v but got: %+v", expectedAccount, account)
	}
}
