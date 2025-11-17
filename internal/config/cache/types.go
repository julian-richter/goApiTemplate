package cache

import (
	"fmt"
	"strconv"
)

type Port int

const (
	PortMin Port = 1
	PortMax Port = 65535
)

func (p Port) Valid() bool {
	return p >= PortMin && p <= PortMax
}

func ParsePort(raw string) (Port, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	p := Port(value)
	if !p.Valid() {
		return 0, fmt.Errorf("port must be between %d and %d, got %d", PortMin, PortMax, p)
	}
	return p, nil
}

type Config struct {
	User     string
	Host     string
	Port     Port
	Password string
	DB       int
}
