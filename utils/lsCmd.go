package utils

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	pb "github.com/depperm/idleOs/proto"
)

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

func HasExe(gameState *GameState, exe string) bool {
	contents := getContentNames(gameState.Player.Dirs, false)
	for _, file := range contents {
		if file == fmt.Sprintf("%s.exe", exe) {
			return true
		}
	}
	return false
}

func LsCmd(gameState *GameState, tokens []string, options map[string]int) {
	_, hiddenFlag := options["a"]
	_, hiddenBig := options["all"]
	_, long := options["l"]
	_, human := options["h"]
	if len(tokens) == 1 {
		contents := getContentNames(gameState.Player.Dirs, false)
		sort.Strings(contents)
		fmt.Println(strings.Join(contents, "  "))
	} else {
		// TODO  details
		if long {
			// fmt.Println(strings.Join(getContentNames(gameState.Dirs, hiddenFlag || hiddenBig), "\n"))
			contents := getContents(gameState.Player.Dirs, hiddenFlag || hiddenBig)
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
				currentTime := time.Unix(0, fileDir.ModifyDate*int64(time.Millisecond))
				if currentTime.Year() != time.Unix(0, time.Now().UnixMilli()*int64(time.Millisecond)).Year() {
					fmt.Printf("%s%s %s %s %s %s %02d HH:MM %s\n",
						d,
						fileDir.Permissions,
						fileDir.Owner,
						fileDir.Owner,
						s,
						currentTime.Month().String()[:3],
						currentTime.Day(),
						fileDir.Name)
				} else {
					fmt.Printf("%s%s %s %s %s %s %02d HH:MM %s\n",
						d,
						fileDir.Permissions,
						fileDir.Owner,
						fileDir.Owner,
						s,
						currentTime.Month().String()[:3],
						currentTime.Day(),
						fileDir.Name)
				}
			}
		} else {
			contents := getContentNames(gameState.Player.Dirs, hiddenFlag || hiddenBig)
			sort.Strings(contents)
			fmt.Println(strings.Join(contents, " "))
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
