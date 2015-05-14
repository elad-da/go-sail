package gosail

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

var apiURLGet = "https://api.sailthru.com/%v?json=%v&api_key=%v&sig=%v&format=%v"
var apiURLPost = "https://api.sailthru.com/%v?format=%v"

//SailThruClient Struct that contains key & hashing locations for sailthru calls
type SailThruClient struct {
	apiKey         string
	secretKey      string
	jsonhashstring string
}

//Job struct that contains json marshalled data about a sailthru Job
type Job struct {
	JobID     string `json:"job_id"`
	Name      string `json:"name"`
	List      string `json:"list"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Filename  string `json:"filename"`
	Expired   bool   `json:"expired"`
	ExportURL string `json:"export_url"`
}

//CreateJobResponse struct that contains the result of a Create Job Request.
type CreateJobResponse struct {
	JobID  string `json:"job_id"`
	Name   string `json:"name"`
	List   string `json:"list"`
	Status string `json:"status"`
}

type jSONBody struct {
	Body    string
	EscBody string
}

//NewSailThruClient func that creates a sailthruclient instance for calls to the SailThruAPI
func NewSailThruClient(apiKey string, secretKey string) SailThruClient {
	sc := SailThruClient{apiKey, secretKey, "%v%vjson%v"}
	return sc
}

func (sc *SailThruClient) getSignatureString(params map[string]string) string {
	stringtohash := ""
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		stringtohash += params[key]
	}
	return fmt.Sprintf(sc.jsonhashstring, sc.secretKey, sc.apiKey, stringtohash)
}

func (sc *SailThruClient) getSigHash(strToHash string) string {
	h := md5.New()
	io.WriteString(h, strToHash)
	sig := fmt.Sprintf("%x", h.Sum(nil))
	return sig
}

func (sc *SailThruClient) getJSONStringBody(items map[string]interface{}) string {
	jsonparams, _ := json.Marshal(items)
	return string(jsonparams)
}

func (sc *SailThruClient) getJSONBody(data map[string]interface{}) jSONBody {
	b := jSONBody{}
	b.Body = sc.getJSONStringBody(data)
	b.EscBody = url.QueryEscape(b.Body)
	return b
}

func (sc *SailThruClient) getSigForJSONBody(params map[string]string) string {
	str := sc.getSignatureString(params)
	hash := sc.getSigHash(str)
	return hash
}

func (sc *SailThruClient) getPostForm(items map[string]interface{}) url.Values {
	jsonb := sc.getJSONBody(items)
	data := map[string]string{"json": jsonb.Body}
	sig := sc.getSigForJSONBody(data)
	form := url.Values{}
	form.Set("api_key", sc.apiKey)
	form.Set("sig", sig)
	form.Set("json", jsonb.Body)
	form.Set("format", "json")
	return form
}

//CreateJob Func that creates a sailthru job.  Call must specify the type of job, the name of the list and the format of the returned data (json|xml)
//Keep in mind that CreateJob does not immediately return the contents of the job, it starts the job and returns a jobID.  The status of the job is checked via the GetJob func
func (sc *SailThruClient) CreateJob(jobType string, listName string, format string) (*CreateJobResponse, error) {
	r := CreateJobResponse{}
	posturl := fmt.Sprintf(apiURLPost, "job", format)
	items := map[string]interface{}{"job": jobType, "list": listName}
	form := sc.getPostForm(items)

	req, reqErr := http.NewRequest("POST", posturl, bytes.NewBufferString(form.Encode()))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}
	errJSON := json.Unmarshal(body, &r)
	if errJSON != nil {
		return nil, errJSON
	}
	return &r, nil
}

//GetJob Func that takes a jobID, which is returned by CreateJob and a format (json|xml) to get back the status of a CreateJob func call
func (sc *SailThruClient) GetJob(jobID string) (*Job, error) {
	items := map[string]interface{}{"job_id": jobID}
	jsonb := sc.getJSONBody(items)
	data := map[string]string{"json": jsonb.Body}
	sig := sc.getSigForJSONBody(data)
	apiurl := fmt.Sprintf(apiURLGet, "job", jsonb.EscBody, sc.apiKey, sig, "json")
	res, _ := http.Get(apiurl)
	output, _ := ioutil.ReadAll(res.Body)
	job := Job{}
	json.Unmarshal([]byte(output), &job)
	return &job, nil
}

//GetCSVData If the job has completed and it has not expired, this call will return the data in the CSV file the job created
func (sc *SailThruClient) GetCSVData(path string) ([]byte, error) {
	res, errGet := http.Get(path)
	if errGet != nil {
		return nil, errGet
	}
	return ioutil.ReadAll(res.Body)
}

//CreateJobAndReturnJob This will create the job, and then return the contents of the job, providing it does not timeout(value is seconds)
func (sc *SailThruClient) CreateJobAndReturnJob(jobType string, listName string, format string, timeout int) ([]byte, error) {
	cjresp, err := sc.CreateJob(jobType, listName, format)
	if err != nil {
		return nil, err
	}
	timer := time.Tick(100 * time.Millisecond)
	start := time.Now()
	for now := range timer {
		_ = now
		j, errJ := sc.GetJob(cjresp.JobID)
		if errJ != nil {
			return nil, errJ
		}
		if j.Status == "completed" && !j.Expired {
			return sc.GetCSVData(j.ExportURL)
		}
		delta := time.Now().Sub(start)
		if delta.Seconds() > float64(timeout) {
			break
		}
	}
	return nil, fmt.Errorf("Timeout Error - Job not ready after %v seconds\n", timeout)
}
