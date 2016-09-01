package godeps

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Parse(filename string) (*Godeps, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	deps := &Godeps{}

	if err := json.Unmarshal(data, deps); err != nil {
		return nil, err
	}

	return deps, nil
}
