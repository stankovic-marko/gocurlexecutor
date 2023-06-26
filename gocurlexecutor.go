package gocurlexecutor

import (
	"errors"
	"fmt"
	"strings"
)

type ArgumentParser func(arguments []string) (map[string]string, error)

var argumentParsers = map[string]ArgumentParser{
	"-X": parseMethodArgument,
	"-d": parseDataArgument,
	"-H": parseHeaderArgument,
}

func Execute(command string) (string, error) {

	options, err := Parse(command)
	if err != nil {
		return "", err
	}
	for key, value := range options {
		fmt.Printf("%s:%s\n", key, value)
	}

	// DOTO: Make an actual request

	return "", nil
}

func Parse(command string) (map[string]string, error) {

	arguments := strings.Split(command, " ")

	if len(arguments) < 2 {
		return nil, errors.New("curl command must have at least 2 arguments")
	}

	if arguments[0] != "curl" {
		return nil, errors.New("missing 'curl' keyword")
	}

	options := make(map[string]string)
	for _, argument := range arguments {
		if !strings.HasPrefix(argument, "-") {
			options["url"] = argument
		}
		if parser, ok := argumentParsers[argument]; ok {
			data, err := parser(arguments)
			if err != nil {
				return nil, err
			}
			for key, value := range data {
				options[key] = value
			}
		}
	}
	if _, ok := options["method"]; !ok {
		options["method"] = "GET"
	}

	return options, nil
}

func parseMethodArgument(arguments []string) (map[string]string, error) {
	index := findArgumentIndex(arguments, "-X")

	method := arguments[index+1]
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return nil, errors.New("invalid method : " + method)
	}

	options := make(map[string]string)
	options["method"] = method

	return options, nil
}

func parseDataArgument(arguments []string) (map[string]string, error) {

	return nil, nil
}

func parseHeaderArgument(arguments []string) (map[string]string, error) {
	// DOTO: Multiple headers
	index := findArgumentIndex(arguments, "-H")
	header := arguments[index+1]
	if !strings.HasPrefix(header, `"`) || !strings.HasSuffix(header, `"`) {
		return nil, errors.New(`header must be contained within ""`)
	}
	header = header[1 : len(header)-1]
	components := strings.Split(header, ":")
	if len(components) != 2 {
		return nil, errors.New("invalid header definition")
	}
	key := components[0]
	value := components[1]
	key = "H-" + key

	options := make(map[string]string)
	options[key] = value

	return options, nil
}

func findArgumentIndex(arguments []string, argument string) int {
	for index := range arguments {
		if arguments[index] == argument {
			return index
		}
	}
	return -1
}
