# What is this
This is a custom Prometheus metrics for monitoring cpu utilization per function in OpenFaaS. It uses mongodb to find total cpu usage per function that is created using [custom template](https://github.com/borodun/of-templates)

`mongo_uri` environmental variable is required. See [mongo uri string](https://docs.mongodb.com/manual/reference/connection-string/)

You need to deploy this metrics in kubernetes cluster using config in **k8s** folder and add a **job** in Prometheus config([example](https://github.com/borodun/of-tool/blob/13818fada2cdbb46ea2f1c3b8b1ff31fb3bf3dde/k8s/Prometheus/config-map.yaml#L188))