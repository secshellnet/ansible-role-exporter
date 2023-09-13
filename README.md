# ansible-role-exporter

This ansible role can be used to install prometheus exporters on the system.  
[node exporter](https://github.com/prometheus/node_exporter) will be installed always.

If the variable `install_docker` is set to true [cadvisor](https://github.com/google/cadvisor) and [docker-exporter](./files/docker_status_exporter_src/) will be installed too.
