package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type operationType func(int, int) (int, error)

var operations = map[string]operationType{
	"+": add,
	"-": sub,
	"/": div,
}

func sub(arg1 int, arg2 int) (int, error) {
	return arg1 - arg2, nil
}

func add(arg1 int, arg2 int) (int, error) {
	return arg1 + arg2, nil
}

func div(arg1 int, arg2 int) (int, error) {
	if arg2 == 0 {
		return 0, errors.New("can't divide by zero")
	}
	return arg1 / arg2, nil
}

func errorHandling(msg error) {
	fmt.Println(msg)
	os.Exit(-1)
}

func main() {
	instructions := [][]string{}
	instructions = append(instructions, []string{"5", "+", "4"})
	instructions = append(instructions, []string{"5", "-", "4"})
	instructions = append(instructions, []string{"8", "/", "2"})

	for _, instruction := range instructions {
		if numArgs := len(instruction); numArgs != 3 {
			errorHandling(errors.New("expected operation to have 3 sections but got: " + strconv.Itoa(numArgs)))
		}

		arg1, err := strconv.Atoi(instruction[0])
		if err != nil {
			errorHandling(err)
		}
		op := instruction[1]
		opFunc, ok := operations[op]
		if !ok {
			errorHandling(errors.New("can find operation function for: " + op))
		}
		arg2, err := strconv.Atoi(instruction[2])
		if err != nil {
			errorHandling(err)
		}

		result, err := opFunc(arg1, arg2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}
