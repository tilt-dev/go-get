package get

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	dir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	downloader := NewDownloader(dir)
	path, err := downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	tiltfile, err := ioutil.ReadFile(filepath.Join(path, "Tiltfile"))
	require.NoError(t, err)
	assert.Contains(t, string(tiltfile), `print("Hello world!")`)
}

func TestGetCached(t *testing.T) {
	dir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	downloader := NewDownloader(dir)
	path, err := downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(path, "Tiltfile"), []byte(`print("Goodbye world!")`), os.FileMode(0644))
	require.NoError(t, err)

	// Make sure that running a second downloader doesn't overwrite.
	//
	// TODO(nick): Do we need a mechanism to force an update? Or should that be up
	// to the app to delete the cached download?
	downloader2 := NewDownloader(dir)
	path2, err := downloader2.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	tiltfile, err := ioutil.ReadFile(filepath.Join(path2, "Tiltfile"))
	require.NoError(t, err)
	assert.Contains(t, string(tiltfile), `print("Goodbye world!")`)
}
