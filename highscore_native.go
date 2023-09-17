// +build !js

package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

func getHighScore() (uint64, error) {
    file, err := os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil { return 0, err }
    defer file.Close()
    fileInfo, err := file.Stat()
    if err != nil { return 0, err }
    if fileInfo.Size() == 0 { return 0, setHighScore(0) }
    scoreBytes := [8]byte{}
    _, err = file.Read(scoreBytes[:])
    if err != nil { return 0, err }
    return binary.BigEndian.Uint64(scoreBytes[:]), nil
}

func setHighScore(score uint64) error {
	file, err := os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil { return err }
	defer file.Close()

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, score); err != nil {
		return err
	}

	_, err = file.WriteAt(buf.Bytes(), 0)
	return err
}
