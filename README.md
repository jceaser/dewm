# dewm
[Forked](https://github.com/rollcat/dewm)

`dewm` is a pure Go autotiling window manager. You may find it
somewhat similar to [dwm][] or [wmii][], but has some ideas of its
own.

This `dewm` was forked from Kamil (Rollcat)'s which in turn was forked from Dave MacFarlane's [dewm][original-dewm],
which was written in [literate style][literate-programming], using
[lmt][]. The fork dropped the original Markdown sources, heavy
refactoring and cleanup was done, bugs were fixed, some features
dropped, more added, arbitrary changes made.

[original-dewm]: https://github.com/driusan/dewm
[fork]: https://github.com/rollcat/dewm
[literate-programming]: https://en.wikipedia.org/wiki/Literate_programming
[lmt]: https://github.com/driusan/lmt
[dwm]: https://dwm.suckless.org/
[wmii]: https://code.google.com/archive/p/wmii/

## Basics

`dewm` is a very simple window manager that seeks to arranges winows on a screen automaticly freeing the user from the need to keep track of or position windows. Most users probobly never has experienced a window manager like this being use to needing to move and size windows manually. In fact, the current state of computing has resulted in most users having so many windows open on the screen that windows actually get lost.

Windows are arranged using three different layouts, Column, Focus, and Monical. 

* Column - All windows displayed in one vertical column
* Focus - Primary window taking most of the screen on the left with a column of other windows on the right.
* Monical - Stacks all the apps on top of one another with the primare window on top. All windows are full screen.

Keybindings then allow you to move windows around in layout specific ways.

## Keybindings

These keybindings are currently hardcoded, but may one day be configurable. Note, if running on Mac OS X, be sure to set the option key to act as an alt key in XQuarts preferences.

### Window Management

* `Alt-right/Alt-left/Alt-Up/Alt-Down` move the current window position.
* `Alt-N` create a new column 
* `Alt-D` delete any empty columns

### Layouts

* `Alt-c` - switch to column layout
* `Alt-f` - switch to focus layout
* `Alt-m` - switch to monical layout

### Other

* `Alt-Enter` spawn an xterm
* `Alt-Shift-Enter` spawn an st
* `Alt-Q` close the current window
* `Alt-Shift-Q` destroy the current window
* `Ctrl-Alt-Shift-Q` quit dewm
