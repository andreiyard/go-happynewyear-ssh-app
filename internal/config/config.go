package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type BaseConfig struct {
	Fps                   int
	SnowflakeRate         int      // Amount of snowflakes that spawn each frame
	SnowflakeLimitPercent int      // Percentage of screen area (0 = no limit)
	SnowflakeChars        []string // Possible chars for snowflakes
}

type ConfigSSH struct {
	BaseConfig
	Host        string
	Port        string
	HostkeyPath string
}

func Load() BaseConfig {
	godotenv.Load()

	return BaseConfig{
		Fps:                   getRequiredInt("SNOWFLAKE_FPS"),
		SnowflakeRate:         getRequiredInt("SNOWFLAKE_RATE"),
		SnowflakeLimitPercent: getOptionalInt("SNOWFLAKE_LIMIT_PERCENT", 10),
		SnowflakeChars: func(n string) []string {
			val := os.Getenv(n)
			if val == "" {
				return []string{"*", "+", "."}
			}
			return strings.Split(val, "")
		}("SNOWFLAKE_CHARS"),
	}
}

func LoadWithSSH() ConfigSSH {
	return ConfigSSH{
		BaseConfig:  Load(),
		Host:        getRequired("SNOWFLAKE_SSH_HOST"),
		Port:        getRequired("SNOWFLAKE_SSH_PORT"),
		HostkeyPath: getOptionalString("SNOWFLAKE_SSH_HOSTKEY", ".ssh/id_ed25519"),
	}
}

func getRequiredInt(varName string) int {
	val := getRequired(varName)
	return mustBeInt(varName, val)
}

func getOptionalString(varName, defaultVal string) string {
	val := os.Getenv(varName)
	if val == "" {
		return defaultVal
	}
	return val
}

func getOptionalInt(varName string, defaultVal int) int {
	val := os.Getenv(varName)
	if val == "" {
		return defaultVal
	}
	return mustBeInt(varName, val)
}

func mustBeInt(varName, varValue string) int {
	res, err := strconv.Atoi(varValue)
	if err != nil {
		panic(fmt.Sprintf("ENV VAR '%s' should be of type int", varName))
	}
	return res
}

func getRequired(varName string) string {
	val := os.Getenv(varName)
	if val == "" {
		panic(fmt.Sprintf("ENV VAR '%s' is required", varName))
	}
	return val
}
