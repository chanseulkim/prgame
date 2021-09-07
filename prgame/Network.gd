extends Node


# Declare member variables here. Examples:
# var a = 2
# var b = "text"

var players = {}
var udp_sock = PacketPeerUDP.new()

func _ready():
	pass
	
remote func add_player(uid, player):
	players[uid] = player
	
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
