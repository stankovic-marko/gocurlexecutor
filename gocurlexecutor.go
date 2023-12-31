package gocurlexecutor

import (
	"errors"
	"net/http"
	"strings"

	"github.com/stankovic-marko/gocurlexecutor/util"
)

// ArgumentParser is used to define a parsing function for each supported option.
type ArgumentParser func(arguments []string, index int) (map[string]string, error)

var argumentParsers = map[string]ArgumentParser{
	"-X": parseMethodArgument,
	"-d": parseDataArgument,
	"-H": parseHeaderArgument,
	"-b": parseCookieArgument,
}

// Execute executes command and returns an HTTP response.
func Execute(command string) (http.Response, error) {

	options, err := Parse(command)
	if err != nil {
		return http.Response{}, err
	}

	response, err := sendRequest(options)
	if err != nil {
		return http.Response{}, nil
	}

	return response, nil
}

func sendRequest(options map[string]string) (http.Response, error) {

	client := &http.Client{}
	data := options["data"]
	request, err := http.NewRequest(options["method"], options["url"], strings.NewReader(data))
	if err != nil {
		return http.Response{}, errors.New("error while creating request")
	}

	headers := util.GetHeaders(options)
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	cookies := util.GetCookies(options)
	for key, value := range cookies {
		cookie := &http.Cookie{
			Name:   key,
			Value:  value,
			MaxAge: 300,
		}
		request.AddCookie(cookie)
	}

	response, err := client.Do(request)
	if err != nil {
		return http.Response{}, errors.New("error while sending request")
	}

	return *response, nil
}

// Parse parses command and returns command options.
func Parse(command string) (map[string]string, error) {

	arguments := strings.Split(command, " ")

	err := validateCommand(arguments)
	if err != nil {
		return nil, err
	}

	options, err := getOptions(arguments)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func validateCommand(arguments []string) error {
	if len(arguments) < 2 {
		return errors.New("curl command must have at least 2 arguments")
	}

	if arguments[0] != "curl" {
		return errors.New("missing 'curl' keyword")
	}
	return nil
}

func getOptions(arguments []string) (map[string]string, error) {
	options := make(map[string]string)
	for index, argument := range arguments {
		if !strings.HasPrefix(argument, "-") {
			options["url"] = argument
		}
		if parser, ok := argumentParsers[argument]; ok {
			data, err := parser(arguments, index)
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

func parseMethodArgument(arguments []string, index int) (map[string]string, error) {

	method := arguments[index+1]
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return nil, errors.New("invalid method : " + method)
	}

	options := make(map[string]string)
	options["method"] = method

	return options, nil
}

func parseDataArgument(arguments []string, index int) (map[string]string, error) {

	data := arguments[index+1]
	if !strings.HasPrefix(data, `"`) || !strings.HasSuffix(data, `"`) {
		return nil, errors.New(`data must be contained within ""`)
	}
	options := make(map[string]string)
	options["data"] = data[1 : len(data)-1]

	return options, nil
}

func parseHeaderArgument(arguments []string, index int) (map[string]string, error) {
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

func parseCookieArgument(arguments []string, index int) (map[string]string, error) {
	data := arguments[index+1]
	if !strings.HasPrefix(data, `"`) || !strings.HasSuffix(data, `"`) {
		return nil, errors.New(`header must be contained within ""`)
	}
	data = data[1 : len(data)-1]
	components := strings.Split(data, "=")
	if len(components) != 2 {
		return nil, errors.New("invalid cookie definition")
	}
	key := components[0]
	value := components[1]
	key = "C-" + key

	options := make(map[string]string)
	options[key] = value

	return options, nil
}
