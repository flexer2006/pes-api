package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	api "github.com/flexer2006/case-person-enrichment-go/internal/service/adapters/enrichment/services"
	domain "github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func createMockResponse(statusCode int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestNewGenderAPIClient(t *testing.T) {
	t.Run("with nil client", func(t *testing.T) {
		client := api.NewGenderAPIClient(nil)
		require.NotNil(t, client, "client should not be nil")
	})

	t.Run("with custom client", func(t *testing.T) {
		mockClient := &MockHTTPClient{}
		client := api.NewGenderAPIClient(mockClient)
		require.NotNil(t, client, "client should not be nil")
	})
}

func TestAPIClient_GetGenderByName(t *testing.T) {
	testCases := []struct {
		name           string
		inputName      string
		mockResponse   func() (*http.Response, error)
		expectedGender string
		expectedProb   float64
		expectedErrMsg string
	}{
		{
			name:           "empty name",
			inputName:      "",
			mockResponse:   func() (*http.Response, error) { return nil, nil },
			expectedGender: "",
			expectedProb:   0,
			expectedErrMsg: "name cannot be empty",
		},
		{
			name:      "http client error",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return nil, errors.New("network error")
			},
			expectedGender: "",
			expectedProb:   0,
			expectedErrMsg: "failed to execute request: network error",
		},
		{
			name:      "non 200 status",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return createMockResponse(404, []byte(`{"error":"not found"}`)), nil
			},
			expectedGender: "",
			expectedProb:   0,
			expectedErrMsg: "API returned non-200 status code: status 404",
		},
		{
			name:      "invalid json response",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return createMockResponse(200, []byte(`invalid json`)), nil
			},
			expectedGender: "",
			expectedProb:   0,
			expectedErrMsg: "failed to decode API response",
		},
		{
			name:      "successful response - male gender",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				resp := domain.GenderResponse{
					Name:        "John",
					Gender:      "male",
					Probability: 0.95,
					Count:       1000,
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedGender: "male",
			expectedProb:   0.95,
			expectedErrMsg: "",
		},
		{
			name:      "successful response - female gender",
			inputName: "Maria",
			mockResponse: func() (*http.Response, error) {
				resp := domain.GenderResponse{
					Name:        "Maria",
					Gender:      "female",
					Probability: 0.98,
					Count:       2000,
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedGender: "female",
			expectedProb:   0.98,
			expectedErrMsg: "",
		},
		{
			name:      "successful response - null gender",
			inputName: "Riley",
			mockResponse: func() (*http.Response, error) {
				resp := domain.GenderResponse{
					Name:        "Riley",
					Gender:      "",
					Probability: 0,
					Count:       500,
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedGender: "",
			expectedProb:   0,
			expectedErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tc.mockResponse()
				},
			}

			client := api.NewGenderAPIClient(mockClient)
			require.NotNil(t, client)

			genderResult, prob, err := client.GetGenderByName(context.Background(), tc.inputName)

			if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedGender, genderResult)
				assert.Equal(t, tc.expectedProb, prob)
			}
		})
	}
}

func TestAPIClient_GetGenderByName_RequestParameters(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodGet, req.Method)
			assert.Equal(t, "https://api.genderize.io", req.URL.Scheme+"://"+req.URL.Host)
			assert.Equal(t, "TestName", req.URL.Query().Get("name"))

			resp := domain.GenderResponse{
				Name:        "TestName",
				Gender:      "male",
				Probability: 0.85,
				Count:       100,
			}
			body, _ := json.Marshal(resp)
			return createMockResponse(200, body), nil
		},
	}

	client := api.NewGenderAPIClient(mockClient)
	_, _, err := client.GetGenderByName(context.Background(), "TestName")
	assert.NoError(t, err)
}

func TestNewAgeAPIClient(t *testing.T) {
	t.Run("with nil client", func(t *testing.T) {
		client := api.NewAgeAPIClient(nil)
		require.NotNil(t, client, "client should not be nil")
	})

	t.Run("with custom client", func(t *testing.T) {
		mockClient := &MockHTTPClient{}
		client := api.NewAgeAPIClient(mockClient)
		require.NotNil(t, client, "client should not be nil")
	})
}

func TestAPIClient_GetAgeByName(t *testing.T) {
	testCases := []struct {
		name           string
		inputName      string
		mockResponse   func() (*http.Response, error)
		expectedAge    int
		expectedProb   float64
		expectedErrMsg string
	}{
		{
			name:           "empty name",
			inputName:      "",
			mockResponse:   func() (*http.Response, error) { return nil, nil },
			expectedAge:    0,
			expectedProb:   0,
			expectedErrMsg: "name cannot be empty",
		},
		{
			name:      "http client error",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return nil, errors.New("network error")
			},
			expectedAge:    0,
			expectedProb:   0,
			expectedErrMsg: "failed to execute request: network error",
		},
		{
			name:      "non 200 status",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return createMockResponse(404, []byte(`{"error":"not found"}`)), nil
			},
			expectedAge:    0,
			expectedProb:   0,
			expectedErrMsg: "API returned non-200 status code: status 404",
		},
		{
			name:      "invalid json response",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return createMockResponse(200, []byte(`invalid json`)), nil
			},
			expectedAge:    0,
			expectedProb:   0,
			expectedErrMsg: "failed to decode API response",
		},
		{
			name:      "successful response with low count",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				resp := domain.AgeResponse{
					Name:  "John",
					Age:   35,
					Count: 500,
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedAge:    35,
			expectedProb:   0.5,
			expectedErrMsg: "",
		},
		{
			name:      "successful response with high count",
			inputName: "Maria",
			mockResponse: func() (*http.Response, error) {
				resp := domain.AgeResponse{
					Name:  "Maria",
					Age:   28,
					Count: 2000,
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedAge:    28,
			expectedProb:   1.0,
			expectedErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tc.mockResponse()
				},
			}

			client := api.NewAgeAPIClient(mockClient)
			require.NotNil(t, client)

			age, prob, err := client.GetAgeByName(context.Background(), tc.inputName)

			if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAge, age)
				assert.Equal(t, tc.expectedProb, prob)
			}
		})
	}
}

func TestAPIClient_GetAgeByName_RequestParameters(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodGet, req.Method)
			assert.Equal(t, "https://api.agify.io", req.URL.Scheme+"://"+req.URL.Host)
			assert.Equal(t, "TestName", req.URL.Query().Get("name"))

			resp := domain.AgeResponse{
				Name:  "TestName",
				Age:   30,
				Count: 100,
			}
			body, _ := json.Marshal(resp)
			return createMockResponse(200, body), nil
		},
	}

	client := api.NewAgeAPIClient(mockClient)
	_, _, err := client.GetAgeByName(context.Background(), "TestName")
	assert.NoError(t, err)
}

func TestNewNationalityAPIClient(t *testing.T) {
	t.Run("with nil client", func(t *testing.T) {
		client := api.NewNationalityAPIClient(nil)
		require.NotNil(t, client, "client should not be nil")
	})

	t.Run("with custom client", func(t *testing.T) {
		mockClient := &MockHTTPClient{}
		client := api.NewNationalityAPIClient(mockClient)
		require.NotNil(t, client, "client should not be nil")
	})
}

func TestAPIClient_GetNationalityByName(t *testing.T) {
	testCases := []struct {
		name                string
		inputName           string
		mockResponse        func() (*http.Response, error)
		expectedNationality string
		expectedProb        float64
		expectedErrMsg      string
	}{
		{
			name:                "empty name",
			inputName:           "",
			mockResponse:        func() (*http.Response, error) { return nil, nil },
			expectedNationality: "",
			expectedProb:        0,
			expectedErrMsg:      "name cannot be empty",
		},
		{
			name:      "http client error",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return nil, errors.New("network error")
			},
			expectedNationality: "",
			expectedProb:        0,
			expectedErrMsg:      "failed to execute request: network error",
		},
		{
			name:      "non 200 status",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return createMockResponse(404, []byte(`{"error":"not found"}`)), nil
			},
			expectedNationality: "",
			expectedProb:        0,
			expectedErrMsg:      "API returned non-200 status code: status 404",
		},
		{
			name:      "invalid json response",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				return createMockResponse(200, []byte(`invalid json`)), nil
			},
			expectedNationality: "",
			expectedProb:        0,
			expectedErrMsg:      "failed to decode API response",
		},
		{
			name:      "empty countries list",
			inputName: "UnknownName",
			mockResponse: func() (*http.Response, error) {
				resp := domain.NationalityResponse{
					Name: "UnknownName",
					// Инициализируем пустой слайс анонимной структуры
					Countries: []struct {
						CountryID   string  `json:"country_id"`
						Probability float64 `json:"probability"`
					}{},
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedNationality: "",
			expectedProb:        0,
			expectedErrMsg:      "",
		},
		{
			name:      "successful response with single country",
			inputName: "John",
			mockResponse: func() (*http.Response, error) {
				resp := domain.NationalityResponse{
					Name: "John",
					Countries: []struct {
						CountryID   string  `json:"country_id"`
						Probability float64 `json:"probability"`
					}{
						{CountryID: "US", Probability: 0.8},
					},
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedNationality: "US",
			expectedProb:        0.8,
			expectedErrMsg:      "",
		},
		{
			name:      "successful response with multiple countries",
			inputName: "Maria",
			mockResponse: func() (*http.Response, error) {
				resp := domain.NationalityResponse{
					Name: "Maria",
					Countries: []struct {
						CountryID   string  `json:"country_id"`
						Probability float64 `json:"probability"`
					}{
						{CountryID: "ES", Probability: 0.4},
						{CountryID: "PT", Probability: 0.2},
						{CountryID: "IT", Probability: 0.7},
						{CountryID: "MX", Probability: 0.3},
					},
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedNationality: "IT",
			expectedProb:        0.7,
			expectedErrMsg:      "",
		},
		{
			name:      "successful response with first country being most probable",
			inputName: "Akira",
			mockResponse: func() (*http.Response, error) {
				resp := domain.NationalityResponse{
					Name: "Akira",
					Countries: []struct {
						CountryID   string  `json:"country_id"`
						Probability float64 `json:"probability"`
					}{
						{CountryID: "JP", Probability: 0.9},
						{CountryID: "KR", Probability: 0.1},
						{CountryID: "CN", Probability: 0.05},
					},
				}
				body, _ := json.Marshal(resp)
				return createMockResponse(200, body), nil
			},
			expectedNationality: "JP",
			expectedProb:        0.9,
			expectedErrMsg:      "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return tc.mockResponse()
				},
			}

			client := api.NewNationalityAPIClient(mockClient)
			require.NotNil(t, client)

			nationalityResult, prob, err := client.GetNationalityByName(context.Background(), tc.inputName)

			if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNationality, nationalityResult)
				assert.Equal(t, tc.expectedProb, prob)
			}
		})
	}
}

func TestAPIClient_GetNationalityByName_RequestParameters(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodGet, req.Method)
			assert.Equal(t, "https://api.nationalize.io", req.URL.Scheme+"://"+req.URL.Host)
			assert.Equal(t, "TestName", req.URL.Query().Get("name"))

			resp := domain.NationalityResponse{
				Name: "TestName",
				Countries: []struct {
					CountryID   string  `json:"country_id"`
					Probability float64 `json:"probability"`
				}{
					{CountryID: "US", Probability: 0.8},
				},
			}
			body, _ := json.Marshal(resp)
			return createMockResponse(200, body), nil
		},
	}

	client := api.NewNationalityAPIClient(mockClient)
	_, _, err := client.GetNationalityByName(context.Background(), "TestName")
	assert.NoError(t, err)
}
