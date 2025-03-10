import { check } from "k6";
import remote from "k6/x/remotewrite";

export let options = {
  iterations: 5000
};

const client = new remote.Client({
  url: "http://localhost:9090/api/v1/write",
});

// avg_over_time(cpu_usage[5m])
// exclude region
// - avg(avg_over_time(cpu_usage[5m])) by (instance)
// - avg(avg_over_time(cpu_usage[5m]) ) without (region)

export default function () {
    sendMetricData("server1", Math.floor(Math.random() * 21) + 80); // 80-100
    sendMetricData("server2", Math.floor(Math.random() * 21) + 0); // 0-20
    sendMetricData("server3", Math.floor(Math.random() * 21) + 80); // 80-100
}

function sendMetricData(instanceValue, value) {
  const regions = ["emea", "amer", "apac"];
  const res = client.store([
    {
      labels: [
        { name: "__name__", value: `cpu_usage` },
        { name: "job", value: "exporter" },
        { name: "region", value: regions[Math.floor(Math.random() * regions.length)] },
        { name: "instance", value: instanceValue },
      ],
      samples: [{ value: value }],
    },
  ]);

  check(res, {
    "is status 204": (r) => r.status === 204,
  });
}
