package docker

import (
	"context"
	"log"

	"git.cryptic.systems/volker.raschek/dockerutils"
)

type Runtime struct {
	dockerClient *dockerutils.Client
}

func (r *Runtime) GetImageLabels(ctx context.Context, image string) (map[string]string, error) {
	log.Printf("Pull image: %v", image)
	if err := r.dockerClient.PullQuiet(ctx, image); err != nil {
		return nil, err
	}

	log.Printf("Image successfully pulled: %v", image)

	imageSpec, _, err := r.dockerClient.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return nil, err
	}

	return imageSpec.Config.Labels, nil
}

func NewRuntime(dockerClient *dockerutils.Client) (*Runtime, error) {
	return &Runtime{
		dockerClient: dockerClient,
	}, nil
}
