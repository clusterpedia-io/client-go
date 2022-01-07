clusterpedia-client supports the use of native `client-go` mode to call the [clusterpedia](https://github.com/clusterpedia-io/clusterpedia) API. clisterpedia uses labelSelector to query data. clusterpedia- client provides a more convenient way to build options for various complex conditions, which is more suitable for users.

### usage

In multiple clusters, build options in a chained stype.

```golang
options := clusterpedia.ListOptionsBuild().Clusters("cluster-01").Namespaces("kube-system").Offset(10).Size(5).OrderBy(Order{"dsad", false}).Options()
```

You can get the `clientset` of client-go connect to clusterpedia.

### example

Here are some [examples](./examples) where clusterpedia-client can be used more easily.