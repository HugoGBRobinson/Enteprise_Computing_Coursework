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

func QueryMicroservices(w http.ResponseWriter, r *http.Request) {
	STTResponse := getSTTResponse(r)
	alphaResponse := getAlphaResponse(STTResponse)
	TTSResponse := getTTSResponse(alphaResponse)
	respond(TTSResponse, w)
}

func requestWriter(post string, URI string, data []byte, client *http.Client) []byte {
	req, err := http.NewRequest(post, URI, bytes.NewReader(data))
	check(err)
	rsp, err2 := client.Do(req)
	check(err2)
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		return body
	}
	return nil
}

func getSTTResponse(r *http.Request) []byte {
	client := &http.Client{}
	t := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if speech, ok := t["speech"]; ok {
			speechForSTT := &Speech{Speech: speech}
			jsonSpeechForSTT, _ := json.Marshal(speechForSTT)
			response := requestWriter("POST", "http://localhost:3002/stt", jsonSpeechForSTT, client)
			return response
		}
	}

	return nil
}
func getAlphaResponse(STTResponse []byte) []byte {
	client := &http.Client{}
	response := requestWriter("POST", "http://localhost:3001/alpha", STTResponse, client)
	return response
}
func getTTSResponse(alphaResponse []byte) []byte {
	client := &http.Client{}
	response := requestWriter("POST", "http://localhost:3003/tts", alphaResponse, client)
	return response
}

func respond(TTSResponse []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	fmt.Print(string(TTSResponse))
	w.Write(TTSResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alexa", QueryMicroservices).Methods("POST")
	http.ListenAndServe(":3000", r)
}
