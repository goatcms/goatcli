package builder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
)

// A FSWriter is the parser and stream to filesystem
type FSWriter struct {
	fs           filesystem.Filespace
	buffor, hash string
	openTag      *regexp.Regexp
	closeTag     *regexp.Regexp
}

// NewFSWriter create new filesystem write stream
func NewFSWriter(fs filesystem.Filespace, hash string) *FSWriter {
	return &FSWriter{
		fs:       fs,
		hash:     hash,
		openTag:  regexp.MustCompile("\n\n" + hash + services.StremDataSeparator + "[A-Za-z0-9/_\\.-]+" + services.StremDataSeparator + "\n\n"),
		closeTag: regexp.MustCompile("\n\n--" + hash + "\n\n"),
	}
}

// Write new data to filesystem writer
func (w *FSWriter) Write(data []byte) (n int, err error) {
	var closeLoc, openLoc []int
	w.buffor += string(data)
	for {
		if closeLoc = w.closeTag.FindStringIndex(w.buffor); closeLoc == nil {
			return len(data), nil
		}
		if openLoc = w.openTag.FindStringIndex(w.buffor); openLoc == nil {
			return 0, fmt.Errorf("builder.FSWriter.Write: Close unknow block: %v", w.buffor)
		}
		content := w.buffor[openLoc[1]:closeLoc[0]]
		openStr := w.buffor[openLoc[0]+len(w.hash)+3 : openLoc[1]]
		path := openStr[:strings.Index(openStr, services.StremDataSeparator)]
		if err = w.fs.WriteFile(path, []byte(content), 0766); err != nil {
			return 0, err
		}
		w.buffor = w.buffor[closeLoc[1]:]
	}
}

// Close writer
func (w *FSWriter) Close() error {
	if openLoc := w.openTag.FindStringIndex(w.buffor); openLoc != nil {
		return fmt.Errorf("builder.FSWriter.Close: Close writer with open tag: %v", w.buffor[openLoc[0]:openLoc[1]])
	}
	return nil
}
