package gosail

import (
	"flag"
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
	t.Skip()
	checkKeys(t)
}

func TestCreateJob(t *testing.T) {
	expectedJobID := "555a21e5a6cba8e27427eb23"
	mc := NewMockClient(NormalJob)
	sc := NewSailThruClient(&mc, "TestAPIKey", "TestSecretKey")
	resp, err := sc.CreateJob("export_list_data", "ad_hoc_test_list_1", "json")
	if err != nil {
		t.Error(err)
	}
	if resp.JobID != expectedJobID {
		t.Errorf("Expected %v, got %v\n", expectedJobID, resp.JobID)
	}
}

func TestCreateInvalidJobType(t *testing.T) {
	expectedErrorStr := "Invalid jobType: invalid_job_type"
	mc := NewMockClient(NormalJob)
	sc := NewSailThruClient(&mc, "TestAPIKey", "TestSecretKey")

	_, err := sc.CreateJob("invalid_job_type", "ad_hoc_test_list_1", "json")
	if err == nil {
		t.Errorf("Expected %v, got %v\n", expectedErrorStr, nil)
	} else {
		if err.Error() != expectedErrorStr {
			t.Errorf("Expected %v, got %v\n", expectedErrorStr, err.Error())
		}
	}
}

func TestGetJobDownloadLink(t *testing.T) {
	expectedJobID := "555a468b975910683a63b666"
	mc := NewMockClient(NormalJob)
	sc := NewSailThruClient(&mc, "TestAPIKey", "TestSecretKey")
	j, err := sc.GetJob(expectedJobID)
	if err != nil {
		t.Error(err)
	}
	if j.JobID != expectedJobID {
		t.Errorf("Expected %v, got %v\n", expectedJobID, j.JobID)
	}
}

func TestGetJobExpired(t *testing.T) {
	expectedJobID := "555a21e5a6cba8e27427eb23"
	mc := NewMockClient(ExpiredJob)
	sc := NewSailThruClient(&mc, "TestAPIKey", "TestSecretKey")
	r, err := sc.GetJob(expectedJobID)
	if err != nil {
		t.Error(err)
	}
	if r.JobID != expectedJobID {
		t.Errorf("Expected jobID %v, got %v\n", expectedJobID, r.JobID)
	}
	if r.Expired != true {
		t.Errorf("Expected status %v, got %v", false, r.Expired)
	}
}

func TestGetInvalidJobID(t *testing.T) {
	expectedJobID := "555a468b975910683a63b667"
	mc := NewMockClient(InvalidJob)
	sc := NewSailThruClient(&mc, "TestAPIKey", "TestSecretKey")
	j, err := sc.GetJob(expectedJobID)
	if err != nil {
		t.Error(err)
	}
	if j.JobID != expectedJobID {
		t.Errorf("Expected %v, got %v\n", expectedJobID, j.JobID)
	}
}
