package images

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/paddycarey/ims/pkg/storage"
)

func Optimize(f storage.File) (storage.File, error) {

	switch f.MimeType() {
	case "image/gif":
		return optimizeGIF(f)
	case "image/jpeg":
		return optimizeJPEG(f)
	case "image/png":
		return optimizePNG(f)
	}

	return nil, errors.New("No optimizer found for format")
}

func runCommand(stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}

	return cmd.Run()
}

func checkForOptimizer(o string) bool {
	if err := runCommand(nil, nil, nil, "which", o); err != nil {
		return false
	}
	return true
}

func optimizeGIF(f storage.File) (storage.File, error) {

	if !checkForOptimizer("gifsicle") {
		logrus.Debug("gifsicle not installed, skipping optimisation")
		return f, nil
	}

	// run gifsicle on the image to optimise it
	var out bytes.Buffer
	if err := runCommand(f, &out, nil, "gifsicle", "-O3"); err != nil {
		return nil, err
	}

	f = &storage.InMemoryFile{bytes.NewReader(out.Bytes()), f.MimeType(), time.Now()}
	return f, nil
}

func optimizeJPEG(f storage.File) (storage.File, error) {

	if !checkForOptimizer("jpegtran") {
		logrus.Debug("jpegtran not installed, skipping optimisation")
		return f, nil
	}

	// run jpegtran on the image to optimise it
	var out bytes.Buffer
	if err := runCommand(f, &out, nil, "jpegtran", "-optimize", "-progressive"); err != nil {
		return nil, err
	}

	f = &storage.InMemoryFile{bytes.NewReader(out.Bytes()), f.MimeType(), time.Now()}
	return f, nil
}

func optimizePNG(f storage.File) (storage.File, error) {

	if !checkForOptimizer("optipng") {
		logrus.Debug("optipng not installed, skipping optimisation")
		return f, nil
	}

	// create a temporary file we can write to before renaming the newly
	// written file into its final place
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	// write data to temporary file, explicitly syncing and closing it
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	tmpFile.Write(b)
	if err := tmpFile.Sync(); err != nil {
		return nil, err
	}
	tmpFile.Close()

	// run optipng on the image to optimise it
	var out bytes.Buffer
	if err := runCommand(nil, &out, &out, "optipng", "-force", "-fix", "-o3", tmpFile.Name()); err != nil {
		return nil, err
	}

	// read the newly optimised file back into memory
	tf, err := os.Open(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(tf)
	if err != nil {
		return nil, err
	}

	f = &storage.InMemoryFile{bytes.NewReader(b), f.MimeType(), time.Now()}
	return f, nil
}
