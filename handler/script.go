package handler

import (
	"os"
)

const (
	TEST = "test"
)

func RegisterScript() {

	input := os.Args
	inputLen := len(input)
	if inputLen < 2 {
		return
	}

	module := input[1]

	switch module {

	case TEST:

	}
}
