package docker

import (
	"context"
	"encoding/json"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) ListContainers(ctx context.Context) ([]container.Summary, error) {
	result, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) Inspect(ctx context.Context, id string) (container.InspectResponse, error) {
	result, err := c.cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return container.InspectResponse{}, err
	}
	return result.Container, nil
}

func (c *Client) Stats(ctx context.Context, id string) (*container.StatsResponse, error) {
	result, err := c.cli.ContainerStats(ctx, id, client.ContainerStatsOptions{Stream: false})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	var stats container.StatsResponse
	if err := json.NewDecoder(result.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (c *Client) Stop(ctx context.Context, id string) error {
	_, err := c.cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	return err
}

func (c *Client) Start(ctx context.Context, id string) error {
	_, err := c.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
	return err
}

func (c *Client) Restart(ctx context.Context, id string) error {
	_, err := c.cli.ContainerRestart(ctx, id, client.ContainerRestartOptions{})
	return err
}

func (c *Client) Remove(ctx context.Context, id string) error {
	_, err := c.cli.ContainerRemove(ctx, id, client.ContainerRemoveOptions{})
	return err
}