extends Node2D

func _ready():
	print("world ready")
	var local_player_id = get_tree().get_network_unique_id()
	if not(get_tree().is_network_server()):
		rpc_id(1, '_request_player_info', local_player_id)
	
	var new_player = preload('res://Player.tscn').instance()
	new_player.name = str(get_tree().get_network_unique_id())
	new_player.set_network_master(get_tree().get_network_unique_id())
	new_player.connect2Server("127.0.0.1", 50080)
	new_player.init(100.0, 100.0)
	add_child(new_player)
	pass
