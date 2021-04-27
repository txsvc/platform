package env

import (
	"os"
	"strconv"
)

// GetString returns the environment variable ENV['env'] or def if 'env' is not set.
// Note: def is only returned if the 'env' is not set, i.e. an EMPTY 'env' is still returned !
func GetString(env, def string) string {
	e, ok := os.LookupEnv(env)
	if ok {
		return e
	}
	return def
}

// GetInt returns the environment variable ENV['env'] or def if 'env' is not set.
func GetInt(env string, def int64) int64 {
	e, ok := os.LookupEnv(env)
	if ok {
		v, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return def
		}
		return v
	}
	return def
}

// Assert verifies that the environment variable 'env' is defined and returns a non-empty value.
func Assert(env string) bool {
	e, ok := os.LookupEnv(env)
	if !ok || e == "" {
		return false
	}
	return true
}
