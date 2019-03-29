// Package jobfile provides for basic operations for manipulating files
// associated with a job
package jobfile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	webdav "github.com/studio-b12/gowebdav"
)

// JobID is the ID of a job run through Antha.
// This can be used to obtain the outputs of that Job.
type JobID string

var (
	errDoesNotExist = errors.New("file does not exist")
)

// A Client manages a connection with the job file service
type Client struct {
	username     string
	password     string
	jobID        string
	webdavClient *webdav.Client
	apiClient    *apiClient
}

// DefaultClient creates a Client based on parameters available in the
// environment
func DefaultClient() (*Client, error) {
	username := os.Getenv("ANTHA_USERNAME")
	password := os.Getenv("ANTHA_PASSWORD")
	webdavEndpoint := os.Getenv("WEBDAV_ENDPOINT")
	anthaEndpoint := os.Getenv("ANTHA_ENDPOINT")
	jobID := os.Getenv("METADATA_JOB_ID")
	ticketID := os.Getenv("ANTHA_UPLOAD_TICKET_ID")

	return &Client{
		username: username,
		password: password,
		jobID:    jobID,

		webdavClient: webdav.NewClient(webdavEndpoint, username, password),

		apiClient: &apiClient{
			c:        http.DefaultClient,
			endpoint: anthaEndpoint,
			ticketID: ticketID,
			username: username,
			password: password,
		},
	}, nil
}

type errWrapper struct {
	Reader io.Reader
	Writer io.Writer
	Closer io.Closer
	Err    error
}

func (w *errWrapper) Write(p []byte) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	return w.Writer.Write(p)
}

func (w *errWrapper) Read(p []byte) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	return w.Reader.Read(p)
}

func (w *errWrapper) Close() error {
	if w.Err != nil {
		return w.Err
	}
	return w.Closer.Close()
}

// A File is data about a file
type File struct {
	Name  string
	Dir   string
	Size  int64
	IsDir bool
}

// dirsToFiles will walk through all sub directories of files and append the list of
// files returned with files found in that sub directory and sub directories of that.
func dirsToFiles(c *Client, files []*File) ([]*File, error) {
	for _, f := range files {
		if f.IsDir {
			subDir := path.Join(f.Dir, f.Name)
			fis, err := c.webdavClient.ReadDir(subDir)
			if err != nil {
				return nil, err
			}
			var subDirFiles []*File
			for _, fi := range fis {
				subDirFiles = append(subDirFiles, &File{
					Name:  fi.Name(),
					Dir:   subDir,
					IsDir: fi.IsDir(),
					Size:  fi.Size(),
				})
			}
			// run recursively to expand any sub directories
			subDirFiles, err = dirsToFiles(c, subDirFiles)
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		}
	}
	return files, nil
}

// ListFiles returns files for a job. If jobID is empty, list files for the
// current job.
func (c *Client) ListFiles(ctx context.Context, jobID JobID) ([]*File, error) {
	if len(jobID) == 0 {
		jobID = JobID(c.jobID)
	}

	dir := fmt.Sprintf("/job/%s", jobID)
	fis, err := c.webdavClient.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, fi := range fis {
		files = append(files, &File{
			Name:  fi.Name(),
			Dir:   dir,
			IsDir: fi.IsDir(),
			Size:  fi.Size(),
		})
	}
	return dirsToFiles(c, files)
}

// NewWriter returns a writer to the filename in the current job
func (c *Client) NewWriter(ctx context.Context, name string) io.WriteCloser {
	writer, err := c.apiClient.WriteStream(ctx, name)
	return &errWrapper{
		Writer: writer,
		Closer: writer,
		Err:    err,
	}
}

// NewReader returns a reader of a filename in the given job. If
// the job id is missing, read a file in the current job.
func (c *Client) NewReader(ctx context.Context, jobID JobID, name string) io.ReadCloser {
	if len(jobID) == 0 {
		jobID = JobID(c.jobID)
	}

	p := fmt.Sprintf("/job/%s/%s", jobID, name)

	// Stat first because ReadStream will blindly return the webDAV response,
	// e.g., "Unauthorized\n" for bad paths
	fi, err := c.webdavClient.Stat(p)
	if err == nil && fi.IsDir() {
		// WebDAV sometimes returns a Dir response for missing files.
		err = errDoesNotExist
	}
	if err != nil {
		return &errWrapper{
			Err: err,
		}
	}

	reader, err := c.webdavClient.ReadStream(p)

	return &errWrapper{
		Reader: reader,
		Closer: reader,
		Err:    err,
	}
}
