package cli

import (
	"time"
)

type Info struct {
	Version string `json:"version" yaml:"version"`
	Time    string `json:"time" yaml:"time"`
}

func NewInfo() *Info {
	return &Info{
		Version: Version,
		Time:    time.Now().Format(time.UnixDate),
	}
}

func (i *Info) GetTime() (time.Time, error) {
	return time.Parse(time.UnixDate, i.Time)
}

func (i *Info) GetAge() (*time.Duration, error) {
	if t, err := i.GetTime(); err == nil {
		duration := time.Since(t)
		return &duration, nil
	} else {
		return nil, err
	}
}
