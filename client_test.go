package gosail

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

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

func getTestConfig() APIConfig {
	c := APIConfig{}
	c.APIKey = "TestAPIKey"
	c.SecretKey = "TestSecretKey"
	c.BaseURL = "https://api.sailthru.com"
	return c
}

func TestCreateJob(t *testing.T) {
	expectedJobID := "555a21e5a6cba8e27427eb23"
	mc := NewMockClient(normalJob)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)

	fields := map[string]map[string]interface{}{"vars": {"user_id": 1}}

	resp, err := sc.CreateJob("export_list_data", "ad_hoc_test_list_1", fields, "json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if resp.JobID != expectedJobID {
		t.Errorf("Expected %v, got %v\n", expectedJobID, resp.JobID)
	}
}

func TestCreateInvalidJobType(t *testing.T) {
	expectedErrorStr := "Invalid jobType: invalid_job_type"
	mc := NewMockClient(normalJob)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)
	fields := map[string]map[string]interface{}{"vars": {"user_id": 1}}
	_, err := sc.CreateJob("invalid_job_type", "ad_hoc_test_list_1", fields, "json")
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
	mc := NewMockClient(normalJob)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)
	j, err := sc.GetJob(expectedJobID)
	if err != nil {
		t.Error(err)
	}
	if j.JobID != expectedJobID {
		t.Errorf("Expected %v, got %v\n", expectedJobID, j.JobID)
	}
}

func TestGetJobInvalidJSON(t *testing.T) {
	expectedJobID := "555a468b975910683a63b666"
	expectedErrorStr := "json: cannot unmarshal"
	mc := NewMockClient(invalidJSON)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)
	_, err := sc.GetJob(expectedJobID)
	if err == nil {
		t.Errorf("Expected %v, got %v\n", expectedErrorStr, nil)
	} else {
		if !strings.HasPrefix(err.Error(), expectedErrorStr) {
			t.Errorf("Expected prefix %v, not found in %v\n", expectedErrorStr, err.Error())
		}
	}
}

func TestGetJobExpired(t *testing.T) {
	expectedJobID := "555a21e5a6cba8e27427eb23"
	mc := NewMockClient(expiredJob)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)
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
	expectedJobID := "InvalidJobID"
	mc := NewMockClient(invalidJob)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)
	_, err := sc.GetJob(expectedJobID)
	if err != nil {
		if err.Error() != "Error Response: 401 Unauthorized" {
			t.Errorf("Expected 401 Unauthorized, got %v", err)
		}
	}
	if err == nil {
		t.Errorf("This should have returned an error, 401 Unauthorized")
	}
}

func TestGetNormalCSV(t *testing.T) {
	exportURL := "https://s3.amazonaws.com/sailthru/export/2015/05/19/UNIQUE_FILE_ID"
	expectedUserIDs := []int{10, 3, 5, 4, 6, 7, 8, 2, 1, 9}
	mc := NewMockClient(normalCSV)
	c := getTestConfig()
	sc := NewSailThruClient(&mc, c)
	r, err := sc.GetCSVData(exportURL)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error("Result data should not be nil")
	}
	data, readErr := ioutil.ReadAll(r)
	if readErr != nil {
		t.Error(readErr)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		t.Error("There should be at least 1 line returned")
	}
	topLine := lines[0]
	var userCol = -1
	for k, v := range strings.Split(topLine, ",") {
		if v == "userid" {
			userCol = k
			break
		}
	}
	if userCol < 0 {
		t.Error("`userid` was not returned in the CSV header row")
	}
	userIDs := []int{}
	for k, line := range strings.Split(string(data), "\n") {
		if k > 0 {
			cols := strings.Split(line, ",")
			if len(cols) >= userCol {
				userID, convErr := strconv.Atoi(cols[userCol])
				if convErr != nil {
					t.Error(convErr)
					break
				}
				userIDs = append(userIDs, userID)
			}
		}
	}
	if fmt.Sprintf("%v", userIDs) != fmt.Sprintf("%v", expectedUserIDs) {
		t.Errorf("Expected %v, got %v\n", expectedUserIDs, userIDs)
	}
}
