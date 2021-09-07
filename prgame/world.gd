extends Node2D

func _ready():
	print("display size : ", OS.get_screen_size())
	print("window size : ", OS.get_real_window_size())
	print("world ready")
	var local_player_id = get_tree().get_network_unique_id()
	if not(get_tree().is_network_server()):
		rpc_id(1, '_request_player_info', local_player_id)
	
	var new_player = load('res://Player.tscn').instance()
	var uid = str(get_tree().get_network_unique_id())
	new_player.name = uid
	new_player.set_network_master(get_tree().get_network_unique_id())
	new_player.init(100.0, 100.0)
	add_child(new_player)
	Network.connect2Server("127.0.0.1", 50080, new_player.name, "(100.0, 100.0)")
	pass
