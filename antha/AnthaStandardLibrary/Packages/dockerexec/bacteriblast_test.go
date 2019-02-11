package dockerexec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBlast(t *testing.T) {

	t.Skip("skipping test")

	query := `
>AJF59511
GTGACCAAACAGGAAAAAACCGCCCTGAACATGGCCCGCTTCATCAGAAGCCAGACATTA
ACCCTGCTGGAGAAGCTCAACGAACTGGACGCGGATGAACAGGCAGACATCTGTGAATCG
CTTCACGACCACGCCGATGAGCTTTACCGCAGTTGCCTCGCACGTTTCGGGGATGACGGT
GAAACCCTCTGA
`
	queryFileName := "query.fa"
	resultsFileName := "blast.out"

	// Create temporary directory to mount in container for file IO
	tmpDir := os.TempDir()
	// !! not shared from OSX and not known to Docker... mounts denied...
	// use cwd for now TODO: resolve
	tmpDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	hostDir, err := ioutil.TempDir(tmpDir, "blast")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(hostDir)

	// Write query to file
	if err := ioutil.WriteFile(filepath.Join(hostDir, queryFileName), []byte(query), 0666); err != nil {
		t.Fatal(err)
	}

	// Docker image
	image := "docker.io/dcouc01/bacteriblast:latest"

	// Container volume in which to mount hostDir
	containerDir := "/mnt"
	volumes := []string{containerDir}

	// Correspondence(s)
	binds := []string{strings.Join([]string{hostDir, containerDir}, ":")}

	// Blast command
	command := []string{
		"blastn",
		"-query", filepath.Join(containerDir, queryFileName),
		"-db", "Escherichia_coli_1303.ASM82998v1.cdna.all",
		"-out", filepath.Join(containerDir, resultsFileName),
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
