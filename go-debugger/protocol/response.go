package protocol

type Response struct {
	Sequence uint        `json:"sequence"`
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}
