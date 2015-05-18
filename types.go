package gosail

import "net/http"

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
