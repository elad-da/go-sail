###go-sail

`go-sail` is a client written in Go (golang) that communicates with the [SailThru API](https://api.sailthru.com).

To include in your project, `git clone` the repo to your $GOPATH.  Put the following in your import section:

```go
  import (
    "github.com/DramaFever/go-sail"
    )
```

You can instantiate a go-client like this:

```go
        httpClient := gosail.HTTPClient{}
				sc := gosail.NewSailThruClient(&httpClient, "YourAPIKey", "YourSecretKey")
```

To create a Job:

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

The `job` variable is an instace of the `Job` struct:

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

If the Status is complete and expired is `false`, then the  JobID can be used to download the data from the job:

```go
data, dataErr := sc.GetCSVData(job.JobID)
```

`GetCSVData` returns a `byte` slice that, if converted to a string, would look like data from a CSV file.

To run tests, run `go-test` inside the **go-sail** directory.
