// Package endpointtests implements users tests for the API layer.
package endpointtests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ArdanStudios/gotraining/11-http/api/app"
	"github.com/ArdanStudios/gotraining/11-http/api/models"
	"github.com/ArdanStudios/gotraining/11-http/api/routes"
	"github.com/ArdanStudios/gotraining/11-http/api/tests"
)

// Test_Users is the entry point for the users tests.
func Test_Users(t *testing.T) {
	c := &app.Context{
		Session:   app.GetSession(),
		SessionID: "TESTING",
	}
	defer c.Session.Close()

	usersCreate200(t, c)
	usersCreate409(t, c)
	usersList200(t, c)
}

// usersCreate200 validates a user can be created with the endpoint.
func usersCreate200(t *testing.T, c *app.Context) {
	u := models.User{
		UserType:  1,
		FirstName: "Bill",
		LastName:  "Kennedy",
		Email:     "bill@ardanstugios.com",
		Company:   "Ardan Labs",
		Addresses: []models.UserAddress{
			{
				Type:    1,
				LineOne: "12973 SW 112th ST",
				LineTwo: "Suite 153",
				City:    "Miami",
				State:   "FL",
				Zipcode: "33172",
				Phone:   "305-527-3353",
			},
		},
	}

	var response struct {
		ID string
	}

	body, _ := json.Marshal(&u)
	r := tests.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	routes.TM.ServeHTTP(w, r)

	t.Log("Given the need to add a new user with the users endpoint.")
	{
		if w.Code != 200 {
			t.Fatalf("\tShould received a status code of 200 for the response. Received[%d] %s", w.Code, tests.Failed)
		}
		t.Log("\tShould received a status code of 200 for the response.", tests.Succeed)

		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal("\tShould be able to unmarshal the response.", tests.Failed)
		}
		t.Log("\tShould be able to unmarshal the response.", tests.Succeed)

		if response.ID == "" {
			t.Fatal("\tShould have a user id in the response.", tests.Failed)
		}
		t.Log("\tShould have a user id in the response.", tests.Succeed)
	}
}

// usersCreate409 validates a user can't be created with the endpoint
// unless a valid user document is submitted.
func usersCreate409(t *testing.T, c *app.Context) {
	u := models.User{
		UserType: 1,
		LastName: "Kennedy",
		Email:    "bill@ardanstugios.com",
		Company:  "Ardan Labs",
	}

	var v []app.Invalid

	body, _ := json.Marshal(&u)
	r := tests.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	routes.TM.ServeHTTP(w, r)

	t.Log("Given the need to validate a new user can't be created with an invalid document.")
	{
		if w.Code != 409 {
			t.Fatalf("\tShould received a status code of 409 for the response. Received[%d] %s", w.Code, tests.Failed)
		}
		t.Log("\tShould received a status code of 409 for the response.", tests.Succeed)

		if err := json.NewDecoder(w.Body).Decode(&v); err != nil {
			t.Fatal("\tShould be able to unmarshal the response.", tests.Failed)
		}
		t.Log("\tShould be able to unmarshal the response.", tests.Succeed)

		if len(v) == 0 {
			t.Fatal("\tShould have validation errors in the response.", tests.Failed)
		}
		t.Log("\tShould have validation errors in the response.", tests.Succeed)

		if v[0].Fld != "FirstName" {
			t.Fatalf("\tShould have a FirstName validation error in the response. Received[%s] %s", v[0].Fld, tests.Failed)
		}
		t.Log("\tShould have a FirstName validation error in the response.", tests.Succeed)

		if v[1].Fld != "Addresses" {
			t.Fatalf("\tShould have an Addresses validation error in the response. Received[%s] %s", v[0].Fld, tests.Failed)
		}
		t.Log("\tShould have an Addresses validation error in the response.", tests.Succeed)
	}
}

// usersList200 validates a users list can be retrieved with the endpoint.
func usersList200(t *testing.T, c *app.Context) {
	var us []models.User

	r := tests.NewRequest("GET", "/v1/users", nil)
	w := httptest.NewRecorder()
	routes.TM.ServeHTTP(w, r)

	t.Log("Given the need to retrieve a list of users with the users endpoint.")
	{
		if w.Code != 200 {
			t.Fatalf("\tShould received a status code of 200 for the response. Received[%d] %s", w.Code, tests.Failed)
		}
		t.Log("\tShould received a status code of 200 for the response.", tests.Succeed)

		if err := json.NewDecoder(w.Body).Decode(&us); err != nil {
			t.Fatal("\tShould be able to unmarshal the response.", tests.Failed)
		}
		t.Log("\tShould be able to unmarshal the response.", tests.Succeed)

		if len(us) == 0 {
			t.Fatal("\tShould have users in the response.", tests.Failed)
		}
		t.Log("\tShould have a users in the response.", tests.Succeed)

		var failed bool
		marks := make([]string, len(us))
		for i, u := range us {
			if u.DateCreated == nil || u.DateModified == nil {
				marks[i] = tests.Failed
				failed = true
			} else {
				marks[i] = tests.Succeed
			}
		}

		if failed {
			t.Fatalf("\tShould have dates in all the user documents. %+v", marks)
		}
		t.Logf("\tShould have dates in all the user documents. %+v", marks)
	}
}
