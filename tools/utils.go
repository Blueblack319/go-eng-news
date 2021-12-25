package tools

import (
	"log"
	"strings"
)

func CleanString(toClean string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(toClean)), " ")
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
