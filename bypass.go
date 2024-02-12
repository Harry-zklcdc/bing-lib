package binglib

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Harry-zklcdc/bing-lib/lib/aes"
	"github.com/Harry-zklcdc/bing-lib/lib/request"
)

func Bypass(bypassServer, cookie, iframeid, IG, convId, rid string) (passResp PassResponseStruct, status int, err error) {
	if IG == "" || len(IG) < 32 {
		return passResp, http.StatusBadRequest, errors.New("IG too short")
	}
	T, err := aes.Encrypt("Harry-zklcdc/go-proxy-bingai", IG)
	if err != nil {
		return passResp, http.StatusInternalServerError, err
	}
	passRequest := passRequestStruct{
		Cookies:  cookie,
		Iframeid: iframeid,
		IG:       IG,
		ConvId:   convId,
		RId:      rid,
		T:        T,
	}
	passResq, err := json.Marshal(passRequest)
	if err != nil {
		return passResp, http.StatusInternalServerError, err
	}

	c := request.NewRequest()
	c.SetMethod("POST").SetBody(bytes.NewReader(passResq)).SetHeader("Content-Type", "application/json").SetHeader("User-Agent", userAgent).Do()

	err = json.Unmarshal(c.Result.Body, &passResp)
	if err != nil {
		return passResp, http.StatusInternalServerError, err
	}
	return passResp, c.Result.Status, nil
}
