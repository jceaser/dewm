package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ColumnLayout arranges Clients into columns.
type ColumnLayout struct {
	columns [][]*Client
	setup map[string]string
}

func (l *ColumnLayout) Name() string {
	return "Column"
}

// Arrange arranges all the windows of the workspace into the screen
// that the workspace is attached to.
func (l *ColumnLayout) Arrange(w *Workspace) {
	nColumns := uint32(len(l.columns))

	// If there are no columns, create one.
	if nColumns == 0 {
		l.addColumn()
		nColumns++
	}

	colWidth := uint32(w.Screen.Width) / nColumns
	for columnIdx, column := range l.columns {
		colHeight := uint32(w.Screen.Height)
		for rowIdx, client := range column {
			client.X = uint32(columnIdx) * colWidth
			client.Y = uint32(uint32(rowIdx) * (colHeight / uint32(len(column))))
			client.W = colWidth
			client.H = uint32(colHeight / uint32(len(column)))
		}
	}
}

// GetClients returns a slice of Client objects managed by this Layout.
func (l *ColumnLayout) GetClients() []*Client {
	clients := make(
		[]*Client,
		0,
		len(l.columns)*3, // reserve some extra capacity
	)
	for _, column := range l.columns {
		clients = append(clients, column...)
	}
	return clients
}

func (l *ColumnLayout) AddClient(c *Client) {
	// No columns? Add one
	if len(l.columns) == 0 {
		l.addColumn()
	}
	// First, look for an empty column to put the client in.
	for i, column := range l.columns {
		if len(column) == 0 {
			l.columns[i] = append(l.columns[i], c)
			return
		} else {
			//rows exist, try to fill all the empty ones first
			for row, value := range l.columns[i] {
				if value == nil {
					l.columns[i][row] = c
					return
				}
			}
		}
	}
	// Failing that, cram the client in the last column.
	l.columns[len(l.columns)-1] = append(l.columns[len(l.columns)-1], c)
}

// RemoveClient removes a Client from the Layout.
func (l *ColumnLayout) RemoveClient(c *Client) {
	for colIdx, column := range l.columns {
		for clIdx, cc := range append([]*Client{}, column...) {
			if c == cc {
				// Found client at at clIdx, so delete it and return.
				l.columns[colIdx] = append(
					column[0:clIdx],
					column[clIdx+1:]...,
				)
				return
			}
		}
	}
}

func (l *ColumnLayout) cleanupColumns() {
restart:
	for {
		for i, c := range l.columns {
			if len(c) == 0 {
				l.columns = append(l.columns[0:i], l.columns[i+1:]...)
				continue restart
			}
		}
		return
	}
}

func (l *ColumnLayout) addColumn() {
	l.columns = append(l.columns, []*Client{})
}

// MoveClient moves the client left/right between columns, or up/down
// change rows within a single column.
func (l *ColumnLayout) MoveClient(c *Client, d Direction) {
	switch d {
	case Up:
		fallthrough
	case Down:
		idx := d.V
		for _, column := range l.columns {
			for i, cc := range column {
				if c == cc {
					// got ya
					if i == 0 && idx < 0 {
						return
					}
					if i == (len(column)-1) && idx > 0 {
						return
					}
					column[i], column[i+idx] = column[i+idx], column[i]
					return
				}
			}
		}

	case Left:
		fallthrough
	case Right:
		idx := d.H
		for colIdx, column := range l.columns {
			for clIdx, cc := range column {
				if c == cc {
					// got ya
					if colIdx == 0 && idx < 0 {
						return
					}
					if colIdx == (len(l.columns)-1) && idx > 0 {
						return
					}
					l.columns[colIdx] = append(
						column[0:clIdx],
						column[clIdx+1:]...,
					)
					l.columns[colIdx+idx] = append(
						l.columns[colIdx+idx],
						c,
					)
					return
				}
			}
		}

	default:
		return
	}
}

/**
Generate a string such as: '3,4'
for the case of 2 columns, with 3 and 4 rows
*/
func (l *ColumnLayout) GetArrangment() map[string]string {
	details := ""
	for idx, value := range l.columns {
		details = fmt.Sprintf ("%s%s%d",
			details,
			iif(idx>0, ",", ""),
			len(value))
	}
	if l.setup == nil {
		l.setup = make (map[string]string)
	}
	l.setup[l.Name()] = details
	return l.setup
}

func (l *ColumnLayout) SetArrangment(incoming_setup map[string]string) {
	data := incoming_setup[l.Name()]
	parts := strings.Split(data, ",")
	for idx, raw_value := range parts {
		value, _ := strconv.ParseInt(raw_value, 0, 64)
		l.addColumn()
		for n:=int64(0) ; n<value; n++ {
			//fill array with blanks to be populated latter
			l.columns[idx] = append(l.columns[idx], nil)
		}
	}
	l.setup = incoming_setup
}
