package core_libs

var last_cmd map[string]map[string]string = make(map[string]map[string]string) // usrid, [cmd, cmd value]

// func GetCommand(usrid string, cmd string, value string) string{
// 	cmdmap := last_cmd[usrid]
// 	val := cmdmap[cmd]
// }
