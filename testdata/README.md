# Run the tests

To create a k6 build to run the test simple execute:

```make build_k6``

To generate service metrics, logs, traces and profiling run:

```./k6 run --iterations 100 testdata/synthetic-traffic.js```

Use the rest of scripts to manage alerts.