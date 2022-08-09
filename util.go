package main

import (
    "os"
    "log"
    "strconv"
    "strings"
    )

func iif (test bool, correct, wrong string) string {
	if test {
		return correct
	} else {
		return wrong
	}
}

func DisplayNumber() string {
    display := os.Getenv("DISPLAY")
    parts := strings.Split(display, ":")
    ret := "x"
    if len(parts)<2 {
        ret = parts[0]
    } else {
        ret = parts[1]
    }
    return ret
}

func Integer(raw string, def int) int {

    num, err := strconv.Atoi(raw)
    result := -1
    if err == nil {
        result = num
    } else {
        result = def
    }
    return result
}

func AppendFile(path string, text string) {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660);
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write([]byte(text + "\n")); err != nil {
        f.Close()
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}
