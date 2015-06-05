package models

import "time"

type UnixTime int64

func (t *UnixTime) Time() time.Time {
	return time.Unix(int64(*t), 0)
}
