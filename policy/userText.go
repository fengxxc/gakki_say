package policy

import (
	"log"
)

func UserText(userText string) Reply {
	var reply Reply = Reply{Type: Failed}
	switch userText {
	case "hello":
		reply.Type = Text
		reply.Body = []byte("Hello, I'm Gakki~")

	case "👍":
		// log.Println("gakki post Like Emoji~")

		img, imgType, err := getImg("./img/smile e03(3).png")
		if err != nil {
			log.Println(err)
			return reply
		}
		// TODO proccess img ...

		reply.Type = Image
		reply.Body = imgToBytes(img, imgType)
	case "ping":
		img, imgType, err := getImg("./img/pingpang.jpg")
		if err != nil {
			log.Println(err)
			return reply
		}
		reply.Type = Image
		reply.Body = imgToBytes(img, imgType)
	default:
		reply.Type = Text
		reply.Body = []byte("已婚，谢谢~")
	}
	return reply
}
