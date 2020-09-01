// gogtmetrix helps you use the GTmetrix API from Golang
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Base url for GTmetrix
const base_url string = "https://gtmetrix.com/api/0.1/"

// Use your username and password to get an authenticated client instance to communicate with the API
func GetClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
		baseURL:  base_url,
	}
}

// An authenticated client instance
type Client struct {
	Username string
	Password string
	baseURL  string
}

// Sends a URL of a website to GTmetrix to be tested. Returns a *TestRefference to get the status of the test
// and test results
func (c *Client) Test(siteURL string) (*TestRefference, error) {
	httpClient := &http.Client{}

	form := url.Values{}
	form.Add("url", siteURL)

	req, err := http.NewRequest("POST", c.baseURL+"test", strings.NewReader(form.Encode()))

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Println(err)
		return &TestRefference{}, err
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	testRefference := &TestRefference{}

	json.Unmarshal(body, testRefference)

	log.Println("Tested " + siteURL + ". Refference: " + testRefference.TestID)
	fmt.Println(testRefference)

	if testRefference.Error != "" {
		return testRefference, errors.New(testRefference.Error)
	}

	return testRefference, nil
}

// TestRefference can be used to poll for results of a queued website test
type TestRefference struct {
	TestID       string `json:"test_id"`
	PollStateURL string `json:"poll_state_url"`
	CreditsLeft  int    `json:"credits_left"`
	Error        string `json:"error"`
}

// Try to fetch results of a given test. It is advised you do this once per second at the most.
func (c *Client) PollResults(tr *TestRefference) (*TestModel, error) {
	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", c.baseURL+"test/"+tr.TestID, nil)

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Println(err)
		return &TestModel{}, err
	} else {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	testModel := &TestModel{}

	json.Unmarshal(body, testModel)

	return testModel, nil
}

// Automatically polls the test results for you in one second intervals until results are ready. Gives up after five minutes
func (c *Client) WaitForResults(tr *TestRefference) (*TestModel, error) {
	startTime := time.Now()

	for time.Since(startTime) < time.Minute*5 {
		result, err := c.PollResults(tr)

		if err != nil {
			return result, err
		}

		log.Println(result)

		if result.State == "completed" || result.State == "error" {
			return result, nil
		}

		time.Sleep(1 * time.Second)
	}

	return &TestModel{}, errors.New("Waited for results for too long. Gave up.")
}

// Accepts a website URL and returns a finished test. Handles polling for you, gives up after 5 minutes.
func (c *Client) TestAndWaitForResults(url string) (*TestModel, error) {
	tr, err := c.Test(url)

	if err != nil {
		return &TestModel{}, err
	}

	return c.WaitForResults(tr)
}

// The model holding results of a test, or the error returned
type TestModel struct {
	State     string        `json:"state"`
	Error     string        `json:"error"`
	Results   TestResults   `json:"results"`
	Resources TestResources `json:"resources"`
}

// The results of a test to be found in the TestModel struct
type TestResults struct {
	ReportURL                string `json:"report_url"`
	PagespeedScore           int    `json:"pagespeed_score"`
	YslowScore               int    `json:"yslow_score"`
	HtmlBytes                int    `json:"html_bytes"`
	HtmlLoadTime             int    `json:"html_load_time"`
	PageBytes                int    `json:"page_bytes"`
	PageLoadTime             int    `json:"page_load_time"`
	PageElements             int    `json:"page_elements"`
	RedirectDuration         int    `json:"redirect_duration"`
	ConnectDuration          int    `json:"connect_duration"`
	BackendDuration          int    `json:"backend_duration"`
	FirstPaintTime           int    `json:"first_paint_time"`
	FirstContentfulPaintTime int    `json:"first_contentful_paint_time"`
	DomInteractiveTime       int    `json:"dom_interactive_time"`
	DomContentLoadedTime     int    `json:"dom_content_loaded_time"`
	DomContentLoadedDuration int    `json:"dom_content_loaded_duration"`
	OnloadTime               int    `json:"onload_time"`
	OnloadDuration           int    `json:"onload_duration"`
	FullyLoadedTime          int    `json:"fully_loaded_time"`
	RUMSpeedIndex            int    `json:"rum_speed_index"`
}

// Additional resources for a successfull test
type TestResources struct {
	Screenshot     string `json:"screenshot"`
	HAR            string `json:"har"`
	Pagespeed      string `json:"pagespeed"`
	PagespeedFiles string `json:"pagespeed_files"`
	Yslow          string `json:"yslow"`
	ReportPDF      string `json:"report_pdf"`
	ReportPDFFull  string `json:"report_pdf_full"`
	Video          string `json:"video"`
	Filmstrip      string `json:"filmstrip"`
}
