load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/mpawlowski/r2modman-headless
gazelle(name = "gazelle")

go_binary(
    name = "r2modman-headless",
    embed = [":r2modman-headless_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "r2modman-headless_lib",
    srcs = ["main.go"],
    importpath = "github.com/mpawlowski/r2modman-headless",
    visibility = ["//visibility:private"],
    deps = [
        "//r2modman",
        "//zip",
        "@org_uber_go_fx//:fx",
    ],
)
