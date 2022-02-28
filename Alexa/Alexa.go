package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Speech provides the necessary structure to create a json speech
//response.
type Speech struct {
	Speech string `json:"speech"`
}

// QueryMicroservices accesses all the other microservices to create
// link them all together, allowing a user to send and receive a
// response from a single microservice
func QueryMicroservices(w http.ResponseWriter, r *http.Request) {
	STTResponse := getSTTResponse(w, r)
	if STTResponse != nil {
		alphaResponse := getAlphaResponse(w, STTResponse)
		if alphaResponse != nil {
			TTSResponse := getTTSResponse(w, alphaResponse)
			if TTSResponse != nil {
				respond(TTSResponse, w)
			}
		}
	}
}

// requestWriter is a general function used to write a POST request
// when provided with the URI and data to be sent.
func requestWriter(w http.ResponseWriter, URI string, data []byte) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("POST", URI, bytes.NewReader(data))
	check(err)
	rsp, err2 := client.Do(req)
	check(err2)
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		return body
		// If the request fails then a response is sent back to the user with the status code
		// and the URI of the failed request
	} else if rsp.StatusCode == http.StatusNotFound {
		w.WriteHeader(rsp.StatusCode)
		w.Write([]byte("Status Not Found response from " + rsp.Request.RequestURI))
	} else if rsp.StatusCode == http.StatusBadRequest {
		w.WriteHeader(rsp.StatusCode)
		w.Write([]byte("Status Bad Request response from " + rsp.Request.RequestURI))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown error response from " + rsp.Request.RequestURI))
	}
	return nil
}

// getSTTResponse gets the response form the STT microservice
func getSTTResponse(w http.ResponseWriter, r *http.Request) []byte {
	t := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if speech, ok := t["speech"]; ok {
			speechForSTT := &Speech{Speech: speech}
			jsonSpeechForSTT, _ := json.Marshal(speechForSTT)
			response := requestWriter(w, "http://localhost:3002/stt", jsonSpeechForSTT)
			if response != nil {
				return response
			} else {
				return nil
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	return nil
}

// getAlphaResponse gets the response from the Alpha microservice
func getAlphaResponse(w http.ResponseWriter, STTResponse []byte) []byte {
	response := requestWriter(w, "http://localhost:3001/alpha", STTResponse)
	if response != nil {
		return response
	} else {
		return nil
	}

}

// getTTSResponse gets the response from the TTS microservice
func getTTSResponse(w http.ResponseWriter, alphaResponse []byte) []byte {
	response := requestWriter(w, "http://localhost:3003/tts", alphaResponse)
	if response != nil {
		return response
	} else {
		return nil
	}

}

// respond is used to send the final .wav response back to the user and
// error handling for this response is done throughout the other functions
func respond(TTSResponse []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(TTSResponse)
}

//main sets up the listen and serve functionality allowing a user to
//request its services.
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alexa", QueryMicroservices).Methods("POST")
	http.ListenAndServe(":3000", r)
}
