package utils

import (
	pb "github.com/depperm/idleOs/proto"
)

type GameState struct {
	CurrentDir   string
	Achievements string
	Player       pb.PlayerInfo
}
