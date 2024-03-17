package utils

import "github.com/spf13/viper"

func ReadDotEnv() error {
	file := viper.GetString("env-file")
	if file == "" {
		return nil
	}

	// Read the file and set the environment variables
	config := viper.New()
	config.SetConfigFile(file)
	err := config.ReadInConfig()
	if err != nil {
		return err
	}

	// Set the viper env to []string
	var configSlice []string
	for _, key := range config.AllKeys() {
		configSlice = append(configSlice, key+"="+config.GetString(key))
	}
	viper.Set("env", configSlice)

	return nil
}

func ConvertEnvToMap() map[string]string {
	env := viper.GetStringSlice("env")
	envMap := make(map[string]string)
	for _, e := range env {
		envMap[e] = e
	}
	return envMap
}

func ConvertEnvToMapP() map[string]*string {
	env := viper.GetStringSlice("env")
	envMap := make(map[string]*string)
	for _, e := range env {
		envMap[e] = &e
	}
	return envMap
}
