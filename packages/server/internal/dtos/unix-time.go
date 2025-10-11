package dtos

import (
	"encoding/json"
	"time"
)

type UnixTime time.Time

func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var ms int64
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}

	t := NewUnixTime(ms)
	ut = &t

	return nil
}

func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ut).Unix() * 1000)
}

func NewUnixTime(ms int64) UnixTime {
	return UnixTime(time.Unix(ms/1000, 0))
}
