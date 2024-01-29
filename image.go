package binglib

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Harry-zklcdc/bing-lib/lib/request"
	"golang.org/x/net/html"
)

const (
	bingImageCreateUrl = "%s/images/create?q=%s&rt=4&FORM=GENCRE"
	bingImageResult    = "%s/images/create/async/results/%s"
)

func NewImage(cookies string) *Image {
	return &Image{
		cookies:     cookies,
		BingBaseUrl: bingBaseUrl,
	}
}

func (image *Image) Clone() *Image {
	return &Image{
		cookies:     image.cookies,
		BingBaseUrl: image.BingBaseUrl,
	}
}

func (image *Image) SetBingBaseUrl(bingBaseUrl string) *Image {
	image.BingBaseUrl = bingBaseUrl
	return image
}

func (image *Image) SetCookies(cookies string) *Image {
	image.cookies = cookies
	return image
}

func (image *Image) SetXFF(xff string) *Image {
	image.xff = xff
	return image
}

func (image *Image) SetBypassServer(bypassServer string) *Image {
	image.bypassServer = bypassServer
	return image
}

func (image *Image) GetBingBaseUrl() string {
	return image.BingBaseUrl
}

func (image *Image) GetCookies() string {
	return image.cookies
}

func (image *Image) GetXFF() string {
	return image.xff
}

func (image *Image) GetBypassServer() string {
	return image.bypassServer
}

func (image *Image) Image(q string) ([]string, string, error) {
	var res []string

	c := request.NewRequest()
	if image.xff != "" {
		c.SetHeader("X-Forwarded-For", image.xff)
	}
	c.Post().SetUrl(bingImageCreateUrl, image.BingBaseUrl, url.QueryEscape(q)).
		SetBody(strings.NewReader(url.QueryEscape(fmt.Sprintf("q=%s&qs=ds", q)))).
		SetContentType("application/x-www-form-urlencoded").
		SetHeader("Cookie", image.cookies).
		SetHeader("User-Agent", userAgent).
		SetHeader("Origin", "https://www.bing.com").
		SetHeader("Referer", "https://www.bing.com/images/create/").
		Do()
	if c.Result.Status != 302 {
		return res, "", fmt.Errorf("status code: %d", c.Result.Status)
	}

	u, _ := url.Parse(fmt.Sprintf("%s%s", image.BingBaseUrl, c.GetHeader("Location")))
	c.Get().SetUrl("%s%s", image.BingBaseUrl, c.GetHeader("Location")).Do()
	if c.Result.Status != 200 {
		return res, "", fmt.Errorf("status code: %d", c.Result.Status)
	}

	id := u.Query().Get("id")
	// fmt.Println(id)

	i := 0
	for i < 120 {
		time.Sleep(1 * time.Second)
		i++
		c.Get().SetUrl(bingImageResult, image.BingBaseUrl, id).Do()
		if len(c.GetBodyString()) > 1 && strings.Contains(c.GetHeader("Content-Type"), "text/html") {
			break
		}
	}

	if i >= 120 {
		return res, "", fmt.Errorf("timeout")
	}

	// fmt.Println(c.GetBodyString())
	body, err := html.Parse(strings.NewReader(c.GetBodyString()))
	if err != nil {
		return res, id, err
	}

	findImgs(body, &res)

	var tmp []string
	for i := range res {
		if !strings.Contains(res[i], "/rp/") {
			url, _ := url.Parse(res[i])
			url.RawQuery = ""
			tmp = append(tmp, url.String())
		}
	}

	return tmp, id, nil
}

func findImgs(n *html.Node, vals *[]string) {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				*vals = append(*vals, a.Val)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findImgs(c, vals)
	}
}
