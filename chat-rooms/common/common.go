package common

type MsgMeta struct {
	CliID   uint64 // who sends the message
	MsgBody string
}

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

type File struct {
	Filename string
	Content  []byte
	Creator  uint64
}
