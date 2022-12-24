package models

import "time"

type Site struct {
	ID                  int
	URL                 string
	Frequency           float64
	LastExecutionDate   *time.Time `json:"last_execution_date"`
	Sucess              *bool
	ResponseTime        *float64   `json:"response_time"`
	ResponseAverageTime *float64   `json:"response_average_time"`
	CreationDate        time.Time  `json:"creation_date"`
	LastUpdatedDate     *time.Time `json:"last_updated_date"`
}
