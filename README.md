# Example of Grafana Alerting in HA

Documentation: https://grafana.com/docs/grafana/latest/alerting/set-up/configure-high-availability/

Based on Grafana examples:

- https://github.com/grafana/alerting-ha-docker-examples/tree/main
- https://github.com/grafana/demo-prometheus-and-grafana-alerts
- https://github.com/grafana/provisioning-alerting-examples


And:

- https://github.com/VahagnMian/golang-microservices-observability

## 🛠 Pre-requisites

Before diving in, ensure you have the following:

1. **`psql`**: Not installed? 🛑 Installation both for Mac and Ubuntu Follow the guide [here](https://www.timescale.com/blog/how-to-install-psql-on-mac-ubuntu-debian-windows/).

2.  Execute the following commands to set up your database:
```bash
docker-compose up -d postgres
psql -h localhost -p 5432 -U postgres -c "CREATE DATABASE users;"
psql -h localhost -p 5432 -U postgres -d users -a -f ./config/resources/users.sql
```


# Features

## Run Grafana in HA

## OTEL collector for ingest metrics and traces

## Send logs to Loki

## Correlate traces with logs and metrics

## Send profiles to Pyroscope

## Grafana alerting

## Grafana - Prometheus alerting integration

