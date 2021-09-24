extends Area2D

var world_node = preload("res://world.tscn")

export var speed = 400  # How fast the player will move (pixels/sec).
var screen_size # Size of the game window.
var default_pos_x = float(100.0)
var default_pos_y = float(100.0)
var timestamp = OS.get_ticks_msec()
var is_connected_to_server = false

# * Pakcet format
# * user_id;command;action;delta-time;

var is_peer = false

func _parse_msg(msg):
	if len(msg) <= 0:
		return
	var ret = []
	var header = []
	while true:
		var f1 = msg.find(';')
		if f1 == -1:
			break
		var data = msg.substr(0, f1)
		msg = msg.substr(f1+1)
		header.push_back(data)
		if data == "m":
			ret.push_back(header)
	return ret

var having_objs = []
var enemy

func _process_msg(msg):
	var headers = _parse_msg(msg)
	if headers == null:
		return
	for header in headers:
		if len(header) < 2:
			continue
		var user_id = String(header[0])
		var command = String(header[1])
		if (command == "objects"):
			var obj_name = header[2]
			if obj_name == "enemy":
				if enemy != null:
					enemy.show()
					continue
				enemy = load('res://Enemy.tscn').instance()
				enemy.name = "enemy"
				var pos_str = header[3]
				var pos = posstr2pos(pos_str)
				enemy.init(pos.x, pos.y)
				print("enemy: "+ header[2] + "," + pos_str)
				$'/root/world'.add_child(enemy)
				having_objs[obj_name] = enemy
			elif enemy != null:
				enemy.hide()
		if (command == "sync"):
			if user_id == name:
				continue
			if len(header) < 3:
				continue
			var x = default_pos_x
			var y = default_pos_y
			var pos_str = String(header[2])
			var pos = posstr2pos(pos_str)
			if len(pos_str) != 0:
				x = pos.x
				y = pos.y
			var new_player = load('res://Player.tscn').instance()
			new_player.name = user_id
			new_player.init(x, y, true)
			$'/root/world'.add_child(new_player)
			Network.add_player(user_id, new_player)
		if (command == "enter") && (user_id != name):
			var new_player = load('res://Player.tscn').instance()
			new_player.name = user_id
			new_player.init(default_pos_x, default_pos_y, true)
			$'/root/world'.add_child(new_player)
			Network.add_player(user_id, new_player)
		
		if len(header) < 4:
			return
		var action = String(header[2])
		var delta_time = String(header[3])
			
		if (command == "move"):
			_handle_move_msg(user_id, action, delta_time)

func _handle_move_msg(user_id, action, delta_time):
	if user_id == name : 
		return
	var velocity = Vector2()  # The player's movement vector.
	if action == "ui_right":
		velocity.x += 1
	if action == "ui_left":
		velocity.x -= 1
	if action == "ui_down":
		velocity.y += 1
	if action == "ui_up":
		velocity.y -= 1
	if velocity.length() > 0:
		velocity = velocity.normalized() * speed
	var player = Network.get_player(user_id)
	var delta = delta_time.to_float()
	player.position += velocity * delta
	player.position.x = clamp(player.position.x, 0, screen_size.x)
	player.position.y = clamp(player.position.y, 0, screen_size.y)
	return

func init(x, y, _is_peer):
	self.position += Vector2(x, y)
	is_peer = _is_peer
	pass

func _ready():
	screen_size = get_viewport_rect().size
	pass

func _process(delta):
	var recved_msg = Network.read()
	if recved_msg != null :
		_process_msg(recved_msg)
	if name != Network.main_uid :
		return
	var msg_tosend = ""
	var velocity = Vector2()  # The player's movement vector.
	if Input.is_action_pressed("ui_right"):
		msg_tosend = "move;ui_right;"
		velocity.x += 1
	if Input.is_action_pressed("ui_left"):
		msg_tosend = "move;ui_left;"
		velocity.x -= 1
	if Input.is_action_pressed("ui_down"):
		msg_tosend = "move;ui_down;"
		velocity.y += 1
	if Input.is_action_pressed("ui_up"):
		msg_tosend = "move;ui_up;"
		velocity.y -= 1
		
	if velocity.length() > 0:
		velocity = velocity.normalized() * speed
		$AnimatedSprite.play()
		position += velocity * delta
#		print(position)
		position.x = clamp(position.x, 0, screen_size.x)
		position.y = clamp(position.y, 0, screen_size.y)
	else:
		$AnimatedSprite.stop()
		
	if velocity.x != 0:
		$AnimatedSprite.animation = "walk"
		$AnimatedSprite.flip_v = false
		# See the note below about boolean assignment
		$AnimatedSprite.flip_h = velocity.x < 0
	elif velocity.y != 0:
		$AnimatedSprite.animation = "up"
		$AnimatedSprite.flip_v = velocity.y > 0
		
	if msg_tosend.length() > 0 :
		msg_tosend += String(delta) + ";" + String(speed) + ";" + str(self.position) + ";" 
		var pac = name + ";" + msg_tosend + "m;"
		Network.send(pac)
	pass

func posstr2pos(pos_str):
	var v = Vector2(0, 0)
	if len(pos_str) != 0:
		var trimedpos = pos_str.substr(1, len(pos_str))
		var p = trimedpos.find(",")
		if p != -1:
			v.x = trimedpos.substr(0, p).to_float()
			v.y = trimedpos.substr(p+1).to_float()
	return v
