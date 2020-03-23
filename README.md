[![CircleCI](https://circleci.com/gh/angular/angular-bazel-example.svg?style=svg)](https://circleci.com/gh/angular/angular-bazel-example)

# Example Bazel + Angular+ Golang + Protocol buffers + gRPC stack

This doesn't work yet.

## Clipboard farm

**All of these commands should be run from the viz directory.**

Some commands:

```shell
bazel run @nodejs//:yarn -- add plotly.js
```

```shell
bazel run @nodejs//:yarn -- add angular-plotly.js
```

```shell
bazel run //:gazelle
```

Example of adding a new repo

```shell
bazel run //:gazelle update-repos google.com/protobuf
```

```shell
bazel aquery @npm//google-protobuf:all --output=textproto
```
