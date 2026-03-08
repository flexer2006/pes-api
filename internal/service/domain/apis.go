package domain

// agify.io
type AgeResponse struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Count int    `json:"count"`
}

// genderize.io
type GenderResponse struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

// nationalize.io
type NationalityResponse struct {
	Countries []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
	Name string `json:"name"`
}
