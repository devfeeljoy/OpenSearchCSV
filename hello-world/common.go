package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// OpenSearch 응답 구조체 정의
type OpenSearchResponse struct {
	Aggregations struct {
		UniqueCategories struct {
			Buckets []struct {
				Key string `json:"key"`
			} `json:"buckets"`
		} `json:"unique_categories"`
	} `json:"aggregations"`
}

// OpenSearch에 쿼리를 요청하고 결과를 반환하는 함수
func queryOpenSearch(query string, url string, auth string) (*OpenSearchResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response OpenSearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// 카테고리 목록을 비교하고 누락된 항목을 찾는 함수
func findMissingCategories(influencersCategories []string, influencerCategories []string) []string {
	missingCategories := []string{}
	categorySet := make(map[string]bool)

	for _, category := range influencerCategories {
		categorySet[category] = true
	}

	for _, category := range influencersCategories {
		if _, found := categorySet[category]; !found {
			missingCategories = append(missingCategories, category)
		}
	}
	return missingCategories
}

// OpenSearch 응답에서 카테고리 키 추출하는 함수
func getCategoryKeys(response *OpenSearchResponse) []string {
	var keys []string
	for _, bucket := range response.Aggregations.UniqueCategories.Buckets {
		keys = append(keys, bucket.Key)
	}
	return keys
}
