package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	defaults "github.com/Synthace/microservice/cmd/defaults/protobuf"
	element "github.com/Synthace/microservice/cmd/element/protobuf"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// An HTTPCall is an RPC over HTTP
type HTTPCall struct {
	Endpoint string
	Path     string
	Token    string
	Request  proto.Message
	Response proto.Message
}

// Call executes the HTTPCall
func (c *HTTPCall) Call() error {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return err
	}
	u.Path = c.Path

	var out bytes.Buffer
	var m jsonpb.Marshaler
	if err := m.Marshal(&out, c.Request); err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", u.String(), &out)
	if err != nil {
		return err
	}

	httpReq.Header.Add("Authorization", "bearer "+c.Token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bs, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("invalid status: %s: %s", resp.Status, string(bs))
	}

	if c.Response != nil {
		bs, _ := ioutil.ReadAll(resp.Body)
		um := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err := um.Unmarshal(bytes.NewReader(bs), c.Response); err != nil {
			return err
		}
	}
	return nil
}

func run() error {
	var token string
	var elementPB string
	var metadataJSON string
	var endpoint string
	var elementSet string
	var protocol string

	flag.StringVar(&token, "token", "", "Authorization token")
	flag.StringVar(&elementPB, "element-proto", "/elements.pb", "Elements to upload")
	flag.StringVar(&metadataJSON, "metadata-json", "/metadata.json", "Elements metadata to upload")
	flag.StringVar(&elementSet, "element-set", "", "Element set name")
	flag.StringVar(&protocol, "protocol", "http", "Http")
	flag.StringVar(&endpoint, "endpoint", "", "Endpoint")
	flag.Parse()

	if protocol == "grpc" {
		return errors.New("protocol not supported")
	}

	{
		bs, err := ioutil.ReadFile(elementPB)
		if err != nil {
			return err
		}

		var req element.CreateElementsRequest
		if err := proto.Unmarshal(bs, &req); err != nil {
			return err
		}

		req.SetName = elementSet
		req.RequiresAnthaCore = true

		call := &HTTPCall{
			Endpoint: endpoint,
			Token:    token,
			Request:  &req,
			Response: &element.Empty{},
			Path:     "/api/pub-v1/element/create-elements",
		}
		if err := call.Call(); err != nil {
			return err
		}
		fmt.Printf("Created element set\n")
	}

	{
		bs, err := ioutil.ReadFile(metadataJSON)
		if err != nil {
			return err
		}

		resp := &defaults.SetResponse{}
		call := &HTTPCall{
			Endpoint: endpoint,
			Token:    token,
			Request: &defaults.ScopedSetRequest{
				ElementSetName: elementSet,
				Payload: &defaults.WorkflowDefaults{
					Rawdoc: bs,
				},
			},
			Response: resp,
			Path:     "/api/pub-v1/defaults/workflow/set",
		}

		if err := call.Call(); err != nil {
			return err
		}
		fmt.Printf("Uploaded defaults %s:%s\n", resp.Op, resp.Id)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
