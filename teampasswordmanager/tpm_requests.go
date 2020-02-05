package teampasswordmanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Project stores the ID and name of the tpm project
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CustomField stores the label and data for each of the custom fields
type CustomField struct {
	Label string `json:"label"`
	Data  string `json:"data"`
}

// Password uses the same structure as the json object returned by TPM
type Password struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Project         Project `json:"project"`
	NotesSnippet    string  `json:"notes_snippet"`
	Tags            string  `json:"tags"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	ExpiryDate      string  `json:"expiry_date"`
	ExpiryStatus    int     `json:"expiry_status"`
	Archived        bool    `json:"archived"`
	Favourite       bool    `json:"favourite§"`
	NumFiles        int     `json:"num_files"`
	Locked          bool    `json:"locked"`
	ExternalSharing bool    `json:"external_sharing"`
	UpdatedOn       string  `json:"updated_on"`

	CustomField1  CustomField `json:"custom_field1"`
	CustomField2  CustomField `json:"custom_field2"`
	CustomField3  CustomField `json:"custom_field3"`
	CustomField4  CustomField `json:"custom_field4"`
	CustomField5  CustomField `json:"custom_field5"`
	CustomField6  CustomField `json:"custom_field6"`
	CustomField7  CustomField `json:"custom_field7"`
	CustomField8  CustomField `json:"custom_field8"`
	CustomField9  CustomField `json:"custom_field9"`
	CustomField10 CustomField `json:"custom_field10"`
}

// PasswordList is a list of Password objects used by GetPasswordList
type PasswordList []Password

func (password *Password) String() string {
	return strconv.Itoa(password.ID) + ": " + password.Name
}

// CustomFields add a function which returns the custom field individual values as an array
func (password *Password) CustomFields() []CustomField {
	var customFields []CustomField
	customFields = append(customFields, password.CustomField1)
	customFields = append(customFields, password.CustomField2)
	customFields = append(customFields, password.CustomField3)
	customFields = append(customFields, password.CustomField4)
	customFields = append(customFields, password.CustomField5)
	customFields = append(customFields, password.CustomField6)
	customFields = append(customFields, password.CustomField7)
	customFields = append(customFields, password.CustomField8)
	customFields = append(customFields, password.CustomField9)
	customFields = append(customFields, password.CustomField10)

	return customFields
}

// CustomField params (name) - Pass in the name of a custom field and this will return the value stored
func (password *Password) CustomField(name string) (string, error) {
	all := password.CustomFields()

	for i := 0; i < 10; i++ {
		if all[i].Label == name {
			return all[i].Data, nil
		}
	}

	return "", errors.New("Could not find the Custom Field label asked for: " + name)
}

// GetPasswordList returns a list of all the Password objects returned by tpm for the key
func (c *Client) GetPasswordList() (PasswordList, error) {
	var output PasswordList
	body, tpmError := c.tpmRequest("passwords.json")

	if tpmError != nil {
		return nil, tpmError
	}

	err := json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetPassword returns a Password struct for a given ID
func (c *Client) GetPassword(id int) (Password, error) {
	var output Password
	body, err := c.tpmRequest("passwords/" + strconv.Itoa(id) + ".json")
	// log.Print(body, err)
	if body == nil || err != nil {
		log.Print("Requested password not found " + strconv.Itoa(id))
		return Password{}, err
	}

	marshallErr := json.Unmarshal(body, &output)
	if marshallErr != nil {
		return Password{}, marshallErr
	}

	return output, nil
}

// GetPasswordByName returns as Password struct for a given name and project
func (c *Client) GetPasswordByName(name string, project string) (Password, error) {
	// get the full list of passwords
	allPasswords, err := c.GetPasswordList()

	if err != nil {
		return Password{}, err
	}

	id := 0
	// find something with the correct name and project
	for _, entry := range allPasswords {
		if entry.Name == name && entry.Project.Name == project {
			id = entry.ID
		}
	}

	if id == 0 {
		return Password{}, errors.New("Requested password not found " + strconv.Itoa(id))
	}

	// return the password by the ID
	return c.GetPassword(id)
}

// return the id of the saved password
// func SavePassword(password Password) int {
//   return 1
// }

func (c *Client) tpmRequest(endpoint string) ([]byte, error) {

	urlString := c.apiURL + endpoint
	newURL, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    newURL,
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Basic "+c.authToken)
	resp, err := c.httpClient.Do(req)

	if err != nil || resp.StatusCode != 200 {
		log.Print("Failed to get " + urlString)
		log.Print(err)
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
