package lib

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

// Map2Struct output must be a pointer to struct
func Map2Struct(input map[string]interface{}, output interface{}) error {
	return mapstructure.Decode(input, output)
}

func GetNowTimeString() string {
	return time.Now().Format("15:04:05")
}

func GetNowDateTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
