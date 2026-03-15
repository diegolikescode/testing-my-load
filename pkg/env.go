package pkg

import (
	"fmt"
	"os"
)

func GetEnvOrDieTrying(k string) string {
	val, exists := os.LookupEnv(k)
	if !exists {
		panic(fmt.Sprintf("Variable %s not set", k))
	}
	return val
}
