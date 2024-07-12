package main

import (
	"errors"
	"fmt"
	"github.com/pkg/term"
	"strings"
	"os"
	"os/exec"
	"runtime"
  "log"

	huh "github.com/charmbracelet/huh"

	spinner "github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var name string

// Code from Stack Overflow
func getch() []byte {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	t.Restore()
	t.Close()
	if err != nil {
		return nil
	}
	return bytes[0:numRead]
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :( \nPlease make an issue.")
	}
}
// End of Stack Overflow Code

func Pause() {
  guh := func ()  {
    var char rune
    _, err := fmt.Scanf("%c", &char)
    if err != nil {
        log.Fatal(err)
    }
  }

  fmt.Println("")
  _ = spinner.New().Title("Press any key to continue...").Action(guh).Run()
}


func RoomTemplate() {
	fmt.Println("ho")
}

func main() {
	style := lipgloss.NewStyle().
		SetString("INSERT RANDOMNESS").
		Padding(1).
		Border(lipgloss.ThickBorder(), true, true).
		Align(lipgloss.Center)

	fmt.Println(style)
	fmt.Println("")

	huh.NewInput().
		Title("What's your name?").
		Description("This will be your character's name in the game.").
		Value(&name).
		// Validating fields is easy. The form will mark erroneous fields
		// and display error messages accordingly.
		Validate(func(str string) error {
			if strings.ReplaceAll(str, " ", "") == "" {
				return errors.New("Your name cannot be empty!")
			}
			return nil
		}).Run()

  style = lipgloss.NewStyle().SetString(fmt.Sprintf("Your name is now %v.", name)).Bold(true)

  fmt.Println(style)

  Pause()
}
