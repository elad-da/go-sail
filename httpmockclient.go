package gosail

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type mockType int

const (
	NormalJob mockType = 1 + iota
	ExpiredJob
	InvalidJob
	NormalCSV
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

var mockTypes = [...]string{
	"Normal Job",
	"Expired Job",
	"Invalid Job",
	"Normal CSV",
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
	case 4:
		mc.doFunc = doNormal
		mc.getFunc = getNormalCSV
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
	resp.StatusCode = 200
	resp.Status = "200 OK"
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
	resp.StatusCode = 200
	resp.Status = "200 OK"
	respString := `{"job_id":"555a468b975910683a63b666","name":"Export All List Data: ad_hoc_test_list_1","list":"ad_hoc_test_list_1","status":"completed","start_time":"Mon, 18 May 2015 16:07:39 -0400","end_time":"Mon, 18 May 2015 16:07:40 -0400","filename":"ad_hoc_test_list_1.csv","export_url":"https:\/\/s3.amazonaws.com\/sailthru\/export\/2015\/05\/18\/4039cfa8f1d782f3af77b46388b55a5b"}`
	resp.Body = nopCloser{bytes.NewBufferString(respString)}
	return &resp, nil
}

func getNormalCSV(url string) (*http.Response, error) {
	resp := http.Response{}
	resp.StatusCode = 200
	resp.Status = "200 OK"
	respString := `"Profile Id","Email Hash",Domain,Engagement,Lists,"Profile Created Date",Signup,Opens,Clicks,Pageviews,"Last Open","Last Click","Last Pageview","Optout Time","List Signup","Geolocation City","Geolocation State","Geolocation Country","Geolocation Zip","Lifetime Message","First Purchase Time","Purchase Count","Purchase Price","Purchase Incomplete","Last Purchase Time","Largest Purchase Item Price","Top Device","Email Status",userid
554bb7153b35d0732c8c0e8a,14b5aebbfaf84afa184df9b67983cb04,greetings.org,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:49","2015/05/07 15:03:49",0,0,0,,,,,"2015/05/07 15:03:49",,,,,0,,0,,0,,,,Valid,10
554bb7143b35d0732c8c0e83,332a70d9324e29a435652f302b1e39fc,noun.com,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,3
554bb7143b35d0732c8c0e85,0d2c4a6d6ea19c4565485c5b70286268,color.edu,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,5
554bb7143b35d0732c8c0e84,67c8f20accfd0cb3821768dcdde74bf7,othernown.org,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,4
554bb7143b35d0732c8c0e86,49676844f8ad384333e947c119100906,color.edu,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,6
554bb7143b35d0732c8c0e87,e1ab9ec1da11867bd6d0fa88f4dcd403,state.gov,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,7
554bb7143b35d0732c8c0e88,6de2e68af618f5965983a045ad048392,state.gov,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,8
554bb7143b35d0732c8c0e82,88725c8e20a171d671e699e0294680e1,that.com,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:48","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,2
554bb7133b35d0732c8c0e81,a2b20ec1c29dca3c775ca49f57379e0e,this.com,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:47","2015/05/07 15:03:48",0,0,0,,,,,"2015/05/07 15:03:48",,,,,0,,0,,0,,,,Valid,1
554bb7153b35d0732c8c0e89,259a2011c9e81d93a850db60ead1af34,thechoseone.com,dormant,ad_hoc_test_list_1,"2015/05/07 15:03:49","2015/05/07 15:03:49",0,0,0,,,,,"2015/05/07 15:03:49",,,,,0,,0,,0,,,,Valid,9`
	resp.Body = nopCloser{bytes.NewBufferString(respString)}
	return &resp, nil
}

func getNormalExpiredJob(url string) (*http.Response, error) {
	resp := http.Response{}
	resp.StatusCode = 200
	resp.Status = "200 OK"
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
