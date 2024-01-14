package binglib

func newChatHub(chatReq ChatReq) *ChatHub {
	return &ChatHub{
		chatReq: chatReq,
		style:   CREATIVE,
	}
}

func (chatHub *ChatHub) Clone() *ChatHub {
	return &ChatHub{
		chatReq: ChatReq{
			ConversationId:                 chatHub.chatReq.ConversationId,
			ClientId:                       chatHub.chatReq.ClientId,
			ConversationSignature:          chatHub.chatReq.ConversationSignature,
			EncryptedConversationSignature: chatHub.chatReq.EncryptedConversationSignature,
		},
		style: chatHub.style,
	}
}

func (chatHub *ChatHub) SetStyle(style string) *ChatHub {
	chatHub.style = style
	return chatHub
}

func (chatHub *ChatHub) SetChatReq(chatReq ChatReq) *ChatHub {
	chatHub.chatReq = chatReq
	return chatHub
}

func (chatHub *ChatHub) SetConversationId(conversationId string) *ChatHub {
	chatHub.chatReq.ConversationId = conversationId
	return chatHub
}

func (chatHub *ChatHub) SetClientId(clientId string) *ChatHub {
	chatHub.chatReq.ClientId = clientId
	return chatHub
}

func (chatHub *ChatHub) SetConversationSignature(conversationSignature string) *ChatHub {
	chatHub.chatReq.ConversationSignature = conversationSignature
	return chatHub
}

func (chatHub *ChatHub) SetEncryptedConversationSignature(encryptedconversationsignature string) *ChatHub {
	chatHub.chatReq.EncryptedConversationSignature = encryptedconversationsignature
	return chatHub
}

func (chatHub *ChatHub) GetStyle() string {
	return chatHub.style
}

func (chatHub *ChatHub) GetChatReq() ChatReq {
	return chatHub.chatReq
}

func (chatHub *ChatHub) GetConversationId() string {
	return chatHub.chatReq.ConversationId
}

func (chatHub *ChatHub) GetClientId() string {
	return chatHub.chatReq.ClientId
}

func (chatHub *ChatHub) GetConversationSignature() string {
	return chatHub.chatReq.ConversationSignature
}

func (chatHub *ChatHub) GetEncryptedConversationSignature() string {
	return chatHub.chatReq.EncryptedConversationSignature
}
