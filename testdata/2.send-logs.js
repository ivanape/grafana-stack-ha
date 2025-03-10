import { check } from "k6";
import loki from "k6/x/loki";

let labels = loki.Labels({
  format: ["logfmt"],
  // detected_level: ["warn", "error"],
  // instance: ["foo", "bar"],
  detected_level: ["error"],
  instance: ["foo"],
  service_name: ["backend"],
});

// Example query:
// count_over_time({detected_level="error", service_name="backend"}[1m])

const conf = new loki.Config("http://localhost:3100", 10000, 1.0, {}, labels);
const client = new loki.Client(conf);

export default () => {
  // push random data (~800-900 log lines) according to the labels
  const res = client.push();
  check(res, { "successful write": (res) => res.status == 204 });
};
