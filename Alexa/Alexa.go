package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Speech struct {
	Speech string `json:"speech"`
}

type Text struct {
	Text string `json:"text"`
}

func QueryMicroservices(w http.ResponseWriter, r *http.Request) {
	STTResponse := getSTTResponse(r)
	//alphaResponse := getAlphaResponse(STTResponse)
	//TTSResponse := getTTSResponse(alphaResponse)
	fmt.Print(STTResponse)
}
func getSTTResponse(r *http.Request) []byte {
	client := &http.Client{}
	t := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if speech, ok := t["speech"]; ok {
			speechForSTT := &Speech{Speech: speech}
			jsonSpeechForSTT, _ := json.Marshal(speechForSTT)
			req, err := http.NewRequest("POST", "http://localhost:3002/stt", bytes.NewReader(jsonSpeechForSTT))
			check(err)
			rsp, err2 := client.Do(req)
			check(err2)
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				body, err3 := ioutil.ReadAll(rsp.Body)
				check(err3)
				fmt.Print(string(body))
				return body
			}
			return nil
		}
	}

	return nil
}
func getAlphaResponse(STTResponse string) {}
func getTTSResponse(alphaResponse string) {}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alexa", QueryMicroservices).Methods("POST")
	http.ListenAndServe(":3000", r)
}
