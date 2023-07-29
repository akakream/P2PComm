package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Pull(imageName string) error {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
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

	// Copy the response output to stdout for progress information (optional)
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return fmt.Errorf("failed to copy response: %v", err)
	}

	fmt.Println("Image pulled successfully!")
	return nil
}
