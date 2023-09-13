package collector

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
)

// DockerContainers implements the Collector interface.
type DockerContainers struct {
	Label *string
	*client.Client
}

var _ prometheus.Collector = (*DockerContainers)(nil)

var (
	containerUptime = prometheus.NewDesc(
		"docker_container_state",
		"Container collected from Docker PS",
		[]string{"container_name", "container_id", "state"}, nil,
	)
)

// Describe provides metric decriptiors.
func (c DockerContainers) Describe(ch chan<- *prometheus.Desc) {
	ch <- containerUptime
}

// Collect scrapes the container information from Docker.
func (c DockerContainers) Collect(ch chan<- prometheus.Metric) {
	containers, err := c.Client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		log.Printf("Error while fetching container list: %s", err)
		return
	}

	for _, container := range containers {
		if container.Labels[*c.Label] == "true" {
			for _, name := range container.Names {
				up := isContainerUp(container)
				ch <- prometheus.MustNewConstMetric(
					containerUptime,
					prometheus.GaugeValue,
					boolToGaugeValue(up),
					name,
					container.ID,
					container.State,
				)
			}
		}
	}
}

func isContainerUp(container types.Container) bool {
	return strings.EqualFold(container.State, "running")
}

func boolToGaugeValue(val bool) float64 {
	if val {
		return 1
	}
	return 0
}
