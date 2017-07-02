package bcontext

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
)

// Out provide output api
type Out struct {
	isFileOpen bool
	hash       string
}

// File define output file
func (c *Out) File(path string) (string, error) {
	if c.isFileOpen {
		return "", fmt.Errorf("bcontext.ContextOut.Endile %s can not be opended. Don't close last file", path)
	}
	return "\n\n" + c.hash + services.StremDataSeparator + path + services.StremDataSeparator + "\n\n", nil
}

// EOF close output file
func (c *Out) EOF() (string, error) {
	if c.isFileOpen {
		return "", fmt.Errorf("bcontext.ContextOut.End: Close file (no file is open)")
	}
	return "\n\n--" + c.hash + "\n\n", nil
}
