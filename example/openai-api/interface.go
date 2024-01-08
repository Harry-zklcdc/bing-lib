package main

import binglib "github.com/Harry-zklcdc/bing-lib"

type chatRequest struct {
	Messages []binglib.Message `json:"messages"`
	Model    string            `json:"model"`
	Stream   bool              `json:"stream"`
}

type chatResponse struct {
	Id                string    `json:"id"`
	Object            string    `json:"object"`
	Create            int       `json:"created"`
	Model             string    `json:"model"`
	SystemFingerprint string    `json:"system_fingerprint"`
	Choices           []choices `json:"choices"`
}

type choices struct {
	Index        int               `json:"index"`
	Delta        []binglib.Message `json:"delta,omitempty"`
	Message      binglib.Message   `json:"message,omitempty"`
	Logprobs     string            `json:"logprobs,omitempty"`
	FinishReason string            `json:"finish_reason,omitempty"`
}

type imageRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
	N      int    `json:"n"`
}

type imageResponse struct {
	Created int         `json:"created"`
	Data    []imageData `json:"data"`
}

type imageData struct {
	Url string `json:"url"`
}
