package util

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"
)

// just to be sure, I know it can be simpler
func LoadSchema(dir string) (res string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	var (
		wg                sync.WaitGroup
		resC, errC, doneC = make(chan []byte), make(chan error), make(chan struct{})
		tmp               bytes.Buffer
	)

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".graphqls" {
			continue
		}
		wg.Add(1)
		go func(f fs.FileInfo) {
			defer wg.Done()
			s, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
			if err != nil {
				errC <- err
			} else {
				resC <- s
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(errC)
		close(resC)
		close(doneC)
	}()

	for {
		select {
		case err := <-errC:
			if err != nil {
				return "", err
			}
		case s := <-resC:
			if _, err = tmp.Write(s); err != nil {
				errC <- err
			}
		case <-doneC:
			return tmp.String(), nil
		}
	}
}
