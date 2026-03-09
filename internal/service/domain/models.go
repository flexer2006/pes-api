package domain

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
	ID                     uuid.UUID `db:"id" json:"id"`
	Name                   string    `db:"name" json:"name"`
	Surname                string    `db:"surname" json:"surname"`
	Patronymic             *string   `db:"patronymic" json:"patronymic,omitempty"`
	Age                    *int      `db:"age" json:"age,omitempty"`
	Gender                 *string   `db:"gender" json:"gender,omitempty"`
	GenderProbability      *float64  `db:"gender_probability" json:"gender_probability,omitempty"`
	Nationality            *string   `db:"nationality" json:"nationality,omitempty"`
	NationalityProbability *float64  `db:"nationality_probability" json:"nationality_probability,omitempty"`
}

type AgeResponse struct { // agify.io
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Count int    `json:"count"`
}

type GenderResponse struct { // genderize.io
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type NationalityResponse struct { // nationalize.io
	Countries []struct {
		Probability float64 `json:"probability"`
		CountryID   string  `json:"country_id"`
	} `json:"country"`
	Name string `json:"name"`
}
