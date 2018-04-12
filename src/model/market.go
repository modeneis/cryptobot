package model

import (
	"time"
)

// The main structure for database
// The objective is to be able to get a standard valid for all markets, if possible.
// We will start with bitTrex and poloniex, then progressively updating it.
type Market struct {
	ID        string    `json:"_id,omitempty"         bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty"        bson:"name,omitempty"`
	Ask       float64   `json:"ask,omitempty"         bson:"ask,omitempty"`
	Bid       float64   `json:"bid,omitempty"         bson:"bid,omitempty"`
	Last      float64   `json:"last,omitempty"        bson:"last,omitempty"`
	Pair      string    `json:"pair,omitempty"        bson:"pair,omitempty"`
	UpdatedAt time.Time `json:"updatedat,omitempty"   bson:"updatedat,omitempty"`
}
