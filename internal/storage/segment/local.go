package segment

import (
	"os"
	"sync"

	"github.com/xitongsys/parquet-go/source"
)

type localParquet struct {
	file *os.File
	mu   sync.Mutex
}

func newLocalParquetWrite(fileName string) (source.ParquetFile, error) {
	return (&localParquet{}).Create(fileName)
}

func newLocalParquetReader(fileName string) (source.ParquetFile, error) {
	return (&localParquet{}).Open(fileName)
}

func (lp *localParquet) Create(fileName string) (source.ParquetFile, error) {
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	myFile := new(localParquet)
	myFile.file = f
	return myFile, err
}

func (lp *localParquet) Write(b []byte) (n int, err error) {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return lp.file.Write(b)
}

func (lp *localParquet) Close() error {
	return lp.file.Close()
}

func (lp *localParquet) Open(name string) (source.ParquetFile, error) {
	var err error
	myFile := new(localParquet)
	myFile.file, err = os.Open(name)
	return myFile, err
}

func (lp *localParquet) Seek(offset int64, pos int) (int64, error) {
	return lp.file.Seek(offset, pos)
}

func (lp *localParquet) Read(b []byte) (cnt int, err error) {
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
