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
		switch e := err.(type) {
		case *RequestError:
			// e is of type *RequestError
			fmt.Println("RequestError", e.Error())
		case *OtherError:
			// e is of type *OtherError
			fmt.Println("other error", e.Error())
		}

		//or using down casting
		if err != nil {
			// This is how you do a type case. the first return is of *RequestError type and the second is a boolean if it was successful.
			if _, ok := err.(*RequestError); ok {
				fmt.Println("RequestError", err)
			} else if _, ok := err.(*OtherError); ok {
				fmt.Println("other error", err)
			}
		}
	}
}
