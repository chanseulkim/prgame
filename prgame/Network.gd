extends Node

var players = {}
var udp_sock = PacketPeerUDP.new()
var main_uid
func _ready():
	pass
	
remote func add_player(uid, player):
	players[uid] = player
remote func get_player(uid):
	if uid in players:
		return players[uid]
	return null
	
func get_main_player():
	return players[main_uid]
	
func connect2Server(serv_ip, serv_port, id, position):
	udp_sock.set_dest_address(serv_ip, serv_port)
	var pac = id + ";" + "enter;" + position + ";" + "m;"
	send(pac)
	print("player " + id + " connected")
	
func send(packet):
	udp_sock.put_packet(packet.to_ascii())
	
func read():
	if udp_sock.get_available_packet_count() > 0:
		return udp_sock.get_packet().get_string_from_ascii()
