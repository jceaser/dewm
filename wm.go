package main

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
)

// WM holds the global window manager state.
type WM struct {
	xc *xgb.Conn

	xroot           xproto.ScreenInfo
	attachedScreens []xinerama.ScreenInfo

	grabs  []*Grab
	keymap [256][]xproto.Keysym

	clients      map[xproto.Window]*Client
	activeClient *Client

	workspaces []*Workspace
	activeWs   int
}

// NewWM allocates internal WM data structures and creates a WM
// instance. No X11 calls are made until WM.Init() is called.
func NewWM() *WM {
	return &WM{
		clients: map[xproto.Window]*Client{},
	}
}

// Init opens the X11 connections, performs the necessary calls to set
// up the internal WM state, and to start managing windows. You should
// also call Deinit before you exit.
func (wm *WM) Init() error {
	var err error
	wm.xc, err = xgb.NewConn()
	if err != nil {
		return err
	}

	if err = wm.initScreens(); err != nil {
		return err
	}
	if err = wm.initAtoms(); err != nil {
		return err
	}
	if err = wm.initWM(); err != nil {
		return err
	}
	if err = wm.initKeys(); err != nil {
		return err
	}
	if err = wm.initWorkspaces(); err != nil {
		return err
	}

	return nil
}

// Deinit cleans up internal WM state before exiting.
func (wm *WM) Deinit() {
	if wm.xc != nil {
		wm.xc.Close()
	}
}

func (wm *WM) initScreens() error {
	setup := xproto.Setup(wm.xc)
	if setup == nil || len(setup.Roots) < 1 {
		return errors.New("Could not parse SetupInfo.")
	}
	if err := xinerama.Init(wm.xc); err != nil {
		return err
	}
	if r, err := xinerama.QueryScreens(wm.xc).Reply(); err != nil {
		return err
	} else {
		if len(r.ScreenInfo) == 0 {
			// If Xinerama does not return useful information, we can
			// still query the root window, and create a fake
			// ScreenInfo structure.
			wm.attachedScreens = []xinerama.ScreenInfo{
				xinerama.ScreenInfo{
					Width:  setup.Roots[0].WidthInPixels,
					Height: setup.Roots[0].HeightInPixels,
				},
			}
		} else {
			wm.attachedScreens = r.ScreenInfo
		}
	}

	coninfo := xproto.Setup(wm.xc)
	if coninfo == nil {
		return errors.New("Could not parse X connection info")
	}
	if len(coninfo.Roots) != 1 {
		return errors.New("Bad number of roots. Did Xinerama initialize correctly?")
	}
	wm.xroot = coninfo.Roots[0]
	return nil
}

func (wm *WM) initWM() error {
	err := xproto.ChangeWindowAttributesChecked(
		wm.xc,
		wm.xroot.Root,
		xproto.CwEventMask,
		[]uint32{
			xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease |
				xproto.EventMaskButtonPress |
				xproto.EventMaskButtonRelease |
				xproto.EventMaskStructureNotify |
				xproto.EventMaskSubstructureRedirect,
		},
	).Check()
	if err != nil {
		if _, ok := err.(xproto.AccessError); ok {
			return errorAnotherWM
		}
	}
	return err
}

func (wm *WM) initWorkspaces() error {
	tree, err := xproto.QueryTree(wm.xc, wm.xroot.Root).Reply()
	if err != nil {
		return err
	}
	if tree == nil {
		return nil
	}
	w := &Workspace{
		Layout: &ColumnLayout{},
	}
	wm.AddWorkspace(w)
	for _, win := range tree.Children {
		if wm.GetClient(win) != nil {
			fmt.Errorf("Window already managed by a client - what happened?")
		}
		c := NewClient(wm.xc, win)
		err := c.Init()
		if err != nil {
			return err
		}
		wm.AddClient(c)
		w.AddClient(c)
	}
	w.Arrange()
	return nil
}

// AddClient adds the client to WM's internal client list.
func (wm *WM) AddClient(c *Client) {
	w := c.window // private!
	wm.clients[w] = c
}

// GetClient returns the Client from associated Window ID, or nil.
func (wm *WM) GetClient(w xproto.Window) *Client {
	c, _ := wm.clients[w]
	return c
}

// AddWorkspace appends the given Workspace at the end of the list,
// and attaches it to the first screen.
func (wm *WM) AddWorkspace(w *Workspace) {
	wm.workspaces = append(wm.workspaces, w)
	w.Screen = &wm.attachedScreens[0]
}

// GetActiveWorkspace returns the Workspace containing the current
// active Client, or nil if no Client is active.
func (wm *WM) GetActiveWorkspace() *Workspace {
	w := wm.workspaces[wm.activeWs]
	if wm.activeClient != nil && !w.HasClient(wm.activeClient) {
		fmt.Errorf("marked active but don't have the active client")
	}
	return w
}

func (wm *WM) valid_workspace(i int) bool {
	if i<0 || i >= len(wm.workspaces) {
		fmt.Printf("only %d workspaces.\n", len(wm.workspaces))
		return false
	}
	return true
}

func (wm *WM) MoveActiveClientToWorkspace(i int) error {
	var ret error = nil
	source := wm.GetActiveWorkspace()
	dest := wm.workspaces[i]
	if source.RemoveClient(wm.activeClient) {
		dest.AddClient(wm.activeClient)
		ret = wm.SetActiveWorkspaceIdx(i)
		source.Arrange()
		dest.Arrange()
		fmt.Println("done arranging")
	} else {
		ret = errors.New("Could not parse SetupInfo.")
	}
	fmt.Print ("error: ")
	fmt.Println (ret)
	return ret
}

// SetActiveWorkspaceIdx switches to the given workspace (by index).
func (wm *WM) SetActiveWorkspaceIdx(i int) error {
	fmt.Printf ("Setting active workspace from %d to %d.\n", wm.activeWs, i)
	if !wm.valid_workspace(i) {
		fmt.Println ("Not a valid workspace")
		return nil
	}
	fmt.Printf("workspace %d -> %d.\n", wm.activeWs, i)
	if wm.activeWs == i {
		fmt.Println ("already on %d", i)
		return nil
	}
	if err := wm.workspaces[wm.activeWs].Hide(); err != nil {
		fmt.Println ("Hide failed:", err)
		return err
	}
	//wm.workspaces[wm.activeWs].Arrange()
	wm.activeWs = i
	if err := wm.workspaces[wm.activeWs].Show(); err != nil {
		fmt.Println ("Show failed:", err)
		return err
	}
	wm.workspaces[wm.activeWs].Arrange()

	return nil
}
