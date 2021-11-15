package common

type MsgMeta struct {
	CliID   uint64 // who sends the message
	MsgBody string
}

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

type ServerDetail struct {
	Tematic    string
	TotalUsers uint64
	IP         string
}

type File struct {
	Filename string
	Content  []byte
	Creator  uint64
}
