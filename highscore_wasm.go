// +build js

package main

import (
	"syscall/js"
    "strconv"
)

var saveFilePath string = "turd_highscore"

func getHighScore() (uint64, error) {
    localStorage := js.Global().Get("localStorage")
    scoreStr := localStorage.Call("getItem", saveFilePath).String()

    // Check for the "<null>" string, indicating the key isn't set in localStorage.
    if scoreStr == "" || scoreStr == "<null>" {
        return 0, setHighScore(0)
    }

    // Convert the string representation of the score to a uint64
    score, err := strconv.ParseUint(scoreStr, 10, 64)
    if err != nil {
        return 0, err
    }
    return score, nil
}

func setHighScore(score uint64) error {
	localStorage := js.Global().Get("localStorage")
	
	// Convert the uint64 score to its string representation
	scoreStr := strconv.FormatUint(score, 10)
	localStorage.Call("setItem", saveFilePath, scoreStr)
	return nil
}
