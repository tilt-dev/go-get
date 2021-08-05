package get

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return dir
}

func TestGet(t *testing.T) {
	dir := setupDir(t)
	downloader := NewDownloader(dir)
	path, err := downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	tiltfile, err := ioutil.ReadFile(filepath.Join(path, "Tiltfile"))
	require.NoError(t, err)
	assert.Contains(t, string(tiltfile), `print("Hello world!")`)
}

func TestDestinationPath(t *testing.T) {
	dir := setupDir(t)
	downloader := NewDownloader(dir)
	result := downloader.DestinationPath("github.com/tilt-dev/tilt-extensions/hello_world")
	assert.True(t, strings.HasSuffix(result, "/github.com/tilt-dev/tilt-extensions/hello_world"))
}

func TestTagSync(t *testing.T) {
	dir := setupDir(t)
	dlr := NewDownloader(dir)
	result, err := dlr.Download("github.com/tilt-dev/tilt-extensions")
	require.NoError(t, err)

	err = dlr.RefSync("github.com/tilt-dev/tilt-extensions",
		"5670e8d8c5925981d019e6065a30058fc50eb299") // initial commit
	require.NoError(t, err)

	contents, err := ioutil.ReadFile(filepath.Join(result, "README.md"))
	require.NoError(t, err)

	assert.Equal(t, "# tilt-extensions", string(contents))
}

func TestHeadRef(t *testing.T) {
	dir := setupDir(t)
	downloader := NewDownloader(dir)
	_, err := downloader.HeadRef("github.com/tilt-dev/tilt-extensions/hello_world")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "no such file or directory")
	}

	_, err = downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	ref, err := downloader.HeadRef("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)
	assert.True(t, regexp.MustCompile("^[0-9a-f]{40}$").MatchString(ref))
}

func TestGetCached(t *testing.T) {
	dir := setupDir(t)
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
