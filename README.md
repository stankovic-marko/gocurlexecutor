# gocurlexecutor

_gocurlexecutor_ is a Go package that provides functionality to parse and execute curl commands.

# Installation

```.bash
go get github.com/stankovic-marko/gocurlexecutor
```

# Usage

```.go
package main

import (
	"github.com/stankovic-marko/gocurlexecutor"
)

func main() {
	response, err := gocurlexecutor.Execute("curl -X GET https://github.com/stankovic-marko")
}
```

# Supported curl options

| Option | Description | Support                  |
| ------ | ----------- | ------------------------ |
| -X     | Method      | :white_check_mark:       |
| -H     | Header      | :white_check_mark:       |
| -d     | Data        | :white_check_mark:       |
| -b     | Cookie data | :white_check_mark:       |
| -x     | Proxy       | :hourglass_flowing_sand: |
