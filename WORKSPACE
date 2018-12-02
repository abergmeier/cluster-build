

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.16.2/rules_go-0.16.2.tar.gz"],
    sha256 = "f87fa87475ea107b3c69196f39c82b7bbf58fe27c62a338684c20ca17d1d8613",
)
http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.15.0/bazel-gazelle-0.15.0.tar.gz"],
    sha256 = "6e875ab4b6bf64a38c352887760f21203ab054676d9c1b274963907e0768740d",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@bazel_gazelle//:deps.bzl", "go_repository")

go_repository(
    name = "org_golang_google_api",
    commit = "faade3cbb06a30202f2da53a8a5e3c4afe60b0c2",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "org_golang_google_genproto",
    commit = "b5d43981345bdb2c233eb4bf3277847b48c6fdc6",
    importpath = "google.golang.org/genproto",
)

http_archive(
    name = "com_google_googleapis",
    #urls = ["https://github.com/googleapis/googleapis/archive/3cef6c237ba87a81470b5fdf082ed187cbb23ded.zip"],
    urls = ["https://github.com/abergmeier/googleapis/archive/30eff4f3990c6bf84159d574f79d8342772779dc.zip"],
    sha256 = "594cc7269feed61a1e14dcc1d117d4c199d27201bdf99f4da05ab7661ad69e6d",
    strip_prefix = "googleapis-30eff4f3990c6bf84159d574f79d8342772779dc",
)

#
# grpc-java repository dependencies (required to by `java_grpc_library` bazel rule)
#
git_repository(
    name = "io_grpc_grpc_java",
    remote = "https://github.com/grpc/grpc-java.git",
    tag = "v1.13.1",
)
