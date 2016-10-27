package itforarchivists

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/richardlehane/siegfried/config"
)

var current = Update{
	Version: config.Version(),
	Path:    "http://www.itforarchivists.com/siegfried/latest",
}

var version = [3]string{strconv.Itoa(current.Version[0]), strconv.Itoa(current.Version[1]), strconv.Itoa(current.Version[2])}

var github = "https://github.com/richardlehane/siegfried/releases/download/v"

func renderDownload(platform string) string {
	return fmt.Sprintf("%s%s.%s.%s/siegfried_%s-%s-%s_%s.zip", github, version[0], version[1], version[2], version[0], version[1], version[2], platform)
}

var downloads = struct {
	Version string
	Win64   string
	Win32   string
	//Linux64  string
	//Linux32  string
	//Darwin64 string
	//Darwin32 string
}{
	strings.Join(version[:], "."),
	renderDownload("win64"),
	renderDownload("win32"),
	//renderDownload("linux64"),
	//renderDownload("linux32"),
	//renderDownload("darwin64"),
	//renderDownload("darwin32"),
}
