package main

import (
	"errors"
	"strings"
	"syscall"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

/*
WM - Window manager
  screens- list
	clients - list
	workspaces - list
		layout
			clients - list

Window Manager
	workspaces
		screen
		layout
				clients
*/


var (
	errorQuit      = errors.New("Quit")
	errorAnotherWM = errors.New("Another WM already running")
	//pipeFile = "~/.dewm-commands"
	pipeFile = "/tmp/dewm-commands"
)

func (wm *WM) closeClientGracefully() error {
	if wm.activeClient == nil {
		log.Println("Tried to close client, but no active client")
		return nil
	}
	return wm.activeClient.CloseGracefully()
}

func (wm *WM) closeClientForcefully() error {
	if wm.activeClient == nil {
		log.Println("Tried to close client, but no active client")
		return nil
	}
	return wm.activeClient.CloseForcefully()
}

func DisplayFileName() string {
	display_number := DisplayNumber()
	pipeFile = fmt.Sprintf("%s.%s", pipeFile, display_number)
	return pipeFile
}

func readCommands(msg chan string, wm *WM) {

	display_number := DisplayNumber()
	pipeFile = fmt.Sprintf("%s.%s", pipeFile, display_number)
	fmt.Println (pipeFile)
	os.Remove(pipeFile)
	pipe_err := syscall.Mkfifo(pipeFile, 0660)
	if pipe_err != nil {
		log.Fatal("Make named pipe file error:", pipe_err)
	}
	defer os.Remove(pipeFile)

	fmt.Println ("Make pipe file")
	file, open_err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)
	if open_err != nil {
		log.Fatal("Open named pipe file error:", open_err)
	}
	fmt.Println ("command pipe open, waiting...")

	reader := bufio.NewReader(file)
	for {
		//fmt.Println ("reading command and putting into buffer")
		line, err := reader.ReadBytes('\n')
		if err == nil {
			//msg <- string(line)
			var action func() error
			action = nil
			clean_line := strings.TrimSpace(string(line))
			parts := strings.Split(clean_line, ":")
			cmd := clean_line
			option_str := ""
			option_num := -1
			if len(parts) < 2{
				cmd = parts[0]
			} else {
				cmd = parts[0]
				option_str = parts[1]
				option_num = Integer(parts[1], 0)
			}
			switch cmd {
				//case "close": closeClientGracefully()
				case "clean": wm.cleanupColumns()
				case "add": wm.addColumn()
				case "quit": wm.closeClientGracefully()
				case "exit": wm.closeClientForcefully()
				case "term": spawner("xterm")
				case "open": spawner(option_str)

				case "focus": wm.setLayoutOnActiveWorkspace(&FocusLayout{})
				case "monocle": wm.setLayoutOnActiveWorkspace(&MonocleLayout{})
				case "column": wm.setLayoutOnActiveWorkspace(&ColumnLayout{})

				case "tag": wm.MoveActiveClientToWorkspace(option_num)
				case "view": wm.SetActiveWorkspaceIdx(option_num)

				case "left": action = wm.moveClientOnActiveWorkspace(Left)
				case "down": action = wm.moveClientOnActiveWorkspace(Down)
				case "up": action = wm.moveClientOnActiveWorkspace(Up)
				case "right": action = wm.moveClientOnActiveWorkspace(Right)
				default: action = nil;
			}
			if action != nil {
				action()
			}
		}
	}
}

func print_manual() {
		fmt.Printf ("To use these commands, appended them to the input file:\n")
		fmt.Println (DisplayFileName())
		fmt.Printf ("\nor use the --command. Parameters are separated with the :\n\n")
		format := "%8s : %-7s %-61s\n"
		fmt.Printf(format, "Command", "Param", "Description")
		fmt.Printf(format, "--------", "-------", "------------")
		fmt.Printf(format, "add", "", "Add a new column")
		fmt.Printf(format, "clean", "", "Clean up columns")
		fmt.Printf(format, "quit", "", "Try to quit nicely")
		fmt.Printf(format, "exit", "", "Forcfully stop program")
		fmt.Printf(format, "term", "", "Open a XTerm")
		fmt.Printf(format, "open", "app", "Open <app>")
		fmt.Printf(format, "focus", "", "Switch to the focus view")
		fmt.Printf(format, "monical", "", "Switch to the monical view")
		fmt.Printf(format, "column", "", "Switch to th ecolumn view")
		fmt.Printf(format, "tag", "number", "Push the top app to screen <number>")
		fmt.Printf(format, "view", "number", "Switch to screen <number>")
		fmt.Printf(format, "left", "", "Move apps left")
		fmt.Printf(format, "down", "", "Move apps down")
		fmt.Printf(format, "up", "", "Move apps up")
		fmt.Printf(format, "right", "", "Move apps right")
}

func main() {
	//flag setup
    manual := flag.Bool("manual", false, "print out a command list")
    command := flag.String("command", "", "print out function manual")

	flag.Parse()

	if *manual {
		print_manual()
		return
	}

	if 0<len(*command) {
		pipeFile = DisplayFileName()
		AppendFile(pipeFile, *command)
		return
	}

	messages := make(chan string, 255)

	var wm = NewWM()
	wm_err := wm.Init()
	if wm_err != nil {
		log.Fatal(wm_err)
	}
	defer wm.Deinit()

	for i := 1; i < 9; i++ {
		fmt.Println("adding workspace")
		wm.AddWorkspace(&Workspace{Layout: &ColumnLayout{}})
	}

	go readCommands(messages, wm)

	for {
		fmt.Printf("Event loop\n")
		err := wm.handleEvent()
		switch err {
		case errorQuit:
			os.Exit(0)
		case nil:
		default:
			log.Print(err)
		}
		//fmt.Println ("about to get message")
		//fmt.Println (<-messages)
		//readCommands(file)
	}
}
