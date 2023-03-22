package machines

import (
	"os"
)

func envVarWithDefault(key string, defValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		val = defValue
	}
	return val
}
