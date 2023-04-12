package dockerapi

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/fsnotify/fsnotify"
)

const (
	dockerRequestTimeout = 3 * time.Second
)

type DockerApi struct {
	watcher    *fsnotify.Watcher
	Containers map[string]string
	Errors     chan error
}

func (api *DockerApi) Close() {
	api = nil
}

func New() (*DockerApi, error) {
	watcher, err := fsnotify.NewWatcher()
	app := DockerApi{
		watcher:    watcher,
		Containers: make(map[string]string, 0),
	}
	return &app, err
}

func (api *DockerApi) RestartWatcher() (err error) {
	api.watcher.Close()
	api.watcher, err = fsnotify.NewWatcher()
	for _, dir := range api.Containers {
		err := api.watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	// Start a loop to monitor for events
	for {
		select {
		case event, ok := <-api.watcher.Events:
			if !ok {
				return
			}
			// Print a message to the console when any changes are detected
			log.Println("event:", event.Name)

			// Restart the container when any changes are detected
			for container, dir := range api.Containers {
				if !strings.Contains(event.Name, dir) {
					continue
				}
				err = restartContainer(dockerClient, container)
				if err != nil {
					log.Println(err)
				}

			}
		case err, ok := <-api.watcher.Errors:
			if !ok {
				return err
			}
			log.Println("error:", err)
		}
	}
}

type restartOptions struct {
	signal         string
	timeout        int
	timeoutChanged bool

	containers []string
}

// Restart the specified Docker container
func restartContainer(dockerClient *client.Client, containerName string) error {
	// Get the container ID by name
	args := filters.Arg("name", containerName)
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(args)})
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		return nil
	}
	//containerID := containers[0].ID
	restartTarget := restartOptions{
		containers: []string{containerName},
	}

	// Restart the container
	err = RunRestart(dockerClient, &restartTarget)
	if err != nil {
		return err
	}

	return nil
}

func RunRestart(dockerCli *client.Client, opts *restartOptions) error {

	var errs []string

	for _, name := range opts.containers {

		err := dockerCli.ContainerRestart(context.Background(), name, container.StopOptions{})

		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
