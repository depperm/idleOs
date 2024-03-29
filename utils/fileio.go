package utils

import (
	"encoding/binary"
	"io"
	"os"
	"sync"
	"time"

	pb "github.com/depperm/idleOs/proto"
	"google.golang.org/protobuf/proto"
)

// Write writes a protobuf message to an io.Writer (e.g., a file).
func Write(w io.Writer, msg []byte) error {
	// Write the length of the message as a 4-byte little-endian integer
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(msg)))
	if _, err := w.Write(buf); err != nil {
		return err
	}

	// Write the actual message bytes
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return nil
}

func SaveProto(myProto *pb.PlayerInfo) error {
	msgBytes, err := proto.Marshal(myProto)
	if err != nil {
		return err
	}

	file, err := os.Create("data.isf")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the protobuf message to the file
	if err := Write(file, msgBytes); err != nil {
		return err
	}
	return nil
}

func AutoSave(playerInfo *pb.PlayerInfo) {
	var mu sync.Mutex
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// fmt.Printf("should save: %+v\n", gameState)
		mu.Lock()
		playerCopy := playerInfo
		mu.Unlock()
		playerCopy.LastSave = time.Now().UnixMilli()

		err := SaveProto(playerCopy)
		if err != nil {
			os.Exit(2)
		}
	}
}

func LoadGame() (*GameState, error) {
	// Open the encoded file for reading
	file, err := os.Open("data.isf")
	if err != nil {
		// File does not exist, new game
		var playerInfo pb.PlayerInfo
		playerInfo.Username = "root"
		playerInfo.Lines = 0
		// gameState.CurrentDir = "/"
		playerInfo.Dirs = &pb.Directory{
			Name: "/",
			Files: []*pb.File{
				{Name: "code.exe", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Size: 5, Contents: ""},
			},
			Permissions: "rw-rw-rw-",
			Owner:       "root",
			ModifyDate:  time.Now().UnixMilli(),
			Dirs: []*pb.Directory{
				{Name: "team", Owner: "root", Permissions: "rw-rw-rw-", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "languages", Owner: "root", Permissions: "rw-rw-rw-", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "tech", Owner: "root", Permissions: "rw-rw-rw-", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "shop", Owner: "root", Permissions: "rw-rw-rw-", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "usr", Owner: "root", Permissions: "rw-rw-rw-", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: ".achievement", Owner: "root", Permissions: "rw-rw-rw-", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
			},
		}
		gameState := GameState{}
		gameState.CurrentDir = "/"
		gameState.Player = playerInfo
		//gameState.achievements
		return &gameState, nil
	}
	defer file.Close()

	// Read the length of the string (as uint32) from the file
	var length uint32
	err = binary.Read(file, binary.LittleEndian, &length)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Read the encoded string data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Decode the string data
	str := string(data)

	// Print the decoded string
	// fmt.Println("Decoded string:", str)

	// Convert the decoded string to a Protocol Buffers message
	gameState := &GameState{}
	err = proto.Unmarshal([]byte(str), &gameState.Player)
	if err != nil {
		// Handle error
		return nil, err
	}
	// fmt.Printf("loaded: %+v\n", gameState)
	// time.Sleep(5 * time.Second)
	return gameState, nil
}
