package main

import (
	"log"
	"os"
	"testing"
)

func TestNotInternal(t *testing.T) {
	log.SetOutput(os.Stdout)

}
