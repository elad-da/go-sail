###go-sail

`go-sail` is a client written in Go (golang) that communicates with the [SailThru API](https://api.sailthru.com).

----------
Currently, the client only allows for:

 1. [Creating](http://getstarted.sailthru.com/new-for-developers-overview/reporting/job/#POST) a job
  1. Currently the only job type allowed is [`export_list_data`](http://getstarted.sailthru.com/new-for-developers-overview/reporting/job/#export_list_data)
 1. [Checking](http://getstarted.sailthru.com/new-for-developers-overview/reporting/job/#GET) of the job's status
 1. Return of the job's data once it has completed.  

More features to come.

----------
To include in your project, `git clone` the repo to your $GOPATH.  Put the following in the import section of the package that will make use of `go-sail`:

```go
  import (
    "github.com/DramaFever/go-sail"
    )
```

Instantiate a go-sail client:
----------

Setup the `APIConfig` struct:
 ```go
  c := APIConfig{}
	c.APIKey = "TestAPIKey"
	c.SecretKey = "TestSecretKey"
	c.BaseURL = "https://api.sailthru.com"
  ```

Then pass that to the `NewSailThruClient` along with a in instance of HTTPClient:

```go
httpClient := gosail.HTTPClient{}
sc := gosail.NewSailThruClient(&httpClient, c)
```


Create a Job:
----------

```go
resp, err := sc.CreateJob("export_list_data", "ad_hoc_test_list_1", "json")
```

The `resp` struct has the following properties:

```go
type CreateJobResponse struct {
	JobID  string `json:"job_id"`
	Name   string `json:"name"`
	List   string `json:"list"`
	Status string `json:"status"`
}
```

Use the jobID to check the status:
```go
job, jobErr := sc.GetJob(jobItem.JobID)
```

The `job` variable is an instance of the `Job` struct:

```go
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
```


Download the data:
----------
If the `Status` is `complete` and `expired` is `false`, then the `JobID` can be used to download the data from the job:

```go
data, dataErr := sc.GetCSVData(job.JobID)
```

`GetCSVData` returns an `io.ReadCloser` that, can be converted to a `slice` of `string` like so:

```go
//r is the returned io.ReadCloser from GetCSVData
data, readErr := ioutil.ReadAll(r)
if readErr != nil {
  //Handle the error as you see fit.
}
lines := strings.Split(string(data), "\n")
```


----------

To run tests, run `go test` inside the **go-sail** directory.
