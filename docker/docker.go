package docker

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Pull(imageName string) error {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Set the context for the API request
	ctx := context.Background()

	// Define the options for the pull operation
	options := types.ImagePullOptions{}

	// Pull the Docker image
	resp, err := cli.ImagePull(ctx, imageName, options)
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %v", imageName, err)
	}
	defer resp.Close()

	fmt.Println("Image pulled successfully!")
	return nil
}

func PullCmd(imageName string) error {
	start := time.Now()
	cmd := exec.Command("docker", "pull", imageName)

	err := cmd.Run()
	if err != nil {
		return err
	}

	t := time.Now()
	elapsedTime := t.Sub(start)
	fmt.Printf("Docker pull took %s\n", elapsedTime)

	return nil
}

func Tag(oldTag string, newTag string) error {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Set the context for the API request
	ctx := context.Background()

	err = cli.ImageTag(ctx, oldTag, newTag)
	if err != nil {
		return err
	}
	return nil
}

func RemoveTag(tag string) error {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Set the context for the API request
	ctx := context.Background()

	// Define the options for the pull operation
	options := types.ImageRemoveOptions{}

	_, err = cli.ImageRemove(ctx, tag, options)
	if err != nil {
		return err
	}
	return nil
}
