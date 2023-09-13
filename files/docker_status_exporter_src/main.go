package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/docker/docker/client"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/secshellnet/docker_status_exporter/collector"
	"net/http"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics.").Default(":9400").String()
	metricsPath   = kingpin.Flag("web.metric-path", "Path under which to expose metrics.").Default("/metrics").String()
	labelFlag     = kingpin.Flag("label.name", "Label to Flag Container for Monitoring.").Default("monitoring").String()
)

func boot() error {
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting node_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	coll := collector.DockerContainers{Label: labelFlag, Client: dockerClient}
	if err := prometheus.Register(&coll); err != nil {
		return err
	}

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Docker Container Exporter</title></head>
			<body>
			<h1>Docker Container Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	return http.ListenAndServe(*listenAddress, nil)
}

func main() {
	kingpin.Parse()
	kingpin.FatalIfError(boot(), "")
}
