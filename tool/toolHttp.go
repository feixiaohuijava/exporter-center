package tool

import (
	"io"
	"io/ioutil"
	"net/http"
)

func Checkerr(e error) {
	if e != nil {
		panic(e)
	}
}

func HttpGetJson(url string, payload io.Reader, token string, params []map[string]string) (string, int) {
	req, err := http.NewRequest("GET", url, payload)
	if token != "" {
		req.Header.Set("Authorization", "JWT "+token)
	}
	if params != nil && len(params) > 0 {
		q := req.URL.Query()
		for _, param := range params {
			for k, v := range param {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	Checkerr(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	Checkerr(err)
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), statusCode
}

type AuthPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HttpPost(httpUrl string, payload io.Reader, authPost AuthPost) (string, int) {
	req, err := http.NewRequest("POST", httpUrl, payload)
	Checkerr(err)
	req.Header.Add("Content-Type", "application/json ;charset=utf-8")
	// add auth
	if (authPost != AuthPost{}) {
		req.SetBasicAuth(authPost.Username, authPost.Password)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	Checkerr(err)
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), statusCode
}
