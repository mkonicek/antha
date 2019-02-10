package dockerexec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrimer3(t *testing.T) {

	t.Skip("skipping test")

	boulder_in := `SEQUENCE_ID=example
SEQUENCE_TEMPLATE=TATTGGTGAAGCCTCAGGTAGTGCAGAATATGAAACTTCAGGATCCAGTGGGCATGCTACTGGTAGTGCTGCCGGCCTTACAGGCATTATGGTGGCAAAGTCGACAGAGTTTA
PRIMER_TASK=generic
PRIMER_PICK_LEFT_PRIMER=1
PRIMER_PICK_INTERNAL_OLIGO=0
PRIMER_PICK_RIGHT_PRIMER=1
PRIMER_OPT_SIZE=20
PRIMER_MIN_SIZE=18
PRIMER_MAX_SIZE=22
PRIMER_PRODUCT_SIZE_RANGE=75-150
PRIMER_EXPLAIN_FLAG=1
=`
	inputFileName := "input.boulder"
	resultsFileName := "output.formatted"

	// Create temporary directory to mount in container for file IO
	tmpDir := os.TempDir()
	// !! not shared from OSX and not known to Docker... mounts denied...
	// use cwd for now TODO: resolve
	tmpDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	hostDir, err := ioutil.TempDir(tmpDir, "primer3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(hostDir)

	// Write input to file
	if err := ioutil.WriteFile(filepath.Join(hostDir, inputFileName), []byte(boulder_in), 0666); err != nil {
		t.Fatal(err)
	}

	// Docker image
	image := "docker.io/dcouc01/primer3:latest"

	// Container volume in which to mount hostDir
	containerDir := "/mnt"
	volumes := []string{containerDir}

	// Correspondence(s)
	binds := []string{strings.Join([]string{hostDir, containerDir}, ":")}

	// Blast command
	command := []string{
		"/home/biouser/primer3/src/primer3_core",
		filepath.Join(containerDir, inputFileName),
		"-output", filepath.Join(containerDir, resultsFileName),
		"--format_output",
	}

	// Create the container, run the command
	err = runDocker(image, command, volumes, binds)
	if err != nil {
		t.Fatal(err)
	}

	// Read the output
	content, err := ioutil.ReadFile(filepath.Join(hostDir, resultsFileName))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)

}
