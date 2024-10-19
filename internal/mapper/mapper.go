package mapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aguiar-sh/tainha/internal/config"
)

func extractPathParams(path string) []string {
	re := regexp.MustCompile(`{([^}]+)}`)
	matches := re.FindAllStringSubmatch(path, -1)
	params := make([]string, len(matches))
	for i, match := range matches {
		params[i] = match[1]
	}
	return params
}

func Map(route config.Route, response []byte) ([]byte, error) {

	var responseData []map[string]interface{}
	if err := json.Unmarshal(response, &responseData); err != nil {
		log.Println("Error parsing JSON:", err)
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	for _, item := range responseData {
		for _, mapping := range route.Mapping {
			pathParams := extractPathParams(mapping.Path)

			for _, param := range pathParams {
				value, exists := item[param]

				if exists {
					// Convert the value to a string representation
					valueStr := fmt.Sprintf("%v", value)
					mappedURL := fmt.Sprintf("%s%s%s", mapping.Service, strings.ReplaceAll(mapping.Path, "{"+param+"}", ""), valueStr)

					log.Printf("Making request to mapping: %s\n", mappedURL)
					log.Printf("Parameters: %s=%v\n", param, value)

					resp, err := http.Get(mappedURL)
					if err != nil {
						log.Printf("Error making request to %s: %v\n", mappedURL, err)
						continue
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Printf("Error reading response from %s: %v\n", mappedURL, err)
						continue
					}

					// Attempt to parse the JSON response as a list of items
					var mappedData interface{}
					if err := json.Unmarshal(body, &mappedData); err != nil {
						log.Printf("Error parsing JSON as list from %s: %v\n", mappedURL, err)
						continue
					}

					log.Printf("Mapped data: %v\n", mappedData)

					// Assign the parsed JSON list to the item using the mapping tag
					item[mapping.Tag] = mappedData

					if mapping.RemoveKeyMapping {
						delete(item, param)
					}
				}
			}
		}
	}

	finalResponse, err := json.Marshal(responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final response: %w", err)
	}

	return finalResponse, nil
}
