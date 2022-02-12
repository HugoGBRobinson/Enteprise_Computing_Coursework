package Alpha

import (
	"Coursework/config"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	AppID = config.GetAlphaKey()
	URI   = "https://api.wolframalpha.com/v1/result?appid=" + AppID
)

func QueryWolframAlpha(text []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", URI, bytes.NewBuffer(text))
	check(err)

	q := req.URL.Query()
	q.Add("i", "population of france")

	req.URL.RawQuery = q.Encode()
	rsp, err2 := client.Do(req)
	check(err2)
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		return body, nil
	} else {
		return nil, errors.New("cannot convert text to speech")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {
	text, err := ioutil.ReadFile("wolf.xml")
	check(err)
	response, err := QueryWolframAlpha(text)
	check(err)
	print(string(response))
}
