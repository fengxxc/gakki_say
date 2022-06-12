package policy

import "log"

func Command(cmd string) Reply {
	var reply Reply = Reply{Type: Failed, Body: []byte("")}
	switch cmd {
	case "start":
		reply.Type = Text
		reply.Body = []byte("初次见面，请多指教，我是图文并茂的Gakki~")
	case "help":
		reply.Type = Text
		reply.Body = []byte("我还没想好怎么帮你")
	case "settings":
		reply.Type = Text
		reply.Body = []byte("这个功能还没做好……再等等")
	case "ping":
		fileName := "./img/pingpang.jpg"
		img, err := imgWriteText(fileName, "pang~", DrawStringConfig{
			ax:          0.5,
			ay:          0.5,
			fontFamily:  "SIMYOU.TTF",
			textBgColor: &RGBA{255, 204, 255, 89},
		})
		if err != nil {
			log.Println(err)
			return reply
		}
		reply.Type = Image
		reply.Body = imgToBytes(img, getImgTypeByFileName(fileName))
	}

	return reply
}
