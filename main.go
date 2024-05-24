package main

import (
	"log"
	"os/exec"
	"time"

	"go.i3wm.org/i3"
)

func main() {
	recv := i3.Subscribe(i3.WindowEventType)
	for recv.Next() {
		ev := recv.Event().(*i3.WindowEvent)
		if isFirefoxEnterFullscreenEvent(ev) {
			time.Sleep(100 * time.Millisecond)
			cmd := exec.Command("xdotool", "key", "--clearmodifiers", "F11")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
			log.Println("Sent cmd")
		}
	}
}

func isFirefoxEnterFullscreenEvent(event *i3.WindowEvent) bool {
	if event.Container.WindowProperties.Class != "firefox" {
		return false
	}

	if event.Change != "fullscreen_mode" {
		return false
	}

	workspaceNode := getFocusedWorkspaceNode()
	// if size is already equals -> exiting fullscreen
	if workspaceNode.Rect.X == event.Container.Rect.X && workspaceNode.Rect.Y == event.Container.Rect.Y {
		return false
	}

	return true
}

func getFocusedWorkspaceNode() *i3.Node {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	focusedWorkplace := tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Type == i3.WorkspaceNode
	})
	if focusedWorkplace == nil {
		log.Fatal("could not locate workspace")
	}

	return focusedWorkplace
}
