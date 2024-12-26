package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP 클라이언트 인터페이스
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	defaultClient HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// SendRequest HTTP 요청 전송 헬퍼 함수
func SendRequest(ctx context.Context, method, url string, body, response interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 성공 상태 코드가 아닌 경우
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 응답 파싱
	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// GetAuthUserInfo Auth Service에서 사용자 정보 조회
func GetAuthUserInfo(ctx context.Context, authServiceURL, userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/users/%s", authServiceURL, userID)
	var response map[string]interface{}

	if err := SendRequest(ctx, http.MethodGet, url, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get auth user info: %w", err)
	}

	return response, nil
}
