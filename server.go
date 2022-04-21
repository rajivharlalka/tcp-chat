package chat

type server struct {
	rooms    map[string]*room
	commands chan command
}
