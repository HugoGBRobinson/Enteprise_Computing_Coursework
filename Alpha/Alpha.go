package main

import (
	"Coursework/config"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var (
	AppID = config.GetAlphaKey()
	URI   = "https://api.wolframalpha.com/v1/result?appid=" + AppID
)

type Text struct {
	Text string `json:"text"`
}

func QueryWolframAlpha(w http.ResponseWriter, r *http.Request) {
	t := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if text, ok := t["text"]; ok {
			client := &http.Client{}
			req, err := http.NewRequest("POST", URI, bytes.NewBuffer([]byte(text)))
			check(err)

			q := req.URL.Query()
			q.Add("i", text)

			req.URL.RawQuery = q.Encode()
			rsp, err2 := client.Do(req)
			check(err2)
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				body, err3 := ioutil.ReadAll(rsp.Body)
				check(err3)
				text := Text{Text: string(body)}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(text)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alpha", QueryWolframAlpha).Methods("POST")
	http.ListenAndServe(":3001", r)
}
