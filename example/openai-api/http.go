package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	binglib "github.com/Harry-zklcdc/bing-lib"
	"github.com/Harry-zklcdc/bing-lib/lib/hex"
)

var (
	cookie        = os.Getenv("COOKIE")
	bingBaseUrl   = os.Getenv("BING_BASE_URL")
	sydneyBaseUrl = os.Getenv("SYDNEY_BASE_URL")
)

var STOPFLAG = "stop"

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resqB, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resq chatRequest
	json.Unmarshal(resqB, &resq)

	if resq.Model != binglib.BALANCED && resq.Model != binglib.BALANCED_OFFLINE && resq.Model != binglib.CREATIVE && resq.Model != binglib.CREATIVE_OFFLINE && resq.Model != binglib.PRECISE && resq.Model != binglib.PRECISE_OFFLINE {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid model"))
		log.Println(r.RemoteAddr, r.Method, r.URL, "400")
		return
	}
	chat := binglib.NewChat(cookie)
	err = chat.NewConversation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(r.RemoteAddr, r.Method, r.URL, "500")
		return
	}
	chat.SetStyle(resq.Model)
	if bingBaseUrl != "" {
		chat.SetBingBaseUrl(bingBaseUrl)
	}
	if sydneyBaseUrl != "" {
		chat.SetSydneyBaseUrl(sydneyBaseUrl)
	}

	prompt, msg := chat.MsgComposer(resq.Messages)
	resp := chatResponse{
		Id:                "chatcmpl-NewBing",
		Object:            "chat.completion.chunk",
		SystemFingerprint: hex.NewHex(12),
		Model:             resq.Model,
		Create:            time.Now().Unix(),
	}

	if resq.Stream {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		text := make(chan string)
		go chat.ChatStream(prompt, msg, text)
		var tmp string

		for {
			tmp = <-text
			resp.Choices = []choices{
				{
					Index: 0,
					Delta: binglib.Message{
						// Role:    "assistant",
						Content: tmp,
					},
				},
			}
			if tmp == "EOF" {
				resp.Choices[0].Delta.Content = ""
				resp.Choices[0].FinishReason = &STOPFLAG
				resData, err := json.Marshal(resp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(r.RemoteAddr, r.Method, r.URL, "500")
					return
				}
				w.Write([]byte("data: "))
				w.Write(resData)
				break
			}
			resData, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(r.RemoteAddr, r.Method, r.URL, "500")
				return
			}
			w.Write([]byte("data: "))
			w.Write(resData)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		}
	} else {
		text, err := chat.Chat(prompt, msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(r.RemoteAddr, r.Method, r.URL, "500")
			return
		}

		resp.Choices = append(resp.Choices, choices{
			Index: 0,
			Message: binglib.Message{
				Role:    "assistant",
				Content: text,
			},
			FinishReason: &STOPFLAG,
		})

		resData, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(r.RemoteAddr, r.Method, r.URL, "500")
			return
		}
		w.Write(resData)
	}
	log.Println(r.RemoteAddr, r.Method, r.URL, "200")

}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println(r.RemoteAddr, r.Method, r.URL, "500")
		return
	}

	resqB, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(r.RemoteAddr, r.Method, r.URL, "500")
		return
	}

	var resq imageRequest
	json.Unmarshal(resqB, &resq)

	image := binglib.NewImage(cookie)
	imgs, _, err := image.Image(resq.Prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(r.RemoteAddr, r.Method, r.URL, "500")
		return
	}

	resp := imageResponse{
		Created: time.Now().Unix(),
	}
	for _, img := range imgs {
		resp.Data = append(resp.Data, imageData{
			Url: img,
		})
	}

	resData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(r.RemoteAddr, r.Method, r.URL, "500")
		return
	}
	w.Write(resData)
	log.Println(r.RemoteAddr, r.Method, r.URL, "200")
}
