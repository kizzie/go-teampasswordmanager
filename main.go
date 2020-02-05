package main

import (
	"fmt"

	"github.com/kizzie/go-teampasswordmanager/teampasswordmanager"
)

func main() {
	config := teampasswordmanager.ClientConfig{
		BaseURL:   "http://localhost/teampasswordmanager",
		AuthToken: "a2F0OnBhc3N3b3Jk",
	}

	client, _ := teampasswordmanager.NewClient(&config)

	fmt.Println(client.GetPasswordList())

	password, _ := client.GetPassword(1)
	fmt.Println(password)
	fmt.Println(password.CustomFields())
	// fmt.Println(client.GetPassword(2))
	// fmt.Println(client.GetPasswordByName("foo", "bar"))
	fmt.Println(password.CustomField("service_username"))
	fmt.Println(password.CustomField("service_password"))
}
