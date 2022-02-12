package TTS

import (
	"Coursework/config"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
)

var KEY = config.GetAzureKey()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Speak struct {
	XMLName xml.Name `xml:"speak"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Lang    string   `xml:"xml:lang,attr"`
	Voice   Voice
}

type Voice struct {
	XMLName xml.Name `xml:"voice"`
	Text    string   `xml:",chardata"`
	Lang    string   `xml:"xml:lang,attr"`
	Name    string   `xml:"name,attr"`
}

type Speech struct {
	Speech []uint8 `json:"speech"`
}

func TextToSpeech(w http.ResponseWriter, r *http.Request) {
	t := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if cont, ok := t["contents"]; ok {
			client := &http.Client{}
			v := &Voice{
				XMLName: xml.Name{},
				Text:    cont,
				Lang:    "en-US",
				Name:    "en-US-JennyNeural",
			}
			s := &Speak{
				XMLName: xml.Name{},
				Text:    "",
				Version: "1.0",
				Lang:    "en-US",
				Voice:   *v,
			}

			m, _ := xml.MarshalIndent(s, "", "  ")
			req, err2 := http.NewRequest("POST", URI, bytes.NewReader(m))
			check(err2)
			req.Header.Set("Content-Type", "application/ssml+xml")
			req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
			req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")

			rsp, err3 := client.Do(req)
			check(err3)
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				body, err4 := ioutil.ReadAll(rsp.Body)
				check(err4)
				w.WriteHeader(http.StatusOK)
				speech := Speech{Speech: body}
				json.NewEncoder(w).Encode(speech)
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
	r.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3003", r)
}
