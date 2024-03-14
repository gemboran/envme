package forms

import (
	"fmt"
	"github.com/spf13/viper"
)

func VRequiredAndSave(key string, msg string) func(value string) error {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf(msg)
		}

		viper.Set(key, value)

		return nil
	}
}

func VSave(key string) func(value string) error {
	return func(value string) error {
		viper.Set(key, value)

		return nil
	}
}
