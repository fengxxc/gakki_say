package policy

func Command(cmd string) string {
	var reply string
	switch cmd {
	case "start":
		reply = "初次见面，请多指教，我是图文并茂的Gakki~"
	case "help":
		reply = "我还没想好怎么帮你"
	case "settings":
		reply = "这个功能还没做好……再等等"
	}
	return reply
}
