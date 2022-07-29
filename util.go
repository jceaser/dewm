package main

func iif (test bool, correct, wrong string) string {
	if test {
		return correct
	} else {
		return wrong
	}
}
