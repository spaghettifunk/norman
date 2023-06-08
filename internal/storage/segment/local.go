package segment

import (
	"os"
	"path/filepath"
	"sync"
)

type LocalParquet struct {
	file *os.File
	mu   sync.Mutex
}

func NewLocalParquet(dir, filename string) (*LocalParquet, error) {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	myFile := new(LocalParquet)

	// Create the file if it doesn't exist
	filePath := filepath.Join(dir, filename)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		myFile.file, err = os.Create(filePath)
		if err != nil {
			return nil, err
		}
	}
	return myFile, err
}

func (lp *LocalParquet) Write(b []byte) (n int, err error) {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return lp.file.Write(b)
}

func (lp *LocalParquet) Close() error {
	return lp.file.Close()
}

func (lp *LocalParquet) Open(name string) (*LocalParquet, error) {
	var err error
	myFile := new(LocalParquet)
	myFile.file, err = os.Open(name)
	return myFile, err
}

func (lp *LocalParquet) Seek(offset int64, pos int) (int64, error) {
	return lp.file.Seek(offset, pos)
}

func (lp *LocalParquet) Read(b []byte) (cnt int, err error) {
	var n int
	ln := len(b)
	for cnt < ln {
		n, err = lp.file.Read(b[cnt:])
		cnt += n
		if err != nil {
			break
		}
	}
	return cnt, err
}
