package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Printf("error reading stdin %s", err)
	}
	lines := strings.Split(string(bs), "\n")

	r, lines, err := processRequest(lines)
	if err != nil {
		log.Printf("error reading stdin %s", err)
	}

	h, lines, err := processHeaders(lines)
	if err != nil {
		log.Printf("error reading stdin %s", err)
	}
	r.Header = h

	b, _, err := processBody(lines)
	r.Body = b

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Printf("error reading stdin %s", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading stdin %s", err)
	}

	fmt.Printf("%s", respBody)
}

func processRequest(lines []string) (*http.Request, []string, error) {
	if len(lines) < 3 {
		return nil, lines, fmt.Errorf("not enough lines to build request, need method and URL lines")
	}
	remainder := lines[3:]
	method := lines[0]
	rawUrl := lines[1]
	emptyLine := lines[2]

	switch method {
	case "POST", "GET":
		// This is fine, carry on
	default:
		return nil, remainder, fmt.Errorf("unrecognised method, expecting either GET or POST found %s", method)
	}

	// TODO would be good to parse the URL and return early errors if it's no good
	//myUrl := url.Parse(string(line[1]))

	r, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		return nil, remainder, err
	}

	if emptyLine != "" {
		return nil, remainder, fmt.Errorf("third line should be empty, read %q", emptyLine)
	}

	return r, remainder, nil
}

func processHeaders(lines []string) (http.Header, []string, error) {
	headers := http.Header{}
	for i, line := range lines {
		if line == "" {
			// Empty strings means end of headers
			return headers, lines[i:], nil
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("malformed header line %v", parts)
		}
		headers.Add(parts[0], parts[1])
	}
	return headers, []string{}, nil
}

func processBody(lines []string) (io.ReadCloser, []string, error) {
	buf := &bytes.Buffer{}
	for i, line := range lines {
		if _, err := buf.WriteString(line); err != nil {
			return nil, []string{}, err
		}
		if i != len(lines)-1 {
			if _, err := buf.WriteString("\n"); err != nil {
				return nil, []string{}, err
			}
		}
	}
	return io.NopCloser(buf), []string{}, nil
}
