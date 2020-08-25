package tests

import (
	"LoadTest/src/util/log"
	"testing"
)

func TestLog(t *testing.T) {
	log.Info.Println("infoooo")
	log.Warning.Println("Warningggg")
	log.Error.Println("Errorrrr")
}
