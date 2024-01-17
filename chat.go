package binglib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Harry-zklcdc/bing-lib/lib/hex"
	"github.com/Harry-zklcdc/bing-lib/lib/request"
	"github.com/gorilla/websocket"
)

const (
	PRECISE          = "Precise"          // 精准
	BALANCED         = "Balanced"         // 平衡
	CREATIVE         = "Creative"         // 创造
	PRECISE_OFFLINE  = "Precise-offline"  // 精准, 不联网搜索
	BALANCED_OFFLINE = "Balanced-offline" // 平衡, 不联网搜索
	CREATIVE_OFFLINE = "Creative-offline" // 创造, 不联网搜索

	PRECISE_G4T          = "Precise-g4t"          // 精准 GPT4-Turbo
	BALANCED_G4T         = "Balanced-g4t"         // 平衡 GPT4-Turbo
	CREATIVE_G4T         = "Creative-g4t"         // 创造 GPT4-Turbo
	PRECISE_G4T_OFFLINE  = "Precise-g4t-offline"  // 精准 GPT4-Turbo, 不联网搜索
	BALANCED_G4T_OFFLINE = "Balanced-g4t-offline" // 平衡 GPT4-Turbo, 不联网搜索
	CREATIVE_G4T_OFFLINE = "Creative-g4t-offline" // 创造 GPT4-Turbo, 不联网搜索
)

const (
	bingCreateConversationUrl = "%s/turing/conversation/create?bundleVersion=1.1467.6"
	sydneyChatHubUrl          = "%s/sydney/ChatHub?sec_access_token=%s"

	spilt = "\x1e"
)

func NewChat(cookies string) *Chat {
	return &Chat{
		cookies:       cookies,
		BingBaseUrl:   bingBaseUrl,
		SydneyBaseUrl: sydneyBaseUrl,
	}
}

func (chat *Chat) Clone() *Chat {
	return &Chat{
		cookies:       chat.cookies,
		BingBaseUrl:   chat.BingBaseUrl,
		SydneyBaseUrl: chat.SydneyBaseUrl,
	}
}

func (chat *Chat) SetCookies(cookies string) *Chat {
	chat.cookies = cookies
	return chat
}

func (chat *Chat) SetStyle(style string) *Chat {
	chat.GetChatHub().SetStyle(style)
	return chat
}

func (chat *Chat) SetBingBaseUrl(bingBaseUrl string) *Chat {
	chat.BingBaseUrl = bingBaseUrl
	return chat
}

func (chat *Chat) SetSydneyBaseUrl(sydneyBaseUrl string) *Chat {
	chat.SydneyBaseUrl = sydneyBaseUrl
	return chat
}

func (chat *Chat) GetCookies() string {
	return chat.cookies
}

func (chat *Chat) GetChatHub() *ChatHub {
	return chat.chatHub
}

func (chat *Chat) GetStyle() string {
	return chat.GetChatHub().GetStyle()
}

func (chat *Chat) GetBingBaseUrl() string {
	return chat.BingBaseUrl
}

func (chat *Chat) GetSydneyBaseUrl() string {
	return chat.SydneyBaseUrl
}

func (chat *Chat) NewConversation() error {
	c := request.NewRequest()
	c.SetUrl(fmt.Sprintf(bingCreateConversationUrl, chat.BingBaseUrl)).
		SetHeader("Cookie", chat.cookies).
		SetHeader("Origin", "https://www.bing.com").
		SetHeader("Referer", "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx&wlexpsignin=1&wlexpsignin=1").
		SetHeader("User-Agent", userAgent).
		SetHeader("X-Ms-Useragent", "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.12.0 OS/Windows").
		SetHeader("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Microsoft Edge\";v=\"120\"").
		SetHeader("Sec-Ch-Ua-Arch", "\"x86\"").
		SetHeader("Sec-Ch-Ua-Bitness", "\"64\"").
		SetHeader("Sec-Ch-Ua-Full-Version", "\"120.0.2210.133\"").
		SetHeader("Sec-Ch-Ua-Full-Version-List", "\"Not_A Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"120.0.6099.217\", \"Microsoft Edge\";v=\"120.0.2210.133\"").
		SetHeader("Sec-Ch-Ua-Mobile", "?0").
		SetHeader("Sec-Ch-Ua-Model", "\"\"").
		SetHeader("Sec-Ch-Ua-Platform", "\"Windows\"").
		SetHeader("Sec-Ch-Ua-Platform-Version", "\"15.0.0\"").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("Sec-Gpc", "CEC540850250EBE21FEBF846CA8560950CAD611F98070E5E14F86FC68E429340").
		SetHeader("Sec-Ms-Gec-Version", "1-120.0.2210.133").
		Do()

	var resp ChatReq
	err := json.Unmarshal(c.GetBody(), &resp)
	if err != nil {
		return err
	}
	resp.ConversationSignature = c.GetHeader("X-Sydney-Conversationsignature")
	resp.EncryptedConversationSignature = c.GetHeader("X-Sydney-Encryptedconversationsignature")

	chat.chatHub = newChatHub(resp)

	return nil
}

func (chat *Chat) MsgComposer(msgs []Message) (prompt string, msg string) {
	systemMsgNum := 0
	for _, t := range msgs {
		if t.Role == "system" {
			systemMsgNum++
		}
	}
	if len(msgs)-systemMsgNum == 1 {
		return "", msgs[0].Content
	}

	var lastRole string
	for _, t := range msgs {
		if lastRole == t.Role {
			msg += "\n" + t.Content
			continue
		} else if lastRole != "" {
			msg += "\n\n"
		}
		switch t.Role {
		case "system":
			prompt += t.Content
		case "user":
			msg += "`me`:\n" + t.Content
		case "assistant":
			msg += "`you`:\n" + t.Content
		}
		if t.Role != "system" {
			lastRole = t.Role
		}
	}
	msg += "\n\n`you`:"
	return prompt, msg
}

func (chat *Chat) optionsSetsHandler(systemContext []SystemContext) []string {
	optionsSets := []string{
		"nlu_direct_response_filter",
		"deepleo",
		"disable_emoji_spoken_text",
		"responsible_ai_policy_235",
		"enablemm",
		"dv3sugg",
		"autosave",
		"iyxapbing",
		"iycapbing",
		"rai289",
		"enflst",
		"enpcktrk",
		"rcaldictans",
		"rcaltimeans",
		"eredirecturl",
	}
	if len(systemContext) > 0 {
		optionsSets = append(optionsSets, "nojbfedge")
	}

	tone := chat.GetStyle()
	if strings.Contains(tone, "g4t") {
		optionsSets = append(optionsSets, "dlgpt4t")
	}
	if strings.Contains(tone, PRECISE) {
		optionsSets = append(optionsSets, "h3precise", "clgalileo", "gencontentv3")
	} else if strings.Contains(tone, BALANCED) {
		optionsSets = append(optionsSets, "galileo", "saharagenconv5")
	} else if strings.Contains(tone, CREATIVE) {
		optionsSets = append(optionsSets, "h3imaginative", "clgalileo", "gencontentv3")
	}
	return optionsSets
}

func (chat *Chat) sliceIdsHandler(systemContext []SystemContext) []string {
	if len(systemContext) > 0 {
		return []string{
			"winmuid1tf",
			"styleoff",
			"ccadesk",
			"smsrpsuppv4cf",
			"ssrrcache",
			"contansperf",
			"crchatrev",
			"winstmsg2tf",
			"creatgoglt",
			"creatorv2t",
			"sydconfigoptt",
			"adssqovroff",
			"530pstho",
			"517opinion",
			"418dhlth",
			"512sprtic1s0",
			"emsgpr",
			"525ptrcps0",
			"529rweas0",
			"515oscfing2s0",
			"524vidansgs0",
		}
	} else {
		return []string{
			"techpillscf",
			"gbaa",
			"gba",
			"gbapa",
			"codecreator",
			"dlidcf",
			"specedge",
			"preall15",
			"suppsm240-t",
			"translref",
			"ardsw_1_9_9",
			"fluxnosearchc",
			"fluxnosearch",
			"1115rai289",
			"1119backoss0",
			"124multi2t",
			"1129gpt4ts0",
			"kchero50cf",
			"cacfastapis",
			"cacdupereccf",
			"cacmuidarb",
			"cacfrwebt2cf",
			"sswebtop2cf",
		}
	}
}

func (chat *Chat) pluginHandler(optionsSets *[]string) []Plugins {
	plugins := []Plugins{}
	tone := chat.GetStyle()
	if !strings.Contains(tone, "offline") {
		plugins = append(plugins, Plugins{Id: "c310c353-b9f0-4d76-ab0d-1dd5e979cf68"})
	} else {
		*optionsSets = append(*optionsSets, "nosearchall")
	}
	return plugins
}

func (chat *Chat) systemContextHandler(prompt string) []SystemContext {
	systemContext := []SystemContext{}
	if prompt != "" {
		systemContext = append(systemContext, SystemContext{
			Author:      "user",
			Description: prompt,
			ContextType: "WebPage",
			MessageType: "Context",
			MessageId:   "discover-web--page-ping-mriduna-----",
		})
	}
	return systemContext
}

func (chat *Chat) requestPayloadHandler(msg string, optionsSets []string, sliceIds []string, plugins []Plugins, systemContext []SystemContext) map[string]any {
	msgId := hex.NewUUID()
	tone := chat.GetStyle()

	data := map[string]any{
		"arguments": []any{
			map[string]any{
				"source":      "cib",
				"optionsSets": optionsSets,
				"allowedMessageTypes": []string{
					"ActionRequest",
					"Chat",
					"ConfirmationCard",
					"Context",
					"InternalSearchQuery",
					"InternalSearchResult",
					"Disengaged",
					"InternalLoaderMessage",
					"InvokeAction",
					"Progress",
					"RenderCardRequest",
					"RenderContentRequest",
					"AdsQuery",
					"SemanticSerp",
					"GenerateContentQuery",
					"SearchQuery",
				},
				"sliceIds":         sliceIds,
				"isStartOfSession": true,
				"verbosity":        "verbose",
				"scenario":         "SERP",
				"plugins":          plugins,
				"previousMessages": systemContext,
				"traceId":          strings.ReplaceAll(hex.NewUUID(), "-", ""),
				"conversationHistoryOptionsSets": []string{
					"autosave",
					"savemem",
					"uprofupd",
					"uprofgen",
				},
				"requestId": msgId,
				"message": map[string]any{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        msg,
					"messageType": "Chat",
					"requestId":   msgId,
					"messageId":   msgId,
				},
				// "conversationSignature": chat.GetChatHub().GetConversationSignature(),
				"tone":           strings.ReplaceAll(strings.ReplaceAll(tone, "-g4t", ""), "-offline", ""),
				"spokenTextMode": "None",
				"participant": map[string]any{
					"id": chat.GetChatHub().GetClientId(),
				},
				"conversationId": chat.GetChatHub().GetConversationId(),
			},
		},
		"invocationId": "0",
		"target":       "chat",
		"type":         4,
	}

	return data
}

func (chat *Chat) wsHandler(data map[string]any) (*websocket.Conn, error) {
	dialer := websocket.DefaultDialer
	dialer.Proxy = http.ProxyFromEnvironment
	headers := http.Header{}
	headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	headers.Add("User-Agent", userAgent)
	if chat.GetSydneyBaseUrl() == sydneyBaseUrl {
		headers.Add("Host", "sydney.bing.com")
		headers.Add("Origin", "https://www.bing.com")
	}

	ws, _, err := dialer.Dial(fmt.Sprintf(sydneyChatHubUrl, chat.SydneyBaseUrl, url.QueryEscape(chat.GetChatHub().GetEncryptedConversationSignature())), headers)
	if err != nil {
		return nil, err
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte("{\"protocol\":\"json\",\"version\":1}"+spilt))
	if err != nil {
		return nil, err
	}

	_, _, err = ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte("{\"type\":6}"+spilt))
	if err != nil {
		return nil, err
	}

	req, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = ws.WriteMessage(websocket.TextMessage, append(req, []byte(spilt)...))
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (chat *Chat) Chat(prompt, msg string) (string, error) {
	systemContext := chat.systemContextHandler(prompt)
	optionsSets := chat.optionsSetsHandler(systemContext)
	sliceIds := chat.sliceIdsHandler(systemContext)
	plugins := chat.pluginHandler(&optionsSets)
	data := chat.requestPayloadHandler(msg, optionsSets, sliceIds, plugins, systemContext)

	ws, err := chat.wsHandler(data)
	if err != nil {
		return "", err
	}
	defer ws.Close()

	text := ""

	i := 0
	for {
		if i >= 15 {
			err := ws.WriteMessage(websocket.TextMessage, []byte("{\"type\":6}"+spilt))
			if err != nil {
				break
			}
			i = 0
		}
		resp := new(ResponsePayload)
		err = ws.ReadJSON(&resp)
		if err != nil {
			if err.Error() != "EOF" {
				return text, err
			}
		}
		if resp.Type == 2 {
			if resp.Item.Result.Value == "CaptchaChallenge" {
				text = "User needs to solve CAPTCHA to continue."
			}
			break
		} else if resp.Type == 1 {
			if len(resp.Arguments) > 0 {
				if len(resp.Arguments[0].Messages) > 0 {
					text = resp.Arguments[0].Messages[0].Text
					// fmt.Println(resp.Arguments[0].Messages[0].Text + "\n\n")
				}
			}
		}
		i++
	}

	return text, nil
}

func (chat *Chat) ChatStream(prompt, msg string, c chan string) (string, error) {
	systemContext := chat.systemContextHandler(prompt)
	optionsSets := chat.optionsSetsHandler(systemContext)
	sliceIds := chat.sliceIdsHandler(systemContext)
	plugins := chat.pluginHandler(&optionsSets)
	data := chat.requestPayloadHandler(msg, optionsSets, sliceIds, plugins, systemContext)

	ws, err := chat.wsHandler(data)
	if err != nil {
		return "", err
	}
	defer ws.Close()

	text := ""

	i := 0
	for {
		if i >= 15 {
			err := ws.WriteMessage(websocket.TextMessage, []byte("{\"type\":6}"+spilt))
			if err != nil {
				break
			}
			i = 0
		}
		resp := new(ResponsePayload)
		err = ws.ReadJSON(&resp)
		if err != nil {
			if err.Error() != "EOF" {
				c <- "EOF"
				// tc <- "EOF"
				return "", err
			}
		}
		if resp.Type == 2 {
			if resp.Item.Result.Value == "CaptchaChallenge" {
				c <- "User needs to solve CAPTCHA to continue."
			}
			break
		} else if resp.Type == 1 {
			if len(resp.Arguments) > 0 {
				if len(resp.Arguments[0].Messages) > 0 {
					if resp.Arguments[0].Messages[0].MessageType == "InternalSearchResult" {
						continue
					}
					if resp.Arguments[0].Messages[0].MessageType == "InternalSearchQuery" || resp.Arguments[0].Messages[0].MessageType == "InternalLoaderMessage" {
						c <- resp.Arguments[0].Messages[0].Text
						c <- "\n\n"
						continue
					}
					if len(resp.Arguments[0].Messages[0].Text) > len(text) {
						c <- strings.ReplaceAll(resp.Arguments[0].Messages[0].Text, text, "")
						text = resp.Arguments[0].Messages[0].Text
					}
					// fmt.Println(resp.Arguments[0].Messages[0].Text + "\n\n")
				}
			}
		}
		i++
	}

	c <- "EOF"

	return text, nil
}
