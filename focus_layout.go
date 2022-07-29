 package main

import (
	"github.com/BurntSushi/xgb/xproto"
)

// FocusLayout displays all Clients as maximized windows.
type FocusLayout struct {
	clients []*Client
	setup map[string]string
}

func (l *FocusLayout) Name() string {
  return "Monocle"
}

// Arrange makes all clients in Workspace maximized.
func (l *FocusLayout) Arrange(w *Workspace) {
	side_count := len(l.clients)-1
	if side_count < 1 {
		side_count = 1
	}
	padding := 2
	screen_width := int(w.Screen.Width)
	screen_height := int(w.Screen.Height)
	side_height := screen_height / side_count
	for i, c := range l.clients {
		if i == 0 {//first take most of the view
			c.X = uint32(padding)
			c.Y = uint32(padding)
			c.W = uint32(int(w.Screen.Width)-100 - 2*padding)
			c.H = uint32(screen_height - (4*padding))
		} else {//stuff all the others on the side
			c.X = uint32(screen_width-100 + padding*2)
			c.Y = uint32( (i-1)*side_height + padding*2 - 1)
			c.W = uint32(100-(padding*4))
			c.H = uint32(side_height-(padding*4))
			c.BorderWidth = 1
			c.StackMode = xproto.StackModeAbove
		}
	}
}

// GetClients returns a slice of Client objects managed by this Layout.
func (l *FocusLayout) GetClients() []*Client {
	return append([]*Client{}, l.clients...)
}

func (l *FocusLayout) AddClient(c *Client) {
	l.clients = append(l.clients, c)
}

// RemoveClient removes a Client from the Layout.
func (l *FocusLayout) RemoveClient(c *Client) {
	for i, cc := range append([]*Client{}, l.clients...) {
		if c == cc {
			// Found client at at idx, so delete it and return.
			l.clients = append(l.clients[0:i], l.clients[i+1:]...)
			return
		}
	}
}

// MoveClient does nothing for the FocusLayout.
func (l *FocusLayout) MoveClient(_ *Client, d Direction) {
    last_client_idx := len(l.clients) -1
	if last_client_idx < 0 {
        return
    }

    switch d {
    case Left:
		//swap primary with last in side
		tmp := l.clients[0]
		l.clients[0] = l.clients[last_client_idx]
		l.clients[last_client_idx] = tmp

    case Right:
		//swap primary with first in side
		tmp := l.clients[0]
		l.clients[0] = l.clients[1]
		l.clients[1] = tmp

    case Up:
		//leave primary, rotate the side up
		l.clients = append(l.clients[0:1],				//first as list
			append(l.clients[2:],						//third to end
				l.clients[1])...)						//second 

    case Down:
		//leave primary, rotate the side down
		l.clients = append(l.clients[0:1],				//first as list
			append(l.clients[last_client_idx:],			//last as list
				l.clients[1:last_client_idx]...)...)	//second to end

    default:
        return
    }
}

func (l *FocusLayout) GetArrangment() map[string]string {
	if l.setup == nil {
		l.setup = make(map[string]string)
	}
	return l.setup
}

func (l *FocusLayout) SetArrangment(incoming_setup map[string]string) {
	l.setup = incoming_setup
}
