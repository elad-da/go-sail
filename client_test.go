package gosail

import (
	"flag"
	"testing"
)

var apikey = flag.String("apikey", "", "API Key for SailThru")
var secretkey = flag.String("secretkey", "", "Secret Key for SailThru")

func checkKeys(t *testing.T) {
	if *apikey == "" {
		t.Error(`Missing or blank required flag "-apikey={key}"`)
	}
	if *secretkey == "" {
		t.Error(`Missing or blank required flag "-secretkey={key}"`)
	}
}

func TestFlag(t *testing.T) {
	checkKeys(t)

}
