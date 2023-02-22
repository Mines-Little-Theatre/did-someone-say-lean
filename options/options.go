package options

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// optional settings configured via environment variables
var (
	FallbackReaction            = os.Getenv("LEAN_FALLBACK_REACTION")
	IgnoreUsers                 = getEnvList("LEAN_IGNORE_USERS")
	Gigglesnort                 = getJsonStringMap("LEAN_GIGGLESNORT_FILE")
	GigglesnortFallbackReaction = os.Getenv("LEAN_GIGGLESNORT_FALLBACK_REACTION")
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

// getJsonStringMap checks the environment variable with the given name
// and tries to load that JSON file as a map[string]string
func getJsonStringMap(key string) map[string]string {
	filename, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println("could not read", key, "file:", err)
		return nil
	}

	var result map[string]string
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println("could not unmarshal", key, "JSON:", err)
		return nil
	}

	return result
}
