 package main

import (
	"github.com/BurntSushi/xgb/xproto"
)

// MonocleLayout displays all Clients as maximized windows.
type MonocleLayout struct {
    clients []*Client
    setup map[string]string
}

func (l *MonocleLayout) Name() string {
  return "Monocle"
}

// Arrange makes all clients in Workspace maximized.
func (l *MonocleLayout) Arrange(w *Workspace) {
    for idx, c := range l.clients {
		c.X = 0
		c.Y = 0
		c.W = uint32(w.Screen.Width)
		c.H = uint32(w.Screen.Height)
		c.BorderWidth = 0
        if idx == 0 {
            c.StackMode = xproto.StackModeAbove
        } else {
            c.StackMode = xproto.StackModeBelow
        }
    }
}

// GetClients returns a slice of Client objects managed by this Layout.
func (l *MonocleLayout) GetClients() []*Client {
	return append([]*Client{}, l.clients...)
}

func (l *MonocleLayout) AddClient(c *Client) {
	l.clients = append(l.clients, c)
}

// RemoveClient removes a Client from the Layout.
func (l *MonocleLayout) RemoveClient(c *Client) {
	for i, cc := range append([]*Client{}, l.clients...) {
		if c == cc {
			// Found client at at idx, so delete it and return.
			l.clients = append(l.clients[0:i], l.clients[i+1:]...)
			return
		}
	}
}

func rotateLeft(arr []*Client, d int) []*Client {
    for ; d > 0 ; d-- {
        left := arr[0]
        arr = arr[1:]
        arr = append(arr, left)
    }
    return arr
}

func rotateRight(ar []*Client,d,n int) []*Client{
    ar = append(ar[d:n],ar[0:d]...)
    return  ar
}

// MoveClient does nothing for the MonocleLayout.
func (l *MonocleLayout) MoveClient(_ *Client, d Direction) {
    if len(l.clients)<1 {
        return
    }

    switch d {
    case Right:
        l.clients = rotateRight(l.clients, 1, len(l.clients))

    case Left:
        l.clients = rotateLeft(l.clients, 1)

    case Up:
        fallthrough
    case Down:
        // windows are stacked from first to end, swap the ends
        last := len(l.clients)-1
        temp := l.clients[last-1]
        l.clients[last-1] = l.clients[last]
        l.clients[last] = temp

    default:
        return
    }

}

func (l *MonocleLayout) GetArrangment() map[string]string {
	if l.setup == nil {
		l.setup = make(map[string]string)
	}
	return l.setup
}

func (l *MonocleLayout) SetArrangment(incoming_setup map[string]string) {
    l.setup = incoming_setup
}