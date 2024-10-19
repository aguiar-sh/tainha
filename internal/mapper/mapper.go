package mapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/aguiar-sh/tainha/internal/config"
	"github.com/aguiar-sh/tainha/internal/util"
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

	var wg sync.WaitGroup
	errChan := make(chan error, len(responseData)*len(route.Mapping))

	for i := range responseData {
		for _, mapping := range route.Mapping {
			wg.Add(1)
			go func(item map[string]interface{}, mapping config.RouteMapping) {
				defer wg.Done()
				pathParams := extractPathParams(mapping.Path)

				for _, param := range pathParams {
					value, exists := item[param]
					if !exists {
						continue
					}

					valueStr := fmt.Sprintf("%v", value)
					mapping.Service = util.PathProtocol(mapping.Service)
					mappedURL := fmt.Sprintf("%s%s%s", mapping.Service, strings.ReplaceAll(mapping.Path, "{"+param+"}", ""), valueStr)

					log.Printf("Mapping: %s\n", mappedURL)

					resp, err := http.Get(mappedURL)
					if err != nil {
						errChan <- fmt.Errorf("error making request to %s: %v", mappedURL, err)
						return
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						errChan <- fmt.Errorf("error reading response from %s: %v", mappedURL, err)
						return
					}

					var mappedData interface{}
					if err := json.Unmarshal(body, &mappedData); err != nil {
						errChan <- fmt.Errorf("error parsing JSON from %s: %v", mappedURL, err)
						return
					}

					item[mapping.Tag] = mappedData

					if mapping.RemoveKeyMapping {
						delete(item, param)
					}
				}
			}(responseData[i], mapping)
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			log.Println(err)
		}
	}

	finalResponse, err := json.Marshal(responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final response: %w", err)
	}

	return finalResponse, nil
}
