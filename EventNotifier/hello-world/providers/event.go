package providers

import (
	"time"

	"../utils"
)

//Event type to be returned to user
type Event struct {
	Title     string
	Location  utils.Location
	Link      string
	StartDate time.Time
	EndDate   time.Time
}
