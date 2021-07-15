load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "concurrent.go",
        "config.go",
        "flags.go",
        "main.go",
        "vcs.go",
    ],
    importpath = "github.com/livegrep/livegrep/cmd/livegrep-index-builder",
    visibility = ["//visibility:private"],
    deps = [
        "//cmd/livegrep-index-builder/codehost:go_default_library",
        "//src/proto:go_config_proto",
        "@com_github_xanzy_go_gitlab//:go_default_library",
        "@in_gopkg_src_d_go_git_v4//:go_default_library",
        "@in_gopkg_src_d_go_git_v4//plumbing/transport:go_default_library",
        "@in_gopkg_src_d_go_git_v4//plumbing/transport/ssh:go_default_library",
        "@org_golang_x_crypto//ssh:go_default_library",
    ],
)

go_binary(
    name = "livegrep-index-builder",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
