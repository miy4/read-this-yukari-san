package yonde

import (
	"bytes"
	"os"
	"os/exec"
)

const CID_YUKARI_V2 string = "2003"

func stripNewLines(lines []string) []byte {
	var newBuffer bytes.Buffer
	for _, s := range lines {
		for _, r := range s {
			if r == 0x0A || r == 0x0D {
				continue
			}
			newBuffer.WriteRune(r)
		}
	}

	return newBuffer.Bytes()
}

func saveToTempFile(msg []byte) (*os.File, error) {
	tmpfile, err := os.CreateTemp(os.TempDir(), "clse.tmp")
	if err != nil {
		return nil, err
	}

	if _, err = tmpfile.Write(msg); err != nil {
		return nil, err
	}

	if err = tmpfile.Close(); err != nil {
		return nil, err
	}

	return tmpfile, nil
}

func dispatchCmd(msg []byte) error {
	f, err := saveToTempFile(msg)
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	//_, err = exec.Command("SeikaSay2", "-cid", CID_YUKARI_V2, "-t", string(msg)).Output()
	_, err = exec.Command("SeikaSay2", "-cid", CID_YUKARI_V2, "-f", f.Name()).Output()
	return nil
}
