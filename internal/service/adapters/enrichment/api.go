package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	domain "github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/ports"
	logger "github.com/flexer2006/case-person-enrichment-go/internal/utilities"

	"go.uber.org/zap"
)

const (
	ageBaseURL         = "https://api.agify.io"
	genderBaseURL      = "https://api.genderize.io"
	nationalityBaseURL = "https://api.nationalize.io"
)

type api struct {
	httpCli                           *http.Client
	ageURL, genderURL, nationalityURL string
}

func NewAPI() ports.API {
	return new(api{
		httpCli:        http.DefaultClient,
		ageURL:         ageBaseURL,
		genderURL:      genderBaseURL,
		nationalityURL: nationalityBaseURL,
	})
}

func (a *api) GetAgeByName(ctx context.Context, name string) (int, float64, error) {
	return predict(a, ctx, name, "age", a.ageURL, func(resp *domain.AgeResponse) (int, float64, error) {
		return resp.Age, min(float64(resp.Count)/1000.0, 1.0), nil
	})
}

func (a *api) GetGenderByName(ctx context.Context, name string) (string, float64, error) {
	return predict(a, ctx, name, "gender", a.genderURL, func(resp *domain.GenderResponse) (string, float64, error) {
		return resp.Gender, resp.Probability, nil
	})
}

func (a *api) GetNationalityByName(ctx context.Context, name string) (string, float64, error) {
	return predict(a, ctx, name, "nationality", a.nationalityURL, func(resp *domain.NationalityResponse) (string, float64, error) {
		if len(resp.Countries) == 0 {
			logger.Debug(ctx, "no nationality data found for name", zap.String("name", name))
			return "", 0, nil
		}
		mostProbable := resp.Countries[0]
		for _, country := range resp.Countries {
			if country.Probability > mostProbable.Probability {
				mostProbable = country
			}
		}
		return mostProbable.CountryID, mostProbable.Probability, nil
	})
}

func getJSON[T any](a *api, ctx context.Context, baseURL string, params map[string]string) (*T, error) {
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		logger.Error(ctx, "failed to parse base URL", zap.Error(err))
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	q := reqURL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		logger.Error(ctx, "failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	//nolint:gosec
	resp, err := a.httpCli.Do(req)
	if err != nil {
		logger.Error(ctx, "request execution failed", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn(ctx, "failed to close response body", zap.Error(cerr))
		}
	}()
	if resp.StatusCode != http.StatusOK {
		logger.Error(ctx, "API returned non-200 status code", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("%w: status %d", domain.ErrNon200Response, resp.StatusCode)
	}
	out := new(T)
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		logger.Error(ctx, "failed to decode response body", zap.Error(err))
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}
	return out, nil
}

func predict[Resp any, Res any](
	a *api,
	ctx context.Context,
	name, kind, baseURL string,
	mapper func(*Resp) (Res, float64, error),
) (res Res, prob float64, err error) {
	logger.Debug(ctx, "getting "+kind+" for name", zap.String("name", name))
	if err = checkName(ctx, name); err != nil {
		return
	}
	var resp *Resp
	resp, err = getJSON[Resp](a, ctx, baseURL, map[string]string{"name": name})
	if err != nil {
		return
	}
	res, prob, err = mapper(resp)
	return
}

func checkName(ctx context.Context, name string) error {
	if name == "" {
		logger.Error(ctx, "empty name provided for prediction")
		return domain.ErrEmptyName
	}
	return nil
}
