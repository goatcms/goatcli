package builder

import (
	"io"
	"strings"
	"testing"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

const (
	testWriterHash   = `01234567890`
	testWriterStream = "\n\n01234567890:path/main.md:\n\nContent\n\n--01234567890\n\n" + "\n\n01234567890:path/main2.md:\n\nContent2\n\n--01234567890\n\n"
)

func TestWriter(t *testing.T) {
	var (
		err     error
		fs      filesystem.Filespace
		writer  *FSWriter
		reader  io.Reader
		content []byte
	)
	t.Parallel()
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	writer = NewFSWriter(fs, testWriterHash)
	reader = strings.NewReader(testWriterStream)
	if _, err = io.Copy(writer, reader); err != nil {
		t.Error(err)
		return
	}
	if !fs.IsFile("path/main.md") {
		t.Errorf("path/main.md file is required")
		return
	}
	if content, err = fs.ReadFile("path/main.md"); err != nil {
		t.Error(err)
		return
	}
	if string(content) != "Content" {
		t.Errorf("file content should be equals to 'Content' and it is '%s'", content)
		return
	}
	if !fs.IsFile("path/main2.md") {
		t.Errorf("path/main2.md file is required")
		return
	}
	if content, err = fs.ReadFile("path/main2.md"); err != nil {
		t.Error(err)
		return
	}
	if string(content) != "Content2" {
		t.Errorf("file2 content should be equals to 'Content2' and it is '%s'", content)
		return
	}
}
