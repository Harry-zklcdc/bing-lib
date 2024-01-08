package main

import (
	"fmt"
	"os"

	binglib "github.com/Harry-zklcdc/bing-lib"
)

var cookie = os.Getenv("COOKIE")

/*
流式输出
*/
func main() {
	c := binglib.NewChat(cookie)
	c.NewConversation()

	text := make(chan string)
	var tmp string
	go c.ChatStream("", "你好", text)

	for {
		tmp = <-text
		if tmp == "EOF" {
			break
		}
		fmt.Print(tmp)
	}
}
