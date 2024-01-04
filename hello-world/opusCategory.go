package main

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
)

// main 함수
func main() {
	openSearchURL := "https://search-buridge-lmwxkhotosmwhumwcno32druge.ap-northeast-2.es.amazonaws.com"
	username := "buridge"
	password := "iFuYdanRBc8oPb*.J!i*PPEsK4@xVX"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	urlOpus := openSearchURL + "/opus/_search"
	urlCategories := openSearchURL + "/opus_categories/_search"
	// opus 인덱스의 firstCategory 카테고리 가져오기
	queryFirstCategoryInfluencers := `{"size":0,"aggs":{"unique_categories":{"terms":{"field":"firstCategory.keyword","size":10000}}}}`
	firstCategoryResultInfluencers, err := queryOpenSearch(queryFirstCategoryInfluencers, urlOpus, auth)
	if err != nil {
		fmt.Println("Error querying OpenSearch for firstCategory in influencers:", err)
		return
	}

	// opus_categories 인덱스의 firstCategory 카테고리 가져오기
	queryFirstCategoryCategories := `{"size":0,"aggs":{"unique_categories":{"terms":{"field":"firstCategory.keyword","size":10000}}}}`
	firstCategoryResultCategories, err := queryOpenSearch(queryFirstCategoryCategories, urlCategories, auth)
	if err != nil {
		fmt.Println("Error querying OpenSearch for firstCategory in influencer_categories:", err)
		return
	}

	// opus 인덱스의 secondCategory 카테고리 가져오기
	querySecondCategoryInfluencers := `{"size":0,"aggs":{"unique_categories":{"terms":{"field":"secondCategory.keyword","size":10000}}}}`
	secondCategoryResultInfluencers, err := queryOpenSearch(querySecondCategoryInfluencers, urlOpus, auth)
	if err != nil {
		fmt.Println("Error querying OpenSearch for secondCategory in influencers:", err)
		return
	}

	// opus_categories 인덱스의 secondCategory 카테고리 가져오기
	querySecondCategoryCategories := `{"size":0,"aggs":{"unique_categories":{"terms":{"field":"secondCategory.keyword","size":10000}}}}`
	secondCategoryResultCategories, err := queryOpenSearch(querySecondCategoryCategories, urlCategories, auth)
	if err != nil {
		fmt.Println("Error querying OpenSearch for secondCategory in influencer_categories:", err)
		return
	}
	// opus 인덱스의 third 카테고리 가져오기
	queryThirdCategoryInfluencers := `{"size":0,"aggs":{"unique_categories":{"terms":{"field":"thirdCategory.keyword","size":10000}}}}`
	thirdCategoryResultInfluencers, err := queryOpenSearch(queryThirdCategoryInfluencers, urlOpus, auth)
	if err != nil {
		fmt.Println("Error querying OpenSearch for secondCategory in influencers:", err)
		return
	}

	// opus_categories 인덱스의 thirdCategory 카테고리 가져오기
	queryThirdCategoryCategories := `{"size":0,"aggs":{"unique_categories":{"terms":{"field":"thirdCategory.keyword","size":10000}}}}`
	thirdCategoryResultCategories, err := queryOpenSearch(queryThirdCategoryCategories, urlCategories, auth)
	if err != nil {
		fmt.Println("Error querying OpenSearch for secondCategory in influencer_categories:", err)
		return
	}

	// 누락된 firstCategory 찾기
	firstCategoryMissing := findMissingCategories(
		getCategoryKeys(firstCategoryResultInfluencers),
		getCategoryKeys(firstCategoryResultCategories),
	)

	// 누락된 secondCategory 찾기
	secondCategoryMissing := findMissingCategories(
		getCategoryKeys(secondCategoryResultInfluencers),
		getCategoryKeys(secondCategoryResultCategories),
	)

	// 누락된 thirdCategory 찾기
	thirdCategoryMissing := findMissingCategories(
		getCategoryKeys(thirdCategoryResultInfluencers),
		getCategoryKeys(thirdCategoryResultCategories),
	)

	// CSV 파일 생성 및 저장
	file, err := os.Create("/Users/joupil/Desktop/highdev/뷰릿지/OpenSearchCSV/csv/누락된_콘텐츠_카테고리.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 헤더 작성
	headers := []string{"Missing First Category", "Missing Second Category", "Missing Third Category"}
	writer.Write(headers)

	// 누락된 firstCategory 작성
	for _, category := range firstCategoryMissing {
		writer.Write([]string{category, ""})
	}

	// 누락된 secondCategory 작성
	for _, category := range secondCategoryMissing {
		writer.Write([]string{"", category})
	}

	// 누락된 thirdCategory 작성
	for _, category := range thirdCategoryMissing {
		writer.Write([]string{"", category})
	}
}
