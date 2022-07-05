load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "github.go",
        "gitlab.go",
        "project.go",
        "static.go",
        "types.go",
    ],
    importpath = "github.com/livegrep/livegrep/cmd/livegrep-index-builder/codehost",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_google_go_github//github:go_default_library",
        "@com_github_xanzy_go_gitlab//:go_default_library",
        "@in_gopkg_yaml_v3//:go_default_library",
        "@org_golang_x_oauth2//:go_default_library",
    ],
)
