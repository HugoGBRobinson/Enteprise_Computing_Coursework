package STT

import (
	"Coursework/config"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
)

var KEY = config.GetAzureKey()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Text struct {
	Text string `json:"text"`
}

type Body struct {
	RecognitionStatus string
	DisplayText       string
	Offset            string
	Duration          string
}

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
func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", SpeechToText).Methods("POST")
	http.ListenAndServe(":3002", r)
}
