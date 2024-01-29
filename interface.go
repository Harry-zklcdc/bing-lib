package binglib

import (
	"time"

	"github.com/google/uuid"
)

const (
	bingBaseUrl   = "https://www.bing.com"
	sydneyBaseUrl = "wss://sydney.bing.com"

	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
)

type Chat struct {
	cookies       string
	xff           string // X-Forwarded-For Header
	bypassServer  string
	chatHub       *ChatHub
	BingBaseUrl   string
	SydneyBaseUrl string
}

type Image struct {
	cookies      string
	xff          string // X-Forwarded-For Header
	bypassServer string
	BingBaseUrl  string
}

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatHub struct {
	style   string
	chatReq ChatReq
}

type ChatReq struct {
	ConversationId                 string `json:"conversationId"`
	ClientId                       string `json:"clientId"`
	ConversationSignature          string `json:"conversationSignature"`
	EncryptedConversationSignature string `json:"encryptedconversationsignature"`
}

type SystemContext struct {
	Author      string `json:"author"`
	Description string `json:"description"`
	ContextType string `json:"contextType"`
	MessageType string `json:"messageType"`
	MessageId   string `json:"messageId"`
}

type Plugins struct {
	Id string `json:"id"`
}

type ResponsePayload struct {
	Type         int    `json:"type"`
	Target       string `json:"target"`
	InvocationId int    `json:"invocationId,string"`
	Arguments    []struct {
		Messages []struct {
			Text   string `json:"text"`
			Author string `json:"author"`
			From   struct {
				Id   string `json:"id"`
				Name any    `json:"name"`
			} `json:"from"`
			CreatedAt     time.Time `json:"createdAt"`
			Timestamp     time.Time `json:"timestamp"`
			Locale        string    `json:"locale"`
			Market        string    `json:"market"`
			Region        string    `json:"region"`
			Location      string    `json:"location"`
			LocationHints []struct {
				Country           string `json:"country"`
				CountryConfidence int    `json:"countryConfidence"`
				State             string `json:"state"`
				City              string `json:"city"`
				CityConfidence    int    `json:"cityConfidence"`
				ZipCode           string `json:"zipCode"`
				TimeZoneOffset    int    `json:"timeZoneOffset"`
				Dma               int    `json:"dma"`
				SourceType        int    `json:"sourceType"`
				Center            struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
					Height    any     `json:"height"`
				} `json:"center"`
				RegionType int `json:"regionType"`
			} `json:"locationHints"`
			MessageId uuid.UUID `json:"messageId"`
			RequestId uuid.UUID `json:"requestId"`
			Offense   string    `json:"offense"`
			Feedback  struct {
				Tag       any    `json:"tag"`
				UpdatedOn any    `json:"updatedOn"`
				Type      string `json:"type"`
			} `json:"feedback"`
			ContentOrigin string `json:"contentOrigin"`
			Privacy       any    `json:"privacy"`
			InputMethod   string `json:"inputMethod"`
			HiddenText    string `json:"hiddenText"`
			MessageType   string `json:"messageType"`
			AdaptiveCards []struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Body    []struct {
					Type    string `json:"type"`
					Version string `json:"version"`
					Body    []struct {
						Type string `json:"type"`
						Text string `json:"text"`
						Wrap bool   `json:"wrap"`
					}
				} `json:"body"`
			} `json:"adaptiveCards"`
			SourceAttributions []struct {
				ProviderDisplayName string `json:"providerDisplayName"`
				SeeMoreUrl          string `json:"seeMoreUrl"`
				SearchQuery         string `json:"searchQuery"`
			} `json:"sourceAttributions"`
			SuggestedResponses []struct {
				Text        string    `json:"text"`
				Author      string    `json:"author"`
				CreatedAt   time.Time `json:"createdAt"`
				Timestamp   time.Time `json:"timestamp"`
				MessageId   string    `json:"messageId"`
				MessageType string    `json:"messageType"`
				Offense     string    `json:"offense"`
				Feedback    struct {
					Tag       any    `json:"tag"`
					UpdatedOn any    `json:"updatedOn"`
					Type      string `json:"type"`
				} `json:"feedback"`
				ContentOrigin string `json:"contentOrigin"`
				Privacy       any    `json:"privacy"`
			} `json:"suggestedResponses"`
			SpokenText string `json:"spokenText"`
		} `json:"messages"`
		FirstNewMessageIndex   int       `json:"firstNewMessageIndex"`
		SuggestedResponses     any       `json:"suggestedResponses"`
		ConversationId         string    `json:"conversationId"`
		RequestId              string    `json:"requestId"`
		ConversationExpiryTime time.Time `json:"conversationExpiryTime"`
		Telemetry              struct {
			Metrics   any       `json:"metrics"`
			StartTime time.Time `json:"startTime"`
		} `json:"telemetry"`
		ShouldInitiateConversation bool `json:"shouldInitiateConversation"`
		Result                     struct {
			Value          string `json:"value"`
			Message        string `json:"message"`
			ServiceVersion string `json:"serviceVersion"`
		} `json:"result"`
	} `json:"arguments,omitempty"`
	Item struct {
		Result struct {
			Value          string `json:"value"`
			Message        string `json:"message"`
			ServiceVersion string `json:"serviceVersion"`
		} `json:"result"`
	} `json:"item,omitempty"`
}

type passRequestStruct struct {
	IG       string `json:"IG,omitempty"`
	Cookies  string `json:"cookies"`
	Iframeid string `json:"iframeid,omitempty"`
	ConvId   string `json:"convId,omitempty"`
	RId      string `json:"rid,omitempty"`
}

type PassResponseStruct struct {
	Result struct {
		Cookies    string `json:"cookies"`
		ScreenShot string `json:"screenshot"`
	} `json:"result"`
	Error string `json:"error"`
}
