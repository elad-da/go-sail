package gosail

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

//HTTPClienter : Interface that contains actual http calls to sailthru interface, or mocking implementation.
type HTTPClienter interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

//HTTPClient : Struct implementation for live SailhThru calls.
type HTTPClient struct {
}

//Do : Func that creates an HttpClient and calls live Sailthru API
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

//Get : Func that calls HTTP get calls against live SailThru API
func (c *HTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type mockType int

const (
	NormalJob mockType = 1 + iota
	ExpiredJob
	InvalidJob
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

var mockTypes = [...]string{
	"Normal Job",
	"Expired Job",
	"Invalid Job",
}

func (mt mockType) String() string {
	return mockTypes[mt-1]
}

//MockClient : Struct implementation for mock SailhThru calls.
type MockClient struct {
	Mocktype mockType
	doFunc
	getFunc
}

//Do : Func that creates an HttpClient and calls mock Sailthru API
func (mc *MockClient) Do(req *http.Request) (*http.Response, error) {
	return mc.doFunc(req)
}

//Get : Func that calls HTTP get against mock SailThru API
func (mc *MockClient) Get(url string) (*http.Response, error) {
	return mc.getFunc(url)
}

//Do : Function type that allows different http client results to be returned to the mock implementation
type doFunc func(req *http.Request) (*http.Response, error)

//Get : Function type that allows different http get results to be returned to the mock implementation
type getFunc func(url string) (*http.Response, error)

//NewMockClient : Returns a MockClient based on mockType parameter
func NewMockClient(mt mockType) MockClient {
	mc := MockClient{}
	mc.Mocktype = mt
	switch mt {
	case 1:
		mc.doFunc = doNormal
		mc.getFunc = getNormalDownloadLink
	case 2:
		mc.doFunc = doNormal
		mc.getFunc = getNormalExpiredJob
	case 3:
		mc.doFunc = doNormal
		mc.getFunc = getInvalidJob
	default:
		panic("Woah!  This type doesn't work!")
	}
	return mc
}

func getMapFromJSONForm(req *http.Request) (map[string]string, error) {
	jsonMap := make(map[string]string)
	errJSON := json.Unmarshal([]byte(req.FormValue("json")), &jsonMap)
	if errJSON != nil {
		return nil, errJSON
	}
	return jsonMap, nil
}

func doNormal(req *http.Request) (*http.Response, error) {
	resp := http.Response{}
	_, err := getMapFromJSONForm(req)
	if err != nil {
		return nil, err
	}
	respString := `{"job_id":"555a21e5a6cba8e27427eb23","name":"Export All List Data: ad_hoc_test_list_1","list":"ad_hoc_test_list_1","status":"pending"}`
	resp.Body = nopCloser{bytes.NewBufferString(respString)}
	return &resp, nil
}

func getNormalDownloadLink(url string) (*http.Response, error) {
	resp := http.Response{}
	respString := `{"job_id":"555a468b975910683a63b666","name":"Export All List Data: ad_hoc_test_list_1","list":"ad_hoc_test_list_1","status":"completed","start_time":"Mon, 18 May 2015 16:07:39 -0400","end_time":"Mon, 18 May 2015 16:07:40 -0400","filename":"ad_hoc_test_list_1.csv","export_url":"https:\/\/s3.amazonaws.com\/sailthru\/export\/2015\/05\/18\/4039cfa8f1d782f3af77b46388b55a5b"}`
	resp.Body = nopCloser{bytes.NewBufferString(respString)}
	return &resp, nil
}

func getNormalExpiredJob(url string) (*http.Response, error) {
	resp := http.Response{}
	respString := `{"job_id":"555a21e5a6cba8e27427eb23","name":"Export All List Data: ad_hoc_test_list_1","list":"ad_hoc_test_list_1","status":"completed","start_time":"Mon, 18 May 2015 13:31:17 -0400","end_time":"Mon, 18 May 2015 13:31:18 -0400","filename":"ad_hoc_test_list_1.csv","expired":true}`
	resp.Body = nopCloser{bytes.NewBufferString(respString)}
	return &resp, nil
}

func getInvalidJob(url string) (*http.Response, error) {
	resp := http.Response{}
	resp.StatusCode = 401
	resp.Status = "401 Unauthorized"
	respString := `{"error" : 99,"errormsg" : "Invalid Job ID: 555a468b975910683a63b667"}`
	resp.Body = nopCloser{bytes.NewBufferString(respString)}
	return &resp, nil
}
