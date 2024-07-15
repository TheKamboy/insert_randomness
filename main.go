package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	huh "github.com/charmbracelet/huh"

	spinner "github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var name string

// Code from Stack Overflow
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

func huhtest() {
  fmt.Println("huh?")
  Pause()
}

func optionsMenu() {
	goback := false

	for !goback {
		CallClear()

		style := lipgloss.NewStyle().
			SetString("OPTIONS").
			Padding(1).
			Border(lipgloss.NormalBorder(), true).
			Align(lipgloss.Center)

		fmt.Println(style)
		fmt.Println("")

		var menu string
		huh.NewSelect[string]().
			Options(
				huh.NewOption("Set Name", "sn"),
				huh.NewOption("huh test", "ht"),
				huh.NewOption("Go Back", "exit"),
			).
			Height(6).
			Value(&menu).Run()

		if menu == "exit" {
			goback = true
		} else if menu == "sn" {
			SetName()

			style = lipgloss.NewStyle().SetString(fmt.Sprintf("Your name is now %v.", name)).Bold(true)

			fmt.Println(style)

      Pause()
		} else if menu == "ht" {
      huhtest()
    }
	}

	mainMenu()
}

func SetName() {
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
}

func Pause() {
	guh := func() {
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

func mainMenu() {
	CallClear()

	style := lipgloss.NewStyle().
		SetString("MAIN MENU").
		Padding(1).
		Border(lipgloss.NormalBorder(), true).
		Align(lipgloss.Center)

	fmt.Println(style)
	fmt.Println("")

	var menu string
	huh.NewSelect[string]().
		Options(
			huh.NewOption("Start Game", "sg"),
			huh.NewOption("Options", "o"),
			huh.NewOption("Exit", "exit"),
		).
		Height(6).
		Value(&menu).Run()

	if menu == "sg" {
		RoomTemplate()
	} else if menu == "o" {
		optionsMenu()
	} else if menu == "exit" {
		CallClear()
		os.Exit(0)
	}
}

func main() {
	style := lipgloss.NewStyle().
		SetString("INSERT RANDOMNESS").
		Padding(1).
		Border(lipgloss.ThickBorder(), true, true).
		Align(lipgloss.Center)

	fmt.Println(style)
	fmt.Println("")

	SetName()

	style = lipgloss.NewStyle().SetString(fmt.Sprintf("Your name is now %v.", name)).Bold(true)

	fmt.Println(style)

	Pause()

	mainMenu()
}
