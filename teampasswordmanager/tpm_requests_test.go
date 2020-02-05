// +build unit

package teampasswordmanager

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomFields(t *testing.T) {
	cf1 := CustomField{
		Label: "one",
		Data:  "1",
	}

	cf2 := CustomField{
		Label: "two",
		Data:  "2",
	}
	// Create a Password struct
	password := Password{
		CustomField1: cf1,
		CustomField2: cf2,
	}

	// check the custom field array is returned
	arrayCustomFields := password.CustomFields()

	if arrayCustomFields[0] != cf1 ||
		arrayCustomFields[1] != cf2 {
		t.Errorf(
			"Array of customfields not initialised correctly, expected (%s) and (%s) but array returned was (%s)",
			cf1,
			cf2,
			arrayCustomFields,
		)
	}
}

func TestGetPasswordByID(t *testing.T) {
	password := Password{
		ID:   1,
		Name: "postgres",
		Project: Project{
			ID:   2,
			Name: "stage.devops__foo--bar",
		},
		Tags:     "tag1, tag2",
		Username: "username",
		Password: "secretpassword",
	}
	passwordAsJSON, err := json.Marshal(password)
	assertNoError(t, err)

	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.String() == "/passwords/1.json" {
				rw.Write([]byte(passwordAsJSON))
			}
		}))
	// Close the server when test finishes
	defer server.Close()

	c := Client{
		apiURL:     server.URL + "/",
		authToken:  "1234",
		httpClient: server.Client(),
	}

	output, err := c.GetPassword(1)
	assertNoError(t, err)

	if output != password {
		t.Errorf("Expected (%s) but got (%s)", password.String(), output.String())
	}

}

func TestGetPasswordByName(t *testing.T) {
	// Create the test data and the server
	pass1 := Password{
		ID:   1,
		Name: "postgres",
		Project: Project{
			ID:   2,
			Name: "stage.devops__foo--bar",
		},
		Tags:     "tag1, tag2",
		Username: "username",
		Password: "secretpassword",
	}

	pass2 := Password{
		ID:   2,
		Name: "foo",
		Project: Project{
			ID:   3,
			Name: "bar",
		},
		Tags:     "tag1, tag2",
		Username: "username",
		Password: "secretpassword",
	}

	passwordList := []Password{}
	passwordList = append(passwordList, pass1)
	passwordList = append(passwordList, pass2)

	passwordListAsJSON, err := json.Marshal(passwordList)
	assertNoError(t, err)
	passwordTwoAsJSON, err := json.Marshal(pass2)
	assertNoError(t, err)

	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			if req.URL.String() == "/passwords.json" {
				rw.Write([]byte(passwordListAsJSON))
			} else if req.URL.String() == "/passwords/2.json" {
				rw.Write([]byte(passwordTwoAsJSON))
			} else {
				panic("Unexpected url requested: " + req.URL.String())
			}
		}))
	// Close the server when test finishes
	defer server.Close()

	c := Client{
		apiURL:     server.URL + "/",
		authToken:  "1234",
		httpClient: server.Client(),
	}

	// should return password 2
	output, err := c.GetPasswordByName("foo", "bar")
	assertNoError(t, err)

	if output != pass2 {
		t.Errorf("Expected (%s) but got (%s)", pass2.String(), output.String())
	}

}

func TestGetCustomFieldByName(t *testing.T) {
	pass1 := Password{
		ID:   1,
		Name: "postgres",
		Project: Project{
			ID:   2,
			Name: "stage.devops__foo--bar",
		},
		Tags:     "tag1, tag2",
		Username: "username",
		Password: "secretpassword",
		CustomField1: CustomField{
			Label: "bar",
			Data:  "baz",
		},
	}

	output, err := pass1.CustomField("foo")
	assert(t, output == "", "Output was not equal to nothing")
	assert(t, err != nil, "Error was not returned")

	output, err = pass1.CustomField("bar")
	assert(t, output == "baz", "Output was not equal to the value of the data for this label")
	assert(t, err == nil, "Error should be nil")
}
