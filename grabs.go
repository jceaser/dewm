package main

import (
	//"fmt"
	"errors"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
)

// Grab represents a key grab and its callback
type Grab struct {
	sym       xproto.Keysym
	modifiers uint16
	codes     []xproto.Keycode
	callback  func() error
}

func (wm *WM) getGrabs() []*Grab {
	return []*Grab{
		{
			sym:       XK_q,
			modifiers: xproto.ModMaskControl | xproto.ModMaskShift | xproto.ModMask1,
			callback:  func() error { return errorQuit },
		},
		{
			sym:       XK_Return,
			modifiers: xproto.ModMask1,
			callback: spawner("xterm"),
		},
		{
			sym:       XK_Return,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback:  spawner("st"),
		},
		{
			sym:       XK_q,
			modifiers: xproto.ModMask1,
			callback:  wm.closeClientGracefully,
		},
		{
			sym:       XK_q,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback:  wm.closeClientForcefully,
		},

		{
			sym:       XK_Left,
			modifiers: xproto.ModMask1,
			callback:  wm.moveClientOnActiveWorkspace(Left),
		},
		{
			sym:       XK_Right,
			modifiers: xproto.ModMask1,
			callback:  wm.moveClientOnActiveWorkspace(Right),
		},

		{
			sym:       XK_Down,
			modifiers: xproto.ModMask1,
			callback:  wm.moveClientOnActiveWorkspace(Down),
		},
		{
			sym:       XK_Up,
			modifiers: xproto.ModMask1,
			callback:  wm.moveClientOnActiveWorkspace(Up),
		},

		{
			sym:       XK_d,
			modifiers: xproto.ModMask1,
			callback:  wm.cleanupColumns,
		},
		{
			sym:       XK_n,
			modifiers: xproto.ModMask1,
			callback:  wm.addColumn,
		},
		{
			sym:	   XK_f,
			modifiers: xproto.ModMask1,
			callback: func() error {
				return wm.setLayoutOnActiveWorkspace(&FocusLayout{})
			},
		},
		{
			sym:       XK_m,
			modifiers: xproto.ModMask1,
			callback: func() error {
				return wm.setLayoutOnActiveWorkspace(&MonocleLayout{})
			},
		},
		{
			sym:       XK_c,
			modifiers: xproto.ModMask1,
			callback: func() error {
				return wm.setLayoutOnActiveWorkspace(&ColumnLayout{})
			},
		},
		{
			sym:		XK_1,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(0)},
		},
		{
			sym:		XK_2,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(1)},
		},
		{
			sym:		XK_3,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(2)},
		},
		{
			sym:		XK_4,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(3)},
		},
		{
			sym:		XK_5,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(4)},
		},
		{
			sym:		XK_6,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(5)},
		},
		{
			sym:		XK_7,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(6)},
		},
		{
			sym:		XK_8,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(7)},
		},
		{
			sym:		XK_9,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(8)},
		},
		{
			sym:		XK_0,
			modifiers: xproto.ModMask1 | xproto.ModMaskShift,
			callback: func() error {return wm.MoveActiveClientToWorkspace(9)},
		},
		{
			sym:       XK_1,
			modifiers: xproto.ModMask1,
			callback:  func() error {return wm.SetActiveWorkspaceIdx(0)},
		},
		{
			sym:       XK_2,
			modifiers: xproto.ModMask1,
			callback:  func() error {return wm.SetActiveWorkspaceIdx(1) },
		},
		{
			sym:       XK_3,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(2) },
		},
		{
			sym:       XK_4,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(3) },
		},
		{
			sym:       XK_5,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(4) },
		},
		{
			sym:       XK_6,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(5) },
		},
		{
			sym:       XK_7,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(6) },
		},
		{
			sym:       XK_8,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(7) },
		},
		{
			sym:       XK_9,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(8) },
		},
		{
			sym:       XK_0,
			modifiers: xproto.ModMask1,
			callback:  func() error { return wm.SetActiveWorkspaceIdx(9) },
		},

	}
}

func (wm *WM) initKeys() error {
	const (
		loKey = 8
		hiKey = 255
	)

	m := xproto.GetKeyboardMapping(wm.xc, loKey, hiKey-loKey+1)
	reply, err := m.Reply()
	if err != nil {
		return err
	}
	if reply == nil {
		return errors.New("Could not load keyboard map")
	}

	for i := 0; i < hiKey-loKey+1; i++ {
		wm.keymap[loKey+i] = reply.Keysyms[i*int(reply.KeysymsPerKeycode) : (i+1)*int(reply.KeysymsPerKeycode)]
	}

	wm.grabs = wm.getGrabs()

	for i, syms := range wm.keymap {
		for _, sym := range syms {
			for c := range wm.grabs {
				if wm.grabs[c].sym == sym {
					wm.grabs[c].codes = append(wm.grabs[c].codes, xproto.Keycode(i))
				}
			}
		}
	}
	for _, grabbed := range wm.grabs {
		for _, code := range grabbed.codes {
			if err := xproto.GrabKeyChecked(
				wm.xc,
				false,
				wm.xroot.Root,
				grabbed.modifiers,
				code,
				xproto.GrabModeAsync,
				xproto.GrabModeAsync,
			).Check(); err != nil {
				return err
			}

		}
	}
	return nil
}

func (wm *WM) cleanupColumns() error {
	w := wm.GetActiveWorkspace()
	switch l := w.Layout.(type) {
	case *ColumnLayout:
		l.cleanupColumns()
	}
	return w.Arrange()
}

func (wm *WM) addColumn() error {
	w := wm.GetActiveWorkspace()
	switch l := w.Layout.(type) {
	case *ColumnLayout:
		l.addColumn()
	}
	return w.Arrange()
}

func (wm *WM) setLayoutOnActiveWorkspace(l Layout) error {
	w := wm.GetActiveWorkspace()
	w.SetLayout(l)
	return w.Arrange()
}

func (wm *WM) moveClientOnActiveWorkspace(d Direction) func() error {
	return func() error {
		w := wm.GetActiveWorkspace()
		if wm.activeClient == nil {
			return nil
		}
		w.Layout.MoveClient(wm.activeClient, d)
		return w.Arrange()
	}
}

func spawner(cmd string, args ...string) func() error {
	return func() error {
		go func() {
			cmd := exec.Command(cmd, args...)
			if err := cmd.Start(); err == nil {
				cmd.Wait()
			}
		}()
		return nil
	}
}
