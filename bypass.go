package binglib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func Bypass(bypassServer, cookie, iframeid, IG, convId, rid string) (passResp PassResponseStruct, err error) {
	passRequest := passRequestStruct{
		Cookies:  cookie,
		Iframeid: iframeid,
		IG:       IG,
		ConvId:   convId,
		RId:      rid,
	}
	passResq, err := json.Marshal(passRequest)
	if err != nil {
		return passResp, err
	}

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, err := http.NewRequest("POST", bypassServer, bytes.NewReader(passResq))
	if err != nil {
		return passResp, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return passResp, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &passResp)
	if err != nil {
		return passResp, err
	}
	return passResp, nil
}
