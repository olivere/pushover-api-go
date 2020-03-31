package pushover

import "os"

// envString finds the first non-empty value in the given environment
// strings. If all values are empty, defaultValue is returned.
func envString(defaultValue string, envVars ...string) string {
	for _, envVar := range envVars {
		if v, ok := os.LookupEnv(envVar); ok && v != "" {
			return v
		}
	}
	return defaultValue
}
