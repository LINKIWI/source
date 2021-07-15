load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "flags.go",
        "main.go",
    ],
    importpath = "github.com/livegrep/livegrep/cmd/livegrep-metrics-exporter",
    visibility = ["//visibility:private"],
    deps = [
        "@info_kevinlin_lib_aperture//:go_default_library",
        "@info_kevinlin_lib_aperture//lib:go_default_library",
    ],
)

go_binary(
    name = "livegrep-metrics-exporter",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
