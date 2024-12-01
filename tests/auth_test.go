package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/beevik/guid"

	"github.com/KAA295/medods/api/types"
)

const waitTime = 30 * time.Second

type GenerateTokensTestCase struct {
	Request        types.GenerateTokensRequest
	ExpectedStatus int
}

func TestGenerateTokens(t *testing.T) {
	GUID := guid.NewString()
	tests := []GenerateTokensTestCase{
		{
			Request: types.GenerateTokensRequest{
				UserID: GUID,
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Request: types.GenerateTokensRequest{
				UserID: GUID,
			},
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Request: types.GenerateTokensRequest{
				UserID: "not guid",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for i, test := range tests {
		body, err := json.Marshal(test.Request)
		if err != nil {
			t.Fatalf("could not marshal to JSON: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "http://service:8000/generate_tokens", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("error while formatting request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf("response number %d was incorrect, got: %d, want: %d.", i, resp.StatusCode, test.ExpectedStatus)
		}

	}
}

type RefreshTokensTestCase struct {
	Request        types.RefreshTokensRequest
	AuthHeader     string
	ExpectedStatus int
}

func TestRefreshTokens(t *testing.T) {
	GUID := guid.NewString()
	request := types.GenerateTokensRequest{
		UserID: GUID,
	}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("could not marshal to JSON: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://service:8000/generate_tokens", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("error while formatting request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	responseStruct := new(types.TokensResponse)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading body: %v", err)
		return
	}

	err = json.Unmarshal(respBody, responseStruct)
	if err != nil {
		t.Fatalf("could not unmarshal JSON: %v", err)
	}
	if responseStruct.AccessToken == "" || responseStruct.RefreshToken == "" {
		t.Fatal("wrong response")
	}
	secondGUID := guid.NewString()

	testCases := []RefreshTokensTestCase{
		{
			Request: types.RefreshTokensRequest{ // Not expired
				UserID:       GUID,
				RefreshToken: responseStruct.RefreshToken,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID: GUID,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID:       secondGUID,
				RefreshToken: responseStruct.RefreshToken,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID:       GUID,
				RefreshToken: responseStruct.RefreshToken,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusOK,
		},
	}

	for i, test := range testCases {
		if test.ExpectedStatus == http.StatusOK {
			time.Sleep(waitTime)
		}
		body, err := json.Marshal(test.Request)
		if err != nil {
			t.Fatalf("could not marshal to JSON: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "http://service:8000/refresh_tokens", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("error while formatting request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.AuthHeader)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf("response number %d was incorrect, got: %d, want: %d.", i, resp.StatusCode, test.ExpectedStatus)
		}

	}
}
