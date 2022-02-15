# onix

Releases comparison tool. Use same metrics for 2 different releases and visualize it.

### Stack

* Go 1.16
* Postgres 10
* React

### Components

* api dashboard-admin
* api dashboard-main
* api system
* daemon metrics-extractor
* stub prometheus
* util load-historical-metrics

### Definitions

* **Service** - single process or group of processes (in the same namespace in k8s, etc).
* **Source** - time series database (Prometheus/Thanos or InfluxDB).
* **Release** - info about new deployment of some service.
* **Criteria** - prometheus/influx query to extract pairs \<timestamp, double-point-value>
* **Measurement** - locally cached metric, object with source_id, criteria_id, timestamp and value.

### Demo

```shell
docker-compose -f docker-compose.demo.yml down --volumes --remove-orphans
docker-compose -f docker-compose.demo.yml up -d
echo "Open http://localhost:23001/"
```

<kbd>
    <img src="docs/demo-screen.png" alt="Demo screen"/>
</kbd>