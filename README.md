[![CircleCI](https://circleci.com/gh/angular/angular-bazel-example.svg?style=svg)](https://circleci.com/gh/angular/angular-bazel-example)

# Example Bazel + Angular+ Golang + Protocol buffers + gRPC stack

This doesn't work yet.

## Clipboard farm

**All of these commands should be run from the viz directory.**

Run the devserver, get an error about failing to load dependencies:

```shell
bazel run //src:devserver
```

Error message:

    Failed to load resource: the server responded with a status of 404 (Not Found)
    zone.min.js:19 Uncaught Error: Script error for "google-protobuf", needed by: examples_angular/httpserver/frontendpb/frontend_pb
        http://requirejs.org/docs/errors.html#scripterror

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
