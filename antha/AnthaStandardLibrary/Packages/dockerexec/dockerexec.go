package dockerexec

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"os"
)

// https://docs.docker.com/develop/sdk/
// https://docs.docker.com/engine/api/v1.39/
// https://godoc.org/github.com/docker/docker/api

const defaultDockerAPIVersion = "v1.39"

func runDocker(image string, command, volumes, binds []string) (err error) {

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		return
	}

	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return
	}
	io.Copy(os.Stdout, reader)

	mapped := make(map[string]struct{})
	for _, volume := range volumes {
		mapped[volume] = struct{}{}
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:   image,
			Cmd:     command,
			Volumes: mapped,
			Tty:     true,
		},
		&container.HostConfig{
			Binds: binds,
		},
		nil,
		"",
	)
	if err != nil {
		return
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return
	}

	io.Copy(os.Stdout, out)

	return
}
