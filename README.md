# What is this
This is a custom Prometheus metrics for monitoring cpu utilization per function in OpenFaaS. It uses mongodb to find total cpu usage per function that is created using [custom template](https://github.com/borodun/of-templates)

`mongo_uri` environmental variable is required. See [mongo uri string](https://docs.mongodb.com/manual/reference/connection-string/)