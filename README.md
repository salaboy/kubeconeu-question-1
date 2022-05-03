# Go HTTP Function

Welcome to your new Go Function! The boilerplate function code can be found in
[`handle.go`](handle.go). This Function responds to HTTP requests.

This function connects to Redis using two environment variables, `REDIS_HOST`
and `REDIS_PASSWORD`. Use a Kubernetes `Secret` to store these values in the
cluster.

kubectl create secret generic redis --from-literal=redis_host=10.200.130.188:6379 --from-literal=redis_password=1234

## Development

Develop new features by adding a test to [`handle_test.go`](handle_test.go) for
each feature, and confirm it works with `go test`.

Update the running analog of the function using the `func` CLI or client
library, and it can be invoked from your browser or from the command line:

```console
curl http://myfunction.example.com/
```

For more, see [the complete documentation]('https://github.com/knative-sandbox/kn-plugin-func/tree/main/docs')


