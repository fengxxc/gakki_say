package policy

import (
	"log"

	"strings"

	mapset "github.com/deckarep/golang-set"
)

func UserText(userText string, symbolMaps *SymbolMaps) Reply {
	var reply Reply = Reply{Type: Failed, Body: []byte(userText)}
	userText = strings.Trim(userText, " ")
	// strings.Split(userText, " ")
	var idx int = strings.Index(userText, " ")
	var firstPiece string
	var secondPiece string = ""
	if idx == -1 {
		firstPiece = userText
	} else {
		firstPiece = userText[0:idx]
		secondPiece = userText[idx+1:]
	}

	// 固定的回复
	if ok := staticReply(firstPiece, &reply); ok {
		return reply
	}

	// gakki图片 + 用户文本 的回复
	// emoji
	if symbolMaps.ContainsEmoji(firstPiece) {
		var imgNameSet mapset.Set = symbolMaps.EmojiMap[firstPiece]
		imgName := getSetRandom(imgNameSet).(string)
		img, err := imgWriteText("./img/"+imgName, secondPiece, 0.5, 0.5, &RGBA{89, 89, 89, 64})
		if err != nil {
			log.Println(err)
			return reply
		}
		reply.Type = Image
		reply.Body = imgToBytes(img, getImgTypeByFileName(imgName))
	}
	return reply
}

func staticReply(text string, reply *Reply) bool {
	var ok bool = false
	switch text {
	case "hello":
		reply.Type = Text
		reply.Body = []byte("Hello, I'm Gakki~")
		ok = true
	case "ping":
		fileName := "./img/pingpang.jpg"
		img, err := imgWriteText(fileName, "pang~", 0.5, 0.5, &RGBA{255, 204, 255, 89})
		if err != nil {
			log.Println(err)
			return false
		}
		reply.Type = Image
		reply.Body = imgToBytes(img, getImgTypeByFileName(fileName))
		ok = true
	}
	return ok
}
