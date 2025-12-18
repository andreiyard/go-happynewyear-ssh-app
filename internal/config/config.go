package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Fps             int
	SnowflakeRate   int      // Amount of snowflakes that spawn each frame
	SnowflakesLimit int      // Max number of snowflakes on the screen
	SnowflakeChars  []string // Possible chars for snowflakes
}

func Load() Config {
	godotenv.Load()

	return Config{
		Fps:             getRequiredInt("SNOWFLAKE_FPS"),
		SnowflakeRate:   getRequiredInt("SNOWFLAKE_RATE"),
		SnowflakesLimit: getOptionalInt("SNOWFLAKE_LIMIT", 800),
		SnowflakeChars: func(n string) []string {
			val := os.Getenv(n)
			if val == "" {
				return []string{"*", "+", "."}
			}
			return strings.Split(val, "")
		}("SNOWFLAKE_CHARS"),
	}
}

func getRequiredInt(varName string) int {
	val := getRequired(varName)
	return mustBeInt(varName, val)
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
