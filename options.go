package main

import (
	"os"
	"strings"
)

// optional settings configured via environment variables
var (
	fallbackReaction = os.Getenv("LEAN_FALLBACK_REACTION")
	ignoreUsers      = getEnvList("LEAN_IGNORE_USERS")
)

// getEnvList splits the value of the environment variable at commas,
// returning an empty slice if the variable is empty or unset
func getEnvList(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		return []string{}
	} else {
		return strings.Split(value, ",")
	}
}
