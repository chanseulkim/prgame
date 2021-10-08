package gcore

import "strings"

const MAX_MSG_HEADER_COUNT = 20

// * Pakcet format
// * requirements : user_id;command;
// * user_id;command;action;delta-time;
const (
	CMD_ENTER  byte = 0
	CMD_LEAVE  byte = 1
	CMD_ACTION byte = 2
	CMD_MOVE   byte = 3 //
	CMD_PING   byte = 4
)
const (
	ACT_UIUP    byte = 0
	ACT_UIDOWN  byte = 1
	ACT_UILEFT  byte = 2
	ACT_UIRIGHT byte = 3
)

type MsgHeader struct {
	Command   byte
	Action    byte
	Timestamp byte
}

func SpliteMsg(msg string) [MAX_MSG_HEADER_COUNT]string {
	var ret [MAX_MSG_HEADER_COUNT]string
	if len(msg) <= 0 {
		return ret
	}
	for i := 0; i < MAX_MSG_HEADER_COUNT; i++ {
		token_pos := strings.Index(msg, ";")
		if token_pos == -1 {
			return ret
		}
		data := msg[0:token_pos]
		msg = msg[token_pos+1:]
		ret[i] = data
	}
	return ret
}
