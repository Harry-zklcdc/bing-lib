package binglib

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Harry-zklcdc/bing-lib/lib/request"
)

func Bypass(bypassServer, cookie, iframeid, IG, convId, rid, T string) (passResp PassResponseStruct, status int, err error) {
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
	c.SetMethod("POST").SetUrl(bypassServer).SetBody(bytes.NewReader(passResq)).SetHeader("Content-Type", "application/json").SetHeader("User-Agent", userAgent).Do()

	err = json.Unmarshal(c.Result.Body, &passResp)
	if err != nil {
		return passResp, http.StatusInternalServerError, err
	}
	return passResp, c.Result.Status, nil
}
