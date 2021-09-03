extends Area2D

export var speed = 400  # How fast the player will move (pixels/sec).
var screen_size  # Size of the game window.
var udp_sock = PacketPeerUDP.new()
var default_pos_x = float(100.0)
var default_pos_y = float(100.0)
var timestamp = OS.get_ticks_msec()
var my_id = "p" + String(timestamp)
var is_connected_to_server = false

# user's id and child index
var user_index_map = {}

# * Pakcet format
# * user_id;command;action;delta-time;

var myindex
var this
var keepalive_thread

func _parse_msg(msg):
	if len(msg) <= 0:
		return
	var ret = []
	var header = []
	while true:
		var f1 = msg.find(';')
		if f1 == -1:
			return ret
		var data = msg.substr(0, f1)
		msg = msg.substr(f1+1)
		header.push_back(data)
		if data == "m":
			ret.push_back(header)
			break
	
	return ret

func _process_msg(msg):
	var headers = _parse_msg(msg)
	for header in headers:
		if len(header) < 2:
			continue
		print(header)
		var user_id = String(header[0])
		var command = String(header[1])
		if (command == "sync"):
			if user_id == my_id:
				continue
			if len(header) < 3:
				continue
			var x = default_pos_x
			var y = default_pos_y
			var pos_str = String(header[2])
			if len(pos_str) != 0:
				var trimedpos = pos_str.substr(1, len(pos_str))
				var p = trimedpos.find(",")
				if p != -1:
					x = trimedpos.substr(0, p).to_float()
					y = trimedpos.substr(p+1).to_float()
			var new_player = load('res://Player.tscn').instance()
			new_player.name = user_id
			new_player.init(x, y)
			add_child(new_player)
			user_index_map[user_id] = get_child_count() - 1
			pass
		if (command == "join") && (user_id != my_id):
			var new_player = load('res://Player.tscn').instance()
			new_player.name = user_id
			new_player.init(default_pos_x, default_pos_y)
			add_child(new_player)
			user_index_map[user_id] = get_child_count() - 1
		if len(header) < 4:
			return
		var action = String(header[2])
		var delta_time = String(header[3])
			
		if (command == "move"):
			_handle_move(user_id, action, delta_time)
		pass
		
	

func _handle_move(user_id, action, delta_time):
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
		$AnimatedSprite.play()
	else:		
		$AnimatedSprite.stop()

	var delta = delta_time.to_float()
	if user_id != my_id:
		var child_index = user_index_map[user_id]
		var child = get_child(child_index)
		child.position += velocity * delta
		child.position.x = clamp(child.position.x, 0, screen_size.x)
		child.position.y = clamp(child.position.y, 0, screen_size.y)
		return
	this.position += velocity * delta
	this.position.x = clamp(this.position.x, 0, screen_size.x)
	this.position.y = clamp(this.position.y, 0, screen_size.y)
		
	if velocity.x != 0:
		$AnimatedSprite.animation = "walk"
		$AnimatedSprite.flip_v = false
		# See the note below about boolean assignment
		$AnimatedSprite.flip_h = velocity.x < 0
	elif velocity.y != 0:
		$AnimatedSprite.animation = "up"
		$AnimatedSprite.flip_v = velocity.y > 0
func init(x, y):
	self.position += Vector2(x, y)
	pass
	
func connect2Server(serv_ip, serv_port):
	udp_sock.set_dest_address(serv_ip, serv_port)
	var pac = my_id + ";" + "join;m;"
	udp_sock.put_packet(pac.to_ascii())
	print("player " + my_id + " ready")
	
# Called when the node enters the scene tree for the first time.
func _ready():
	name = my_id
	screen_size = get_viewport_rect().size
	myindex = get_index()
	this = get_child(myindex)
	keepalive_thread = Thread.new()
	keepalive_thread.start(self, "_keepalive_loop")
	
func _keepalive_loop():
	while true:
		print("ping")
		_wait(10)
		if udp_sock != null:
			var pac = my_id + ";ping;m;"
			udp_sock.put_packet(pac.to_ascii())
		pass
	keepalive_thread.wait_to_finish()
	pass
	
	

const RED = Color(1.0, 0, 0, 0.4)
const GREEN = Color(0, 1.0, 0, 0.4)
var color = GREEN
var radius = 80
var center = Vector2(default_pos_x, default_pos_y)
var rotation_angle = 50
var angle_from = 75
var angle_to = 195

var last_mouse_pos = Vector2()
func _input(event):
	if event is InputEventMouseButton:
		print("Mouse Click/Unclick at: ", event.position)
	elif event is InputEventMouseMotion:
		print("Mouse Motion at: ", event.position)
		# 캐릭터 위치기준 마우스의 위치 변화
		if (last_mouse_pos.x - event.position.x) < 0 :
			angle_from += 3
			angle_to += 3
		if (last_mouse_pos.y - event.position.y) > 0 :
			angle_from -= 3
			angle_to -= 3
		if (last_mouse_pos.y - event.position.y) < 0 :
			angle_from += 3
			angle_to += 3
		last_mouse_pos= event.position
		
		
#	angle_from += rotation_angle
#	angle_to += rotation_angle
	# We only wrap angles when both of them are bigger than 360.
	if angle_from > 360 and angle_to > 360:
		angle_from = wrapf(angle_from, 0, 360)
		angle_to = wrapf(angle_to, 0, 360)
	update()
#	print("Viewport Resolution is: ", get_viewport_rect().size)

func _draw():
	draw_circle_arc_poly(Vector2(position.x, position.y), radius, angle_from, angle_to, color)
	#draw_circle_arc(Vector2(position.x, position.y), radius, angle_from, angle_to, color)
	
func draw_circle_arc_poly(center, radius, angle_from, angle_to, color):
	var nb_points = 32
	var points_arc = PoolVector2Array()
	points_arc.push_back(center)
	var colors = PoolColorArray([color])

	for i in range(nb_points + 1):
		var angle_point = deg2rad(angle_from + i * (angle_to - angle_from) / nb_points - 90)
		points_arc.push_back(center + Vector2(cos(angle_point), sin(angle_point)) * radius)
	draw_polygon(points_arc, colors)

func draw_circle_arc(center, radius, angle_from, angle_to, color):
	var nb_points = 32
	var points_arc = PoolVector2Array()
	for i in range(nb_points + 1):
		var angle_point = deg2rad(angle_from + i * (angle_to-angle_from) / nb_points - 90)
		points_arc.push_back(center + Vector2(cos(angle_point), sin(angle_point)) * radius)

	for index_point in range(nb_points):
		draw_line(points_arc[index_point], points_arc[index_point + 1], color)
		
# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	
	if this == null:
		return
	if udp_sock.get_available_packet_count() > 0:
		var recved_msg = udp_sock.get_packet().get_string_from_ascii()
		_process_msg(recved_msg)
	
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
	if Input.is_key_pressed(KEY_K):
		msg_tosend = "key;pressed_K;"
		
	if velocity.length() > 0:
		velocity = velocity.normalized() * speed
		$AnimatedSprite.play()
	else:
		$AnimatedSprite.stop()
		
	if msg_tosend.length() > 0 :
		msg_tosend += String(delta) + ";" + String(speed) + ";" + str(this.position) + ";" 
		var pac = my_id + ";" + msg_tosend + "m;"
		udp_sock.put_packet(pac.to_ascii())
	pass


signal timer_end

func _wait( seconds ):
	self._create_timer(self, seconds, true, "_emit_timer_end_signal")
	yield(self,"timer_end")

func _emit_timer_end_signal():
	emit_signal("timer_end")

func _create_timer(object_target, float_wait_time, bool_is_oneshot, string_function):
	var timer = Timer.new()
	timer.set_one_shot(bool_is_oneshot)
	timer.set_timer_process_mode(0)
	timer.set_wait_time(float_wait_time)
	timer.connect("timeout", object_target, string_function)
	self.add_child(timer)
	timer.start()
