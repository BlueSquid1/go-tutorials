package main

import (
	"errors"
	"fmt"
)

type RequestError struct {
	error
}

type OtherError struct {
	error
}

func MakeRequest(url string) (string, error) {
	if url == "doesn't exist" {
		return "", &RequestError{errors.New("does not exist")}
	}
	return "", &OtherError{errors.New("other error")}
}

func main() {
	args := []string{"doesn't exist", "test"}

	for _, arg := range args {
		_, err := MakeRequest(arg)
		switch err.(type) {
		case *RequestError:
			fmt.Println("RequestError", err)
		case *OtherError:
			fmt.Println("other error", err)
		}

		//or
		if err != nil {
			if _, ok := err.(*RequestError); ok {
				fmt.Println("RequestError", err)
			} else if _, ok := err.(*OtherError); ok {
				fmt.Println("other error", err)
			}
		}
	}
}
