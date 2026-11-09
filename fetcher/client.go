package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"groupie-tracker/data"
)

func FetchData[T any](
	client *http.Client,
	url string,
	target *T,
) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch url %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unexpected status code: %d for url: %s",
			resp.StatusCode,
			url,
		)
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	return nil
}

type GroupieData struct {
	Artists []data.Artist

	Locations struct {
		Index []data.Location `json:"index"`
	}

	Dates struct {
		Index []data.Date `json:"index"`
	}

	Relations struct {
		Index []data.Relation `json:"index"`
	}
}

func FetchAllData(baseURL string) (*GroupieData, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	dataStruct := &GroupieData{}

	var wg sync.WaitGroup
	wg.Add(4)

	errChan := make(chan error, 4)

	// Fetch Artists
	go func() {
		defer wg.Done()

		if err := FetchData(
			client,
			baseURL+"/artists",
			&dataStruct.Artists,
		); err != nil {
			errChan <- fmt.Errorf("artists error: %w", err)
		}
	}()

	// Fetch Locations
	go func() {
		defer wg.Done()

		if err := FetchData(
			client,
			baseURL+"/locations",
			&dataStruct.Locations,
		); err != nil {
			errChan <- fmt.Errorf("locations error: %w", err)
		}
	}()

	// Fetch Dates
	go func() {
		defer wg.Done()

		if err := FetchData(
			client,
			baseURL+"/dates",
			&dataStruct.Dates,
		); err != nil {
			errChan <- fmt.Errorf("dates error: %w", err)
		}
	}()

	// Fetch Relations
	go func() {
		defer wg.Done()

		if err := FetchData(
			client,
			baseURL+"/relation",
			&dataStruct.Relations,
		); err != nil {
			errChan <- fmt.Errorf("relations error: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return dataStruct, nil
}
