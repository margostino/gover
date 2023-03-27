package common

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetValueFor(key string) string {
	var value string
	value = viper.GetString(key)
	if &value == nil {
		Fatal(fmt.Sprintf("Cannot get value for key %s", key))
	}
	return value
}
