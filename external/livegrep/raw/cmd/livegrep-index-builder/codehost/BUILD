load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "gitlab.go",
        "project.go",
        "static.go",
        "types.go",
    ],
    importpath = "github.com/livegrep/livegrep/cmd/livegrep-index-builder/codehost",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_xanzy_go_gitlab//:go_default_library",
        "@in_gopkg_yaml_v3//:go_default_library",
    ],
)
