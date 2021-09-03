package packetizer

import "strings"

const MAX_MSG_HEADER_COUNT = 20

// * Pakcet format
// * requirements : user_id;command;
// * user_id;command;action;delta-time;
const (
	CMD_Enter byte = 0
	CMD_Leave      = 1
	CMD_Action
	CMD_Move
	CMD_Ping
)

type msg_header struct {
	Command   byte
	Action    byte
	Timestamp byte
}

func ParseMsg(msg string) [MAX_MSG_HEADER_COUNT]string {
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
