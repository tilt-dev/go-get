package get

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tmpdir(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err, "Could not create tmpdir")
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return dir
}

func setupDir(t *testing.T) string {
	t.Helper()

	// force the tests to run in an empty working directory
	// because otherwise the working directory will be the `go-get` local
	// repo path, which is itself a git command, so commands that don't
	// properly chdir could inadvertently run against this repo
	fakeWd := tmpdir(t)
	cwd, err := os.Getwd()
	require.NoError(t, err, "Could not get cwd")
	require.NoError(t, os.Chdir(fakeWd), "Could not change cwd")
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Errorf("failed to restore current working directory: %v", err)
		}
	})

	return tmpdir(t)
}

func TestGet(t *testing.T) {
	dir := setupDir(t)
	downloader := NewDownloader(dir)
	path, err := downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	tiltfile, err := ioutil.ReadFile(filepath.Join(path, "Tiltfile"))
	require.NoError(t, err)
	assert.Contains(t, string(tiltfile), `print("Hello world!")`)

	path, err = downloader.Download("gitlab.com/nicks6/tilt-extension-experiment")
	require.NoError(t, err)

	tiltfile, err = ioutil.ReadFile(filepath.Join(path, "tilt_extension", "Tiltfile"))
	require.NoError(t, err)
	assert.Contains(t, string(tiltfile), `def dotnet_build(ref,context = ".", version = "latest"):`)
	assert.Contains(t, string(tiltfile), `def vue_build(ref,context = ".", version = "latest"):`)
}

func TestDestinationPath(t *testing.T) {
	dir := setupDir(t)
	downloader := NewDownloader(dir)
	result := downloader.DestinationPath("github.com/tilt-dev/tilt-extensions/hello_world")
	assert.True(t, strings.HasSuffix(result, filepath.Join("github.com", "tilt-dev", "tilt-extensions", "hello_world")))

	result = downloader.DestinationPath("gitlab.com/nicks6/tilt-extension-experiment")
	assert.True(t, strings.HasSuffix(result, filepath.Join("gitlab.com", "nicks6", "tilt-extension-experiment")))
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

	result, err = dlr.Download("gitlab.com/nicks6/tilt-extension-experiment")
	require.NoError(t, err)

	err = dlr.RefSync("gitlab.com/nicks6/tilt-extension-experiment",
		"93df1cfd37237550d0e10b2384a45a889200c1ef")
	require.NoError(t, err)

	contents, err = ioutil.ReadFile(filepath.Join(result, "README.md"))
	require.NoError(t, err)

	assert.Contains(t, string(contents), "# tilt-extension-experiment")
}

func TestVersionTagSync(t *testing.T) {
	dir := setupDir(t)
	dlr := NewDownloader(dir)
	result, err := dlr.Download("github.com/tilt-dev/tilt-extensions")
	require.NoError(t, err)

	err = dlr.RefSync("github.com/tilt-dev/tilt-extensions", "v0.25.0")
	require.NoError(t, err)

	contents, err := ioutil.ReadFile(filepath.Join(result, "README.md"))
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(string(contents), "# Tilt Extensions"))

	result, err = dlr.Download("gitlab.com/brshore/tilt-extension-experiment")
	require.NoError(t, err)

	err = dlr.RefSync("gitlab.com/brshore/tilt-extension-experiment", "v0.0.1")
	require.NoError(t, err)

	contents, err = ioutil.ReadFile(filepath.Join(result, "README.md"))
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(string(contents), "# tilt-extension-experiment"))
}

func TestHeadRef(t *testing.T) {
	dir := setupDir(t)
	downloader := NewDownloader(dir)
	_, err := downloader.HeadRef("github.com/tilt-dev/tilt-extensions/hello_world")
	if assert.Error(t, err) {
		var pathErr *os.PathError // use *os.PathError type, error string differs between OS
		assert.True(t, errors.As(err, &pathErr))
	}

	_, err = downloader.Download("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)

	ref, err := downloader.HeadRef("github.com/tilt-dev/tilt-extensions/hello_world")
	require.NoError(t, err)
	assert.True(t, regexp.MustCompile("^[0-9a-f]{40}$").MatchString(ref))

	_, err = downloader.Download("gitlab.com/nicks6/tilt-extension-experiment")
	require.NoError(t, err)

	ref, err = downloader.HeadRef("gitlab.com/nicks6/tilt-extension-experiment")
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

	path, err = downloader.Download("gitlab.com/nicks6/tilt-extension-experiment")
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(path, "tilt_extension", "Tiltfile"), []byte(`print("Goodbye world!")`), os.FileMode(0644))
	require.NoError(t, err)

	path2, err = downloader2.Download("gitlab.com/nicks6/tilt-extension-experiment")
	require.NoError(t, err)

	tiltfile, err = ioutil.ReadFile(filepath.Join(path2, "tilt_extension", "Tiltfile"))
	require.NoError(t, err)
	assert.Contains(t, string(tiltfile), `print("Goodbye world!")`)
}
