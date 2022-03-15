package teampasswordmanager

import "fmt"

type ClientConfig2 struct {
	BaseURL   string
	AuthToken string
}

func doStuff() (bool, error) {
	output := 2+2
	other_stuff := output + 2
	fmt.Println(other_stuff)
	return false, nil
}