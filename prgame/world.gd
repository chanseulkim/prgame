extends Node2D

func _ready():
	print("display size : ", OS.get_screen_size())
	print("window size : ", OS.get_real_window_size())
	print("world ready")
	var local_player_id = get_tree().get_network_unique_id()
	if not(get_tree().is_network_server()):
		rpc_id(1, '_request_player_info', local_player_id)

	var timestamp = OS.get_ticks_msec()
	var my_id = "p" + String(timestamp)
	
	Network.connect2Server("127.0.0.1", 51081, my_id, "(100.0, 100.0)")
	
	var player = load('res://Player.tscn').instance()
	player.name = my_id
	player.set_network_master(get_tree().get_network_unique_id())
	player.init(100.0, 100.0, false)
	add_child(player)
	Network.add_player(my_id, player)
	Network.main_uid = my_id
	#get_tree().change_scene('res://Game.tscn')
	
#	var noti_msg = "obj;noti;"
#	var cd = get_children()
#	for c in cd:
#		if "field" in c.name:
#			noti_msg += c.name + ";" + str(c.position) + ";m;"
#			var n = c.get_child()
#			Network.send(noti_msg)
#			c.hide()
#	pass


