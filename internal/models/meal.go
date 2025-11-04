package models

import "time"

type Meal struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	Comment *string   `json:"comment,omitempty"`
}
