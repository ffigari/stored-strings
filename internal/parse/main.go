package parse

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func BodyParams(r *http.Request) map[string]string {
	var parsedParams = map[string]string{}

	if r.Body == nil {
		return parsedParams
	}

	buf := new(strings.Builder)

	if _, err := io.Copy(buf, r.Body); err != nil {
		log.Println("[parse] could not copy body", err)
		return nil
	}

	for _, param := range strings.Split(buf.String(), "&") {
		input := strings.Split(param, "=")
		if len(input) == 2 {
			parsedParams[input[0]] = input[1]
		}
	}

	return parsedParams
}
