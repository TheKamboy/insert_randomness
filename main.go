package main

import (
	"errors"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	huh "github.com/charmbracelet/huh"

	spinner "github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

// debug mode will auto set name, and other things to make it quicker to test
const DEBUG = true

// Enables Nerd Fonts
const NERD = true

// Character Name
var name string

// Inventory Variables

// Terminal Sizing
var termWidth = 0

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

func convertToSymb(symbol rune) string {
	symb := symbol
	foreground := lipgloss.Color("8")

	if NERD {
		if symb == 'O' {
			symb = ''
		}
	}

	style := lipgloss.NewStyle().Foreground(foreground)

	if NERD {
		return style.Render(string(symb))
	} else {
		return style.Render("")
	}
}

// Text Adventure input stuff (pi in this case means Processed Input, oi means original input)
// If it fails to find a command, pi will be used to show the command they sent
// Commands are in 2 parts
//
// examine object
// ^^^^^^^ ^^^^^^
// Verb    Noun
//
// Nouns have to be 1 word, or else I will need to code in more things for handling more than 2 words
func playerInput(roomTitle string) (pi string, oi string, err string) {
	huh.NewInput().
		Title(roomTitle).
		Value(&oi).
		// Validating fields is easy. The form will mark erroneous fields
		// and display error messages accordingly.
		Validate(func(str string) error {
			if strings.ReplaceAll(str, " ", "") == "" {
				return errors.New("empty input")
			}
			return nil
		}).Run()

	oi = strings.ToLower(oi)

	before, after, found := strings.Cut(oi, " ")

	// TODO: Finish commands
	if found {
		if before == "examine" || before == "look" {
			pi = fmt.Sprintf("exam %v", after)
		}
	} else {
		if oi == "look" {
			pi = "cls"
		} else if oi == "examine" {
			pi = oi
			err = "\"examine\" requires an noun. (if you are trying to clear the screen use \"look\")"
		} else if oi == "quit" || oi == "exit" {
			pi = "exit"
		} else {
			pi = oi
			err = fmt.Sprintf("Failed to process command: %v", pi)
		}
	}

	return
}

func debugMsg() {
	var style = lipgloss.NewStyle()

	if DEBUG {
		fmt.Println(style.Render("DEBUG MODE"))
	}
}

func huhtest() {
	fmt.Println("huh?")
	Pause()
}

func optionsMenu() {
	goback := false

	for !goback {
		CallClear()
		debugMsg()

		style := lipgloss.NewStyle().
			SetString(fmt.Sprintf("%v OPTIONS", convertToSymb('O'))).
			Padding(1).
			Border(lipgloss.RoundedBorder(), true).
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
				return errors.New("your name cannot be empty")
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

	// damn charm got lazy with the adaptive colors
	_ = spinner.New().TitleStyle(lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#00020A", Dark: "#FFFDF5"})).Title("Press any key to continue...").Action(guh).Run()
}

func RoomTemplate() {
	arrowInputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F780E2"))
	inputStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"})
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"})

	CallClear()
	for {
		// Get Input
		pi, oi, err := playerInput("This is a room template")

		// Display Input
		fmt.Println(arrowInputStyle.Render("> ") + inputStyle.Render(oi))

		// Detect Error
		if err != "" {
			fmt.Println(errorStyle.Render(" * " + err))
		} else {
			fmt.Println(pi)
		}

		// Use Input
		// TODO: Finish making commands
		if oi == "quit" {
			os.Exit(0)
		}

		if pi == "exit" {
			os.Exit(0)
		}

    if pi == "cls" {
      CallClear()
    }

		fmt.Println("")
	}
}

func mainMenu() {
	CallClear()
	debugMsg()

	style := lipgloss.NewStyle().
		SetString("MAIN MENU").
		Padding(1).
		Border(lipgloss.RoundedBorder(), true).
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
	w, _, e := term.GetSize(int(os.Stdin.Fd()))

	if e != nil {
		log.Fatalln(e)
	}

	termWidth = w

	style := lipgloss.NewStyle().
		SetString("INSERT RANDOMNESS").
		Padding(1).
		Border(lipgloss.ThickBorder(), true, true).
		Align(lipgloss.Center)

	fmt.Println(style)
	fmt.Println("")

	if !DEBUG {
		SetName()

		style = lipgloss.NewStyle().SetString(fmt.Sprintf("Your name is now %v.", name)).Bold(true)

		fmt.Println(style)

		Pause()
	} else {
		name = "Debug"
	}

	mainMenu()
}
