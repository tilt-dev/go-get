package get_test

import (
	"fmt"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/tilt-dev/go-get"
)

func ExampleGet(t *testing.T) {
	usr, _ := user.Current()
	downloader := get.NewDownloader(filepath.Join(usr.HomeDir, ".cache", "go-get"))
	path, err := downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Downloaded to %s\n", path)
}
