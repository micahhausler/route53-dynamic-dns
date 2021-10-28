package pkg

import (
	"encoding/json"
	"io"
)

type RecordType string

const (
	RecordTypeA    = "A"
	RecordTypeAAAA = "AAAA"
)

type Config struct {
	ZoneID      string       `json:"zoneId"`
	Records     []string     `json:"records"`
	RecordTypes []RecordType `json:"recordTypes"`
	TTL         int64        `json:"ttl"`
}

func (c *Config) ARecordAllowed() bool {
	return c.recordAllowed(RecordTypeA)
}

func (c *Config) AAAARecordAllowed() bool {
	return c.recordAllowed(RecordTypeAAAA)
}

func (c *Config) recordAllowed(rtype RecordType) bool {
	resp := false
	for _, rt := range c.RecordTypes {
		if rt == rtype {
			resp = true
		}
	}
	return resp
}

func ParseConfigFile(reader io.Reader) ([]Config, error) {
	cfgs := []Config{}
	err := json.NewDecoder(reader).Decode(&cfgs)
	if err != nil {
		return nil, err
	}
	return cfgs, nil
}
