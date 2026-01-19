package pkg

type ChatModel uint

const (
	DeepSeekChat     ChatModel = iota
	DeepSeekReasoner ChatModel = iota
)

func (c ChatModel) String() string {
	switch c {
	case DeepSeekChat:
		return "deepseek-chat"
	case DeepSeekReasoner:
		return "deepseek-reasoner"
	default:
		return ""
	}
}

type FileType uint

const (
	FileTypeDoc  FileType = iota
	FileTypeCode FileType = iota
)

func (c FileType) String() string {
	switch c {
	case FileTypeDoc:
		return "手册"
	case FileTypeCode:
		return "代码"
	default:
		return ""
	}
}
