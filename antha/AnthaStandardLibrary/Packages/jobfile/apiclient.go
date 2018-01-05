package jobfile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	regenerateTicketPath = "/api/pub-v1/datarepo/regenerate-upload-ticket"
	jsonContentType      = "application/json"
)

var (
	errNoMatchingURL = errors.New("no matching url")
)

type apiClient struct {
	c        *http.Client
	endpoint string
	ticketID string
	username string
	password string
}

type regenerateRequest struct {
	FileNames []string `json:"file_names"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Resumable bool     `json:"resumable"`
	TicketID  string   `json:"ticket_id"`
}

type regenerateResponse struct {
	ID        string   `json:"id"`
	SignedURL []string `json:"signed_url"`
}

func (c *apiClient) regenerateTicket(ctx context.Context, req *regenerateRequest) (*regenerateResponse, error) {
	u, err := url.Parse(c.endpoint + regenerateTicketPath)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	postResp, err := c.c.Post(u.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer postResp.Body.Close() // nolint

	if postResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error code %d", postResp.StatusCode)
	}

	bs, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return nil, err
	}

	var resp regenerateResponse
	if err := json.Unmarshal(bs, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func matchURL(signedURLs []string, name string) (string, error) {
	for _, s := range signedURLs {
		p, err := url.Parse(s)
		if err != nil {
			continue
		}
		n := path.Base(p.Path)
		if n == name {
			return s, nil
		}
	}

	return "", errNoMatchingURL
}

func (c *apiClient) WriteStream(ctx context.Context, name string) (io.WriteCloser, error) {
	resp, err := c.regenerateTicket(ctx, &regenerateRequest{
		FileNames: []string{name},
		Username:  c.username,
		Password:  c.password,
		TicketID:  c.ticketID,
	})
	if err != nil {
		return nil, err
	}

	signedURL, err := matchURL(resp.SignedURL, name)
	if err != nil {
		return nil, err
	}

	reader, writer := io.Pipe()

	go func() {
		defer reader.Close() // nolint

		req, err := http.NewRequest("PUT", signedURL, reader)
		if err != nil {
			reader.CloseWithError(err) // nolint
			return
		}

		resp, err := c.c.Do(req)
		if err != nil {
			reader.CloseWithError(err) // nolint
			return
		}
		defer resp.Body.Close() // nolint

		if resp.StatusCode != http.StatusOK {
			reader.CloseWithError(fmt.Errorf("error code %d", resp.StatusCode)) // nolint
			return
		}
	}()

	return writer, nil
}
