package idleOs

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	//  . "filesamedir"

	pb "github.com/depperm/idleOs/proto"
	"github.com/inancgumus/screen"
	"google.golang.org/protobuf/proto"
)

type Man struct {
	Name        string     `json:"name"`
	Synopsis    string     `json:"synopsis"`
	Description string     `json:"description"`
	Options     [][]string `json:"options"`
	Examples    [][]string `json:"examples"`
}
type FileDir struct {
	Name        string
	Permissions string
	Owner       string
	ModifyDate  int64
	Size        int32
	IsDir       bool
}
type ByName []FileDir

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func formatBytes(size int32) string {
	units := []string{"K", "M", "G", "T", "P", "E", "Z", "Y", "R", "Q"}

	// Calculate the appropriate unit
	unitIndex := 0
	floatSize := float64(size)
	for floatSize >= 1024 && unitIndex < len(units)-1 {
		floatSize /= 1024
		unitIndex++
	}

	// Round the size to two decimal places
	roundedSize := math.Round(floatSize*10) / 10

	// Construct the human-readable string
	return fmt.Sprintf("%4.1f%s", roundedSize, units[unitIndex])
}

func handleInput(input string, gameState *pb.GameState) {
	tokens := strings.Split(strings.TrimSpace(input), " ")
	cmd := tokens[0]
	options := make(map[string]int)
	var positional []string
	if len(tokens) > 1 {
		// get options
		for j := 1; j < len(tokens); j++ {
			if strings.HasPrefix(tokens[j], "--") {
				options[tokens[j][2:]] = 1
				// todo some options have # with -t 5
				// grab from positional later on?
				// j += 1
			} else if strings.HasPrefix(tokens[j], "-") {
				if len(tokens[j]) == 2 {
					options[tokens[j][1:]] = 1
				} else {
					for _, flag := range tokens[j][1:] {
						options[string(flag)] = 1
					}
				}
			} else {
				positional = append(positional, tokens[j])
			}
		}
	}
	fmt.Println(tokens)
	switch cmd {
	case "":
		fmt.Print("")
	case "man":
		basicCmd.manCmd(tokens[1])
		if len(tokens) == 2 {
			file, err := os.Open(strings.Join([]string{"man/", tokens[1], ".json"}, ""))
			if err != nil {
				// no man page for given command
				fmt.Printf("bash: %s: command not found", tokens[0])
				return
			}
			defer file.Close()
			data, err := io.ReadAll(file)
			if err != nil {
				// Handle error
				os.Exit(5)
			}
			var manPage Man
			if err := json.Unmarshal(data, &manPage); err != nil {
				os.Exit(6)
			}
			// fmt.Printf("got: %+v\n", manPage)
			// fmt.Println(manPage.Description)
			// fmt.Printf("%s", data)
			// w, _ := screen.Size()
			// fmt.Println(w)
		}
	case "help":
		fmt.Println("should print something")
	case "whoami":
		fmt.Println(gameState.Username)
	case "cd":
		fmt.Print("")
		if len(tokens) == 1 {
			// change to root dir
			gameState.CurrentDir = gameState.Dirs.Name
		} else {
			dst := strings.Split(strings.TrimRight(tokens[1], "/"), "/")
			fmt.Print(dst)
		}
	case "pwd":
		fmt.Println(gameState.CurrentDir)
	case "ls":
		_, hiddenFlag := options["a"]
		_, hiddenBig := options["all"]
		_, long := options["l"]
		_, human := options["h"]
		if len(tokens) == 1 {
			contents := getContentNames(gameState.Dirs, false)
			sort.Strings(contents)
			fmt.Println(strings.Join(contents, "  "))
		} else {
			// TODO  details
			if long {
				// fmt.Println(strings.Join(getContentNames(gameState.Dirs, hiddenFlag || hiddenBig), "\n"))
				contents := getContents(gameState.Dirs, hiddenFlag || hiddenBig)
				sort.Sort(ByName(contents))
				for _, fileDir := range contents {
					d := "-"
					if fileDir.IsDir {
						d = "d"
					}
					var s string
					if human {
						s = formatBytes(fileDir.Size)
					} else {
						s = fmt.Sprintf("%7d", fileDir.Size)
					}
					// if fileDir.IsDir {
					// 	fmt.Printf("d%s %s %s", fileDir.Permissions, fileDir.Owner, fileDir.Owner)
					// } else {
					// 	fmt.Printf("-%s %s %s", fileDir.Permissions, fileDir.Owner, fileDir.Owner)
					// }
					fmt.Printf("%s%s %s %s %s MON DY HH:MM %s\n", d, fileDir.Permissions, fileDir.Owner, fileDir.Owner, s, fileDir.Name)
				}
			} else {
				contents := getContentNames(gameState.Dirs, hiddenFlag || hiddenBig)
				sort.Strings(contents)
				fmt.Println(strings.Join(contents, " "))
			}
		}
	}
}

func getContents(dir *pb.Directory, hidden bool) []FileDir {
	var result []FileDir
	for _, d := range dir.Dirs {
		if (d.Name[0] == '.' && hidden) || d.Name[0] != '.' {
			result = append(result, FileDir{
				Owner:       d.Owner,
				Name:        d.Name,
				Permissions: d.Permissions,
				ModifyDate:  d.ModifyDate,
				Size:        getSize(d),
				IsDir:       true,
			})
		}
	}
	for _, f := range dir.Files {
		if (f.Name[0] == '.' && hidden) || f.Name[0] != '.' {
			result = append(result, FileDir{
				Owner:       f.Owner,
				Name:        f.Name,
				Permissions: f.Permissions,
				ModifyDate:  f.ModifyDate,
				Size:        f.Size,
				IsDir:       false,
			})
		}
	}
	return result
}

func getSize(directory *pb.Directory) int32 {
	result := int32(4)
	for _, dir := range directory.Dirs {
		result += getSize(dir)
	}
	for _, file := range directory.Files {
		result += file.Size
	}
	return result
}
func getContentNames(dir *pb.Directory, hidden bool) []string {
	var s []string
	for _, d := range dir.Dirs {
		if (d.Name[0] == '.' && hidden) || d.Name[0] != '.' {
			s = append(s, d.Name)
		}
	}
	for _, f := range dir.Files {
		if (f.Name[0] == '.' && hidden) || f.Name[0] != '.' {
			s = append(s, f.Name)
		}
	}
	// todo sort
	return s
}

func GameLoop(gameState *pb.GameState) {
	var userInput string
	screen.Clear()
	// TODO
	// IDLE OS info
	// basic info, help, man, examples, etc

	// fmt.Printf("loaded: %+v\n", gameState)
	screen.MoveTopLeft()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("[%s@IDLE %s]$ ", gameState.Username, gameState.CurrentDir)
		scanner.Scan()
		userInput = scanner.Text()
		if userInput == "exit" || userInput == "quit" || userInput == "logout" {
			break
		}
		handleInput(userInput, gameState)
	}
}

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

func SaveProto(myProto *pb.GameState) error {
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

func AutoSave(gameState *pb.GameState) {
	var mu sync.Mutex
	for range time.Tick(10 * time.Second) {
		// fmt.Printf("should save: %+v\n", gameState)
		mu.Lock()
		gameCopy := gameState
		mu.Unlock()
		gameCopy.LastSave = time.Now().UnixMilli()

		err := SaveProto(gameCopy)
		if err != nil {
			os.Exit(2)
		}
	}
}

func LoadGame() (*pb.GameState, error) {
	// Open the encoded file for reading
	file, err := os.Open("data.isf")
	if err != nil {
		// File does not exist, new game
		var gameState pb.GameState
		gameState.Username = "root"
		gameState.Money = 0
		gameState.CurrentDir = "/"
		gameState.Dirs = &pb.Directory{
			Name:        "/",
			Files:       []*pb.File{},
			Permissions: "rwxrwxrwx",
			Owner:       "root",
			ModifyDate:  time.Now().UnixMilli(),
			Dirs: []*pb.Directory{
				{Name: "team", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "languages", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "tech", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "shop", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: "usr", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
				{Name: ".achievement", Owner: "root", Permissions: "rwxrwxrwx", ModifyDate: time.Now().UnixMilli(), Files: []*pb.File{}, Dirs: []*pb.Directory{}},
			},
		}
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
	gameState := &pb.GameState{}
	err = proto.Unmarshal([]byte(str), gameState)
	if err != nil {
		// Handle error
		return nil, err
	}
	// fmt.Printf("loaded: %+v\n", gameState)
	// time.Sleep(5 * time.Second)
	return gameState, nil
}

func main() {
	gameState, err := LoadGame()
	if err != nil {
		os.Exit(1)
	}
	go AutoSave(gameState)
	GameLoop(gameState)
}
