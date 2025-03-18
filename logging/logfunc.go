package logging

import (
	"fmt"
	"log"

	"github.com/canonical/go-dqlite/v3/client"
)

func LogFunc(level client.LogLevel, format string, a ...any) {
	log.Printf(fmt.Sprintf("%s: %s\n", level.String(), format), a...)
}
