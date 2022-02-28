package main

import (
	"Coursework/config"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//REGION and URI provide constant a URI to connect to
//Microsoft Cognitive Services.
const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
)

//KEY is the provided access key, it is obtained through the
//config file.
var KEY = config.GetAzureKey()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Text provides the necessary structure to create a json text
//response.
type Text struct {
	Text string `json:"text"`
}

//Body provides the necessary structure to unmarshall the xml
//response from Azure
type Body struct {
	RecognitionStatus string
	DisplayText       string
	Offset            string
	Duration          string
}

//SpeechToText is the primary function of this file as it takes in
//the request and response writer, decodes the json request, decodes
//the speech, marshals a xml response and encodes a new json
//response back to Alexa.
func SpeechToText(w http.ResponseWriter, r *http.Request) {
	t := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if speech, ok := t["speech"]; ok {
			client := &http.Client{}
			decoded, _ := b64.StdEncoding.DecodeString(speech)
			req, err := http.NewRequest("POST", URI, bytes.NewReader(decoded))
			check(err)
			req.Header.Set("Content-Type",
				"audio/wav;codecs=audio/pcm;samplerate=16000")
			req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
			rsp, err2 := client.Do(req)
			check(err2)
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				body, err3 := ioutil.ReadAll(rsp.Body)
				check(err3)
				w.WriteHeader(http.StatusOK)
				var UMBody Body
				json.Unmarshal(body, &UMBody)
				text := Text{Text: UMBody.DisplayText}
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

//main sets up the listen and serve functionality allowing Alexa to
//request its services.
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", SpeechToText).Methods("POST")
	http.ListenAndServe(":3002", r)
}
