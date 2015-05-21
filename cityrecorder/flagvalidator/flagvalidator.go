package flagvalidator

import (
	"flag"
	"log"
)

func ValidateFlags(names []string) {
	for _, name := range names {
		validateFlag(name, flag.Lookup(name).Value.String())
	}
}

func validateFlag(name string, flag string) {
	if len(flag) <= 0 {
		log.Fatal(name + " cannot be empty")
	}
}
