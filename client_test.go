package gosail

import (
	"flag"
	"log"
	"testing"
)

var apiKey = flag.String("apikey", "", "API Key for SailThru")
var secretKey = flag.String("secretkey", "", "Secret Key for SailThru")

func checkKeys(t *testing.T) {
	if *apiKey == "" {
		t.Error(`Missing or blank required flag "-apikey={key}"`)
	}
	if *secretKey == "" {
		t.Error(`Missing or blank required flag "-secretkey={key}"`)
	}
}

func TestFlag(t *testing.T) {
	checkKeys(t)
}

func TestCreateJob(t *testing.T) {
	sc := NewSailThruClient(*apiKey, *secretKey)
	resp, err := sc.CreateJob("export_list_data", "ad_hoc_test_list_1", "json")
	if err != nil {
		t.Error(err)
	}
	if resp.JobID == "" {
		log.Println(resp)
		t.Errorf("JobID should not be empty string")
	}
}
