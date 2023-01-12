package writers

import (
	"os"
)

func WriteToFile(fpath string, data string) error {
	encoded_data := []byte(data)
	err := os.WriteFile(fpath, encoded_data, 0644)

	return err
}
