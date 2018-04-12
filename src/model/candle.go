package model

import (
	"fmt"
	"time"
)

type Candle struct {
	ID         string     `json:"_id,omitempty"            bson:"_id,omitempty"`
	Open       float64    `json:"O,omitempty"              bson:"open,omitempty"`
	Close      float64    `json:"C,omitempty"              bson:"close,omitempty"`
	High       float64    `json:"H,omitempty"              bson:"high,omitempty"`
	Low        float64    `json:"L,omitempty"              bson:"low,omitempty"`
	Volume     float64    `json:"V,omitempty"              bson:"volume,omitempty"`
	BaseVolume float64    `json:"B,omitempty"              bson:"basevolume,omitempty"`
	TimeStamp  CandleTime `json:"T,omitempty"              bson:"timestamp,omitempty"`
}

type Candles struct {
	CandleLS []Candle `json:"ticks"`
}

type CandleTime struct {
	time.Time
}

func (t *CandleTime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("could not parse time %s", string(b))
	}

	result, err := time.Parse("2006-01-02T15:04:05", string(b[1:len(b)-1]))
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	t.Time = result
	return nil
}
