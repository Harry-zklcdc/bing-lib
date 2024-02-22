package binglib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Harry-zklcdc/bing-lib/lib/aes"
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

	PRECISE_18K          = "Precise-18k"          // 精准
	BALANCED_18K         = "Balanced-18k"         // 平衡
	CREATIVE_18K         = "Creative-18k"         // 创造
	PRECISE_18K_OFFLINE  = "Precise-18k-offline"  // 精准, 不联网搜索
	BALANCED_18K_OFFLINE = "Balanced-18k-offline" // 平衡, 不联网搜索
	CREATIVE_18K_OFFLINE = "Creative-18k-offline" // 创造, 不联网搜索

	PRECISE_G4T_18K  = "Precise-g4t-18k"  // 精准 GPT4-Turbo
	BALANCED_G4T_18K = "Balanced-g4t-18k" // 平衡 GPT4-Turbo
	CREATIVE_G4T_18K = "Creative-g4t-18k" // 创造 GPT4-Turbo
)

var ChatModels = [21]string{BALANCED, BALANCED_OFFLINE, CREATIVE, CREATIVE_OFFLINE, PRECISE, PRECISE_OFFLINE, BALANCED_G4T, BALANCED_G4T_OFFLINE, CREATIVE_G4T, CREATIVE_G4T_OFFLINE, PRECISE_G4T, PRECISE_G4T_OFFLINE,
	BALANCED_18K, BALANCED_18K_OFFLINE, CREATIVE_18K, CREATIVE_18K_OFFLINE, PRECISE_18K, PRECISE_18K_OFFLINE, BALANCED_G4T_18K, CREATIVE_G4T_18K, PRECISE_G4T_18K}

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
		xff:           chat.xff,
		bypassServer:  chat.bypassServer,
		BingBaseUrl:   chat.BingBaseUrl,
		SydneyBaseUrl: chat.SydneyBaseUrl,
	}
}

func (chat *Chat) SetCookies(cookies string) *Chat {
	chat.cookies = cookies
	return chat
}

func (chat *Chat) SetXFF(xff string) *Chat {
	chat.xff = xff
	return chat
}

func (chat *Chat) SetBypassServer(bypassServer string) *Chat {
	chat.bypassServer = bypassServer
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

func (chat *Chat) GetXFF() string {
	return chat.xff
}

func (chat *Chat) GetBypassServer() string {
	return chat.bypassServer
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
	if chat.GetXFF() != "" {
		c.SetHeader("X-Forwarded-For", chat.xff)
	}
	if chat.GetBingBaseUrl() == bingBaseUrl {
		c.SetHeader("Host", "www.bing.com")
		c.SetHeader("Origin", "https://www.bing.com")
	}
	c.SetUrl(fmt.Sprintf(bingCreateConversationUrl, chat.BingBaseUrl)).
		SetUserAgent(userAgent).
		SetCookies(chat.cookies).
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US;q=0.9").
		SetHeader("Referer", "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx").
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
		SetHeader("X-Ms-Useragent", "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.12.3 OS/Windows").
		SetHeader("X-Ms-Client-Request-Id", hex.NewUUID()).
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
		"iyxapbing",
		"iycapbing",
		"enable_user_consent",
		"fluxmemcst",
		"gldcl1p",
		"uquopt",
		"langdtwb",
		"enflst",
		"enpcktrk",
		"rcaldictans",
		"rcaltimeans",
		"gndbfptlw",
	}
	if len(systemContext) > 0 {
		optionsSets = append(optionsSets, "nojbfedge", "rai278")
	}

	tone := chat.GetStyle()
	if strings.Contains(tone, "g4t") {
		optionsSets = append(optionsSets, "dlgpt4t")
	}
	if strings.Contains(tone, "18k") {
		optionsSets = append(optionsSets, "prjupy")
	}
	if strings.Contains(tone, PRECISE) {
		optionsSets = append(optionsSets, "h3precise", "clgalileo", "gencontentv3")
	} else if strings.Contains(tone, BALANCED) {
		if strings.Contains(tone, "18k") {
			optionsSets = append(optionsSets, "clgalileo", "saharagenconv5")
		} else {
			optionsSets = append(optionsSets, "galileo", "saharagenconv5")
		}
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
			"qnacnc",
			"fluxsunoall",
			"mobfdbk",
			"v6voice",
			"cmcallcf",
			"specedge",
			"tts5",
			"advperfcon",
			"designer2cf",
			"defred",
			"msgchkcf",
			"thrdnav",
			"0212boptpsc",
			"116langwb",
			"124multi2t",
			"927storev2s0",
			"0131dv1",
			"1pgptwdes",
			"0131gndbfpr",
			"brndupdtcf",
			"enter4nl",
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

func (chat *Chat) requestPayloadHandler(msg string, optionsSets []string, sliceIds []string, plugins []Plugins, systemContext []SystemContext) (data map[string]any, msgId string) {
	msgId = hex.NewUUID()
	tone := chat.GetStyle()

	data = map[string]any{
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
				"gptId":            "copilot",
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
					"author":        "user",
					"inputMethod":   "Keyboard",
					"text":          msg,
					"messageType":   "Chat",
					"requestId":     msgId,
					"messageId":     msgId,
					"userIpAddress": chat.GetXFF(),
					"locale":        "zh-CN",
					"market":        "en-US",
					"region":        "US",
					"location":      "lat:47.639557;long:-122.128159;re=1000m;",
					"locationHints": []any{
						map[string]any{
							"country":           "United States",
							"state":             "California",
							"city":              "Los Angeles",
							"timezoneoffset":    8,
							"countryConfidence": 8,
							"Center": map[string]any{
								"Latitude":  78.4156,
								"Longitude": -101.4458,
							},
							"RegionType": 2,
							"SourceType": 1,
						},
					},
				},
				// "conversationSignature": chat.GetChatHub().GetConversationSignature(),
				"tone":           strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(tone, "-18k", ""), "-g4t", ""), "-offline", ""),
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

	return
}

func (chat *Chat) wsHandler(data map[string]any) (*websocket.Conn, error) {
	dialer := websocket.DefaultDialer
	dialer.Proxy = http.ProxyFromEnvironment
	headers := http.Header{}
	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", userAgent)
	headers.Set("Referer", "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx")
	headers.Set("Cookie", chat.cookies)
	if chat.GetXFF() != "" {
		headers.Set("X-Forwarded-For", chat.xff)
	}
	if chat.GetSydneyBaseUrl() == sydneyBaseUrl {
		headers.Set("Host", "sydney.bing.com")
		headers.Set("Origin", "https://www.bing.com")
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
	data, msgId := chat.requestPayloadHandler(msg, optionsSets, sliceIds, plugins, systemContext)

	ws, err := chat.wsHandler(data)
	if err != nil {
		return "", err
	}
	defer ws.Close()

	text := ""
	verifyStatus := false

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
			} else {
				return text, nil
			}
		}
		if resp.Type == 2 {
			if resp.Item.Result.Value == "CaptchaChallenge" || resp.Item.Result.Value == "Throttled" {
				if resp.Item.Result.Value == "CaptchaChallenge" {
					text = "User needs to solve CAPTCHA to continue."
				} else if resp.Item.Result.Value == "Throttled" {
					text = "Request is throttled."
				} else {
					text = "Unknown error."
				}
				if chat.GetBypassServer() != "" && !verifyStatus {
					IG := hex.NewUpperHex(32)
					T, err := aes.Encrypt("Harry-zklcdc/go-proxy-bingai", IG)
					if err != nil {
						break
					}
					r, status, err := Bypass(chat.GetBypassServer(), chat.GetCookies(), "local-gen-"+hex.NewUUID(), IG, chat.GetChatHub().GetConversationId(), msgId, T)
					if err != nil || status != http.StatusOK {
						break
					}
					verifyStatus = true
					chat.SetCookies(r.Result.Cookies)
					ws.Close()
					data["invocationId"] = "1"
					ws, err = chat.wsHandler(data)
					if err != nil {
						break
					}
					defer ws.Close()
				} else {
					break
				}
			} else if resp.Item.Result.Value == "Success" {
				if len(resp.Item.Messages) > 1 {
					for i, v := range resp.Item.Messages[len(resp.Item.Messages)-1].SourceAttributions {
						text += "\n[^" + strconv.Itoa(i+1) + "^]: [" + v.ProviderDisplayName + "](" + v.SeeMoreUrl + ")"
					}
				}
				break
			} else {
				break
			}

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
	data, msgId := chat.requestPayloadHandler(msg, optionsSets, sliceIds, plugins, systemContext)

	ws, err := chat.wsHandler(data)
	if err != nil {
		c <- "EOF"
		return "", err
	}
	defer ws.Close()

	text := ""
	verifyStatus := false

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
				return "", err
			} else {
				c <- "EOF"
				return "", nil
			}
		}
		if resp.Type == 2 {
			if resp.Item.Result.Value == "CaptchaChallenge" || resp.Item.Result.Value == "Throttled" {
				if chat.GetBypassServer() != "" && !verifyStatus {
					c <- "Bypassing... Please Wait.\n\n"
					IG := hex.NewUpperHex(32)
					T, err := aes.Encrypt("Harry-zklcdc/go-proxy-bingai", IG)
					if err != nil {
						c <- "Bypass Fail!"
						break
					}
					r, status, err := Bypass(chat.GetBypassServer(), chat.GetCookies(), "local-gen-"+hex.NewUUID(), IG, chat.GetChatHub().GetConversationId(), msgId, T)
					if err != nil || status != http.StatusOK {
						c <- "Bypass Fail!"
						break
					}
					verifyStatus = true
					chat.SetCookies(r.Result.Cookies)
					ws.Close()
					data["invocationId"] = "1"
					ws, err = chat.wsHandler(data)
					if err != nil {
						c <- "Bypass Fail!"
						break
					}
					defer ws.Close()
				} else {
					if resp.Item.Result.Value == "CaptchaChallenge" {
						c <- "User needs to solve CAPTCHA to continue."
					} else if resp.Item.Result.Value == "Throttled" {
						c <- "Request is throttled."
					} else {
						c <- "Unknown error."
					}
					break
				}
			} else if resp.Item.Result.Value == "Success" {
				if len(resp.Item.Messages) > 1 {
					for i, v := range resp.Item.Messages[len(resp.Item.Messages)-1].SourceAttributions {
						c <- "\n[^" + strconv.Itoa(i+1) + "^]: [" + v.ProviderDisplayName + "](" + v.SeeMoreUrl + ")"
					}
				}
				break
			} else {
				break
			}
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
