load(
    "@bazel_gazelle//:deps.bzl",
    "go_repository",
)

def _normalize_repo_name(repo):
    return repo.replace("/", "_").replace("-", "_").replace(".", "_")

def _github(repo, commit):
    name = "com_github_" + _normalize_repo_name(repo)
    importpath = "github.com/" + repo
    return struct(
        name = name,
        commit = commit,
        importpath = importpath,
    )

def _golang_x(pkg, commit):
    name = "org_golang_x_" + pkg
    importpath = "golang.org/x/" + pkg
    return struct(
        name = name,
        commit = commit,
        importpath = importpath,
    )

def _gopkg(repo, commit):
    name = "in_gopkg_" + _normalize_repo_name(repo)
    importpath = "gopkg.in/" + repo
    return struct(
        name = name,
        commit = commit,
        importpath = importpath,
    )

_externals = [
    _golang_x("crypto", "4b2356b1ed79e6be3deca3737a3db3d132d2847a"),
    _golang_x("net", "244492dfa37ae2ce87222fd06250a03160745faa"),
    _golang_x("text", "a9a820217f98f7c8a207ec1e45a874e1fe12c478"),
    _golang_x("oauth2", "a6bd8cefa1811bd24b86f8902872e4e8225f74c4"),
    _golang_x("sys", "33540a1f603772f9d4b761f416f5c10dade23e96"),
    struct(
        name = "org_golang_google_appengine",
        commit = "170382fa85b10b94728989dfcf6cc818b335c952",
        importpath = "google.golang.org/appengine/",
        remote = "https://github.com/golang/appengine",
        vcs = "git",
    ),
    _github("bmizerany/pat", "c068ca2f0aacee5ac3681d68e4d0a003b7d1fd2c"),
    _github("emirpasic/gods", "4e23915b9a82f35f320a68a395a7a5045c826932"),
    _github("facebookgo/clock", "600d898af40aa09a7a93ecb9265d87b0504b6f03"),
    _github("facebookgo/limitgroup", "6abd8d71ec01451d7f1929eacaa263bbe2935d05"),
    _github("facebookgo/muster", "fd3d7953fd52354a74b9f6b3d70d0c9650c4ec2a"),
    _github("google/go-github", "e0066688b631702f66e0435ee1633f9d0091e4b9"),
    _github("google/go-querystring", "53e6ce116135b80d037921a7fdd5138cf32d7a8a"),
    _github("honeycombio/libhoney-go", "a8716c5861ae19c1e2baaad52dd59ba64b902bde"),
    _github("jbenet/go-context", "d14ea06fba99483203c19d92cfcd13ebe73135f4"),
    _github("kevinburke/ssh_config", "01f96b0aa0cdcaa93f9495f89bbc6cb5a992ce6e"),
    _github("mitchellh/go-homedir", "af06845cf3004701891bf4fdb884bfe4920b3727"),
    _github("nelhage/go.cli", "2aeb96ef8025f3646befae8353b90f95e9e79bdc"),
    _github("pelletier/go-buffruneio", "25c428535bd3f55a16f149a9daebd3fa4c5a562b"),
    _github("sergi/go-diff", "58c5cb1602ee9676b5d3590d782bedde80706fcc"),
    _github("src-d/gcfg", "1ac3a1ac202429a54835fe8408a92880156b489d"),
    _github("xanzy/ssh-agent", "6a3e2ff9e7c564f36873c2e36413f634534f1c44"),
    _github("xanzy/go-gitlab", "3f1f63decbcbf8e2ec9bbe1417863a8d51718cd5"),
    _gopkg("alexcesaro/statsd.v2", "7fea3f0d2fab1ad973e641e51dba45443a311a90"),
    _gopkg("check.v1", "20d25e2804050c1cd24a7eea1e7a6447dd0e74ec"),
    _gopkg("src-d/go-billy.v4", "780403cfc1bc95ff4d07e7b26db40a6186c5326e"),
    _gopkg("warnings.v0", "ec4a0fea49c7b46c2aeb0b51aac55779c607e52b"),
    _gopkg("yaml.v3", "a6ecf24a6d716d933bcbc255a2f5d492285b54f5"),
    struct(
        name = "in_gopkg_src_d_go_git_v4",
        version = "v0.0.0-20210509073045-5623572584f2",
        sum = "h1:s+zKHYk7gSXj5SEy3wCL5BrlmZAUQQHodKWp78vM3SQ=",
        importpath = "gopkg.in/src-d/go-git.v4",
        replace = "lib.kevinlin.info/external/go-git",
    ),
    struct(
        name = "org_golang_google_grpc",
        commit = "f74f0337644653eba7923908a4d7f79a4f3a267b",
        importpath = "google.golang.org/grpc",
    ),
    struct(
        name = "info_kevinlin_lib_aperture",
        version = "v0.0.0-20210508182034-84a61db2c3ec",
        sum = "h1:nOHTC7mSHPxWS6q6VZeFBvpGEirD8xBtsZWL4bg0E2g=",
        importpath = "lib.kevinlin.info/aperture",
    ),
]

def go_externals():
    for ext in _externals:
        if hasattr(ext, "vcs"):
            go_repository(
                name = ext.name,
                commit = ext.commit,
                importpath = ext.importpath,
                remote = ext.remote,
                vcs = ext.vcs,
            )
        elif hasattr(ext, "version") and hasattr(ext, "sum") and hasattr(ext, "replace"):
            go_repository(
                name = ext.name,
                version = ext.version,
                sum = ext.sum,
                importpath = ext.importpath,
                replace = ext.replace,
            )
        elif hasattr(ext, "version") and hasattr(ext, "sum"):
            go_repository(
                name = ext.name,
                version = ext.version,
                sum = ext.sum,
                importpath = ext.importpath,
            )
        else:
            go_repository(
                name = ext.name,
                commit = ext.commit,
                importpath = ext.importpath,
            )
