package model

import (
	"time"
)

// RSVP represents an RSVP entry in the domain model
type RSVP struct {
	ID        string
	Name      string
	Attending bool
	Message   string
	GroupName string
	CreatedAt time.Time
}
