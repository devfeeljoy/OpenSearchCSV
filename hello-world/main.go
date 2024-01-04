package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func processBatch(records [][]string, headers []string, client *http.Client, req *http.Request) error {
	var bulkRequestBody bytes.Buffer

	for _, record := range records {
		metaData := map[string]map[string]interface{}{
			"index": {"_index": "product_categories"},
		}

		metaJSON, _ := json.Marshal(metaData)
		bulkRequestBody.Write(metaJSON)
		bulkRequestBody.WriteByte('\n')

		data := make(map[string]interface{})
		for i, header := range headers {
			data[header] = record[i]
		}

		dataJSON, _ := json.Marshal(data)
		bulkRequestBody.Write(dataJSON)
		bulkRequestBody.WriteByte('\n')
	}

	req.Body = ioutil.NopCloser(&bulkRequestBody)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
	return nil
}

func HandleRequest() {
	openSearchURL := "https://search-buridge-lmwxkhotosmwhumwcno32druge.ap-northeast-2.es.amazonaws.com"
	username := "buridge"
	password := "iFuYdanRBc8oPb*.J!i*PPEsK4@xVX"
	// Lambda 실행 환경에서의 파일 경로
	filePath := "/Users/joupil/Desktop/highdev/뷰릿지/OpenSearchCSV/hello-world/상품카테고리.csv"
	file, err := os.Open(filePath)

	defer file.Close()

	// CSV Reader 생성
	reader := csv.NewReader(file)
	reader.Comma = ',' // 필요에 따라 구분자 변경
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// OpenSearch 색인화를 위한 Bulk API 요청 생성
	var bulkRequestBody bytes.Buffer
	headers := records[0] // 첫 번째 행은 헤더
	client := &http.Client{}

	// OpenSearch에 데이터 색인화
	req, err := http.NewRequest("POST", openSearchURL+"/_bulk", &bulkRequestBody)

	// ID와 패스워드를 결합하고 Base64로 인코딩합니다.
	auth := username + ":" + password
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))

	// Authorization 헤더를 설정합니다.
	req.Header.Set("Authorization", "Basic "+authEncoded)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// 배치 크기 정의
	batchSize := 50 // 예시 크기
	for i := 0; i < len(records[1:]); i += batchSize {
		end := i + batchSize
		if end > len(records[1:]) {
			end = len(records[1:])
		}

		batch := records[1:][i:end]
		err := processBatch(batch, headers, client, req)
		if err != nil {
			panic(err)
		}

		//// 선택적으로 요청 간 지연 시간 추가
		//time.Sleep(1 * time.Second)
	}

}

func main() {
	HandleRequest()
}
