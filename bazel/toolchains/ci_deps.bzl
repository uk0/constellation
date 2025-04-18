"""CI dependencies"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def ci_deps():
    """Install CI dependencies"""
    _shellcheck_deps()
    _terraform_deps()
    _actionlint_deps()
    _gofumpt_deps()
    _tfsec_deps()
    _golangci_lint_deps()
    _buf_deps()
    _talos_docgen_deps()
    _helm_deps()
    _ghh_deps()

def _shellcheck_deps():
    http_archive(
        name = "com_github_koalaman_shellcheck_linux_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/6c881ab0698e4e6ea235245f22832860544f17ba386442fe7e9d629f8cbedf87",
            "https://github.com/koalaman/shellcheck/releases/download/v0.10.0/shellcheck-v0.10.0.linux.x86_64.tar.xz",
        ],
        sha256 = "6c881ab0698e4e6ea235245f22832860544f17ba386442fe7e9d629f8cbedf87",
        strip_prefix = "shellcheck-v0.10.0",
        build_file_content = """exports_files(["shellcheck"], visibility = ["//visibility:public"])""",
        type = "tar.xz",
    )
    http_archive(
        name = "com_github_koalaman_shellcheck_linux_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/324a7e89de8fa2aed0d0c28f3dab59cf84c6d74264022c00c22af665ed1a09bb",
            "https://github.com/koalaman/shellcheck/releases/download/v0.10.0/shellcheck-v0.10.0.linux.aarch64.tar.xz",
        ],
        sha256 = "324a7e89de8fa2aed0d0c28f3dab59cf84c6d74264022c00c22af665ed1a09bb",
        strip_prefix = "shellcheck-v0.10.0",
        build_file_content = """exports_files(["shellcheck"], visibility = ["//visibility:public"])""",
        type = "tar.xz",
    )
    http_archive(
        name = "com_github_koalaman_shellcheck_darwin_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/ef27684f23279d112d8ad84e0823642e43f838993bbb8c0963db9b58a90464c2",
            "https://github.com/koalaman/shellcheck/releases/download/v0.10.0/shellcheck-v0.10.0.darwin.x86_64.tar.xz",
        ],
        sha256 = "ef27684f23279d112d8ad84e0823642e43f838993bbb8c0963db9b58a90464c2",
        strip_prefix = "shellcheck-v0.10.0",
        build_file_content = """exports_files(["shellcheck"], visibility = ["//visibility:public"])""",
        type = "tar.xz",
    )

def _terraform_deps():
    http_archive(
        name = "com_github_hashicorp_terraform_linux_amd64",
        build_file_content = """exports_files(["terraform"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/ad0c696c870c8525357b5127680cd79c0bdf58179af9acd091d43b1d6482da4a",
            "https://releases.hashicorp.com/terraform/1.5.5/terraform_1.5.5_linux_amd64.zip",
        ],
        sha256 = "ad0c696c870c8525357b5127680cd79c0bdf58179af9acd091d43b1d6482da4a",
        type = "zip",
    )
    http_archive(
        name = "com_github_hashicorp_terraform_linux_arm64",
        build_file_content = """exports_files(["terraform"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/b055aefe343d0b710d8a7afd31aeb702b37bbf4493bb9385a709991e48dfbcd2",
            "https://releases.hashicorp.com/terraform/1.5.5/terraform_1.5.5_linux_arm64.zip",
        ],
        sha256 = "b055aefe343d0b710d8a7afd31aeb702b37bbf4493bb9385a709991e48dfbcd2",
        type = "zip",
    )
    http_archive(
        name = "com_github_hashicorp_terraform_darwin_amd64",
        build_file_content = """exports_files(["terraform"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/6d61639e2141b7c23a9219c63994f729aa41f91110a1aa08b8a37969fb45e229",
            "https://releases.hashicorp.com/terraform/1.5.5/terraform_1.5.5_darwin_amd64.zip",
        ],
        sha256 = "6d61639e2141b7c23a9219c63994f729aa41f91110a1aa08b8a37969fb45e229",
        type = "zip",
    )
    http_archive(
        name = "com_github_hashicorp_terraform_darwin_arm64",
        build_file_content = """exports_files(["terraform"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/c7fdeddb4739fdd5bada9d45fd786e2cbaf6e9e364693eee45c83e95281dad3a",
            "https://releases.hashicorp.com/terraform/1.5.5/terraform_1.5.5_darwin_arm64.zip",
        ],
        sha256 = "c7fdeddb4739fdd5bada9d45fd786e2cbaf6e9e364693eee45c83e95281dad3a",
        type = "zip",
    )

def _actionlint_deps():
    http_archive(
        name = "com_github_rhysd_actionlint_linux_amd64",
        build_file_content = """exports_files(["actionlint"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/023070a287cd8cccd71515fedc843f1985bf96c436b7effaecce67290e7e0757",
            "https://github.com/rhysd/actionlint/releases/download/v1.7.7/actionlint_1.7.7_linux_amd64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "023070a287cd8cccd71515fedc843f1985bf96c436b7effaecce67290e7e0757",
    )
    http_archive(
        name = "com_github_rhysd_actionlint_linux_arm64",
        build_file_content = """exports_files(["actionlint"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/401942f9c24ed71e4fe71b76c7d638f66d8633575c4016efd2977ce7c28317d0",
            "https://github.com/rhysd/actionlint/releases/download/v1.7.7/actionlint_1.7.7_linux_arm64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "401942f9c24ed71e4fe71b76c7d638f66d8633575c4016efd2977ce7c28317d0",
    )
    http_archive(
        name = "com_github_rhysd_actionlint_darwin_amd64",
        build_file_content = """exports_files(["actionlint"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/28e5de5a05fc558474f638323d736d822fff183d2d492f0aecb2b73cc44584f5",
            "https://github.com/rhysd/actionlint/releases/download/v1.7.7/actionlint_1.7.7_darwin_amd64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "28e5de5a05fc558474f638323d736d822fff183d2d492f0aecb2b73cc44584f5",
    )
    http_archive(
        name = "com_github_rhysd_actionlint_darwin_arm64",
        build_file_content = """exports_files(["actionlint"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/2693315b9093aeacb4ebd91a993fea54fc215057bf0da2659056b4bc033873db",
            "https://github.com/rhysd/actionlint/releases/download/v1.7.7/actionlint_1.7.7_darwin_arm64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "2693315b9093aeacb4ebd91a993fea54fc215057bf0da2659056b4bc033873db",
    )

def _gofumpt_deps():
    http_file(
        name = "com_github_mvdan_gofumpt_linux_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/11604bbaf7321abcc2fca2c6a37b7e9198bb1e76e5a86f297c07201e8ab1fda9",
            "https://github.com/mvdan/gofumpt/releases/download/v0.8.0/gofumpt_v0.8.0_linux_amd64",
        ],
        executable = True,
        downloaded_file_path = "gofumpt",
        sha256 = "11604bbaf7321abcc2fca2c6a37b7e9198bb1e76e5a86f297c07201e8ab1fda9",
    )
    http_file(
        name = "com_github_mvdan_gofumpt_linux_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/787c1d3d4d20e6fe2b0bf06a5a913ac0f50343dbf9a71540724a2b8092a0e6ca",
            "https://github.com/mvdan/gofumpt/releases/download/v0.8.0/gofumpt_v0.8.0_linux_arm64",
        ],
        executable = True,
        downloaded_file_path = "gofumpt",
        sha256 = "787c1d3d4d20e6fe2b0bf06a5a913ac0f50343dbf9a71540724a2b8092a0e6ca",
    )
    http_file(
        name = "com_github_mvdan_gofumpt_darwin_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/0dda6600cf263b703a5ad93e792b06180c36afdee9638617a91dd552f2c6fb3e",
            "https://github.com/mvdan/gofumpt/releases/download/v0.8.0/gofumpt_v0.8.0_darwin_amd64",
        ],
        executable = True,
        downloaded_file_path = "gofumpt",
        sha256 = "0dda6600cf263b703a5ad93e792b06180c36afdee9638617a91dd552f2c6fb3e",
    )
    http_file(
        name = "com_github_mvdan_gofumpt_darwin_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/7e66e92b7a67d1d12839ab030fb7ae38e5e2273474af3762e67bc7fe9471fcd9",
            "https://github.com/mvdan/gofumpt/releases/download/v0.8.0/gofumpt_v0.8.0_darwin_arm64",
        ],
        executable = True,
        downloaded_file_path = "gofumpt",
        sha256 = "7e66e92b7a67d1d12839ab030fb7ae38e5e2273474af3762e67bc7fe9471fcd9",
    )

def _tfsec_deps():
    http_archive(
        name = "com_github_aquasecurity_tfsec_linux_amd64",
        build_file_content = """exports_files(["tfsec"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/f643c390e6eabdf4bd64e807ff63abfe977d4f027c0b535eefe7a5c9f391fc28",
            "https://github.com/aquasecurity/tfsec/releases/download/v1.28.13/tfsec_1.28.13_linux_amd64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "f643c390e6eabdf4bd64e807ff63abfe977d4f027c0b535eefe7a5c9f391fc28",
    )
    http_archive(
        name = "com_github_aquasecurity_tfsec_linux_arm64",
        build_file_content = """exports_files(["tfsec"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/4aed1b122f817b684cc48da9cdc4b98b891808969441914570c089883c85ac50",
            "https://github.com/aquasecurity/tfsec/releases/download/v1.28.13/tfsec_1.28.13_linux_arm64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "4aed1b122f817b684cc48da9cdc4b98b891808969441914570c089883c85ac50",
    )
    http_archive(
        name = "com_github_aquasecurity_tfsec_darwin_amd64",
        build_file_content = """exports_files(["tfsec"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/966c7656f797c120dceb56a208a50dbf6a363c30876662a28e1c65505afca10d",
            "https://github.com/aquasecurity/tfsec/releases/download/v1.28.13/tfsec_1.28.13_darwin_amd64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "966c7656f797c120dceb56a208a50dbf6a363c30876662a28e1c65505afca10d",
    )
    http_archive(
        name = "com_github_aquasecurity_tfsec_darwin_arm64",
        build_file_content = """exports_files(["tfsec"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/a381580c81d3413bb3fe07aa91ab89e51c1bbbd33c848194a2b43e9be3729c16",
            "https://github.com/aquasecurity/tfsec/releases/download/v1.28.13/tfsec_1.28.13_darwin_arm64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "a381580c81d3413bb3fe07aa91ab89e51c1bbbd33c848194a2b43e9be3729c16",
    )

def _golangci_lint_deps():
    http_archive(
        name = "com_github_golangci_golangci_lint_linux_amd64",
        build_file = "//bazel/toolchains:BUILD.golangci.bazel",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/bc16fd1ef25bce2c600de0122600100ab26d6d75388cc5369c5bb916cb2b82e3",
            "https://github.com/golangci/golangci-lint/releases/download/v2.1.2/golangci-lint-2.1.2-linux-amd64.tar.gz",
        ],
        strip_prefix = "golangci-lint-2.1.2-linux-amd64",
        type = "tar.gz",
        sha256 = "bc16fd1ef25bce2c600de0122600100ab26d6d75388cc5369c5bb916cb2b82e3",
    )
    http_archive(
        name = "com_github_golangci_golangci_lint_linux_arm64",
        build_file = "//bazel/toolchains:BUILD.golangci.bazel",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/46e86f1c4a94236e4d0bb35252c72939bed9f749897aaad54b576d430b1bb6d4",
            "https://github.com/golangci/golangci-lint/releases/download/v2.1.2/golangci-lint-2.1.2-linux-arm64.tar.gz",
        ],
        strip_prefix = "golangci-lint-2.1.2-linux-arm64",
        type = "tar.gz",
        sha256 = "46e86f1c4a94236e4d0bb35252c72939bed9f749897aaad54b576d430b1bb6d4",
    )
    http_archive(
        name = "com_github_golangci_golangci_lint_darwin_amd64",
        build_file = "//bazel/toolchains:BUILD.golangci.bazel",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/ed02ba3ad28466d61d2ae2b80cc95671713fa842c353da37842b1b89e36cb3ce",
            "https://github.com/golangci/golangci-lint/releases/download/v2.1.2/golangci-lint-2.1.2-darwin-amd64.tar.gz",
        ],
        strip_prefix = "golangci-lint-2.1.2-darwin-amd64",
        type = "tar.gz",
        sha256 = "ed02ba3ad28466d61d2ae2b80cc95671713fa842c353da37842b1b89e36cb3ce",
    )
    http_archive(
        name = "com_github_golangci_golangci_lint_darwin_arm64",
        build_file = "//bazel/toolchains:BUILD.golangci.bazel",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/1cff60651d7c95a4248fa72f0dd020bffed1d2dc4dd8c2c77aee89a0731fa615",
            "https://github.com/golangci/golangci-lint/releases/download/v2.1.2/golangci-lint-2.1.2-darwin-arm64.tar.gz",
        ],
        strip_prefix = "golangci-lint-2.1.2-darwin-arm64",
        type = "tar.gz",
        sha256 = "1cff60651d7c95a4248fa72f0dd020bffed1d2dc4dd8c2c77aee89a0731fa615",
    )

def _buf_deps():
    http_archive(
        name = "com_github_bufbuild_buf_linux_amd64",
        strip_prefix = "buf/bin",
        build_file_content = """exports_files(["buf"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/ce439b98c0e655ba1821fc529442f677135a9a44206a85cf2081a445bde02554",
            "https://github.com/bufbuild/buf/releases/download/v1.52.1/buf-Linux-x86_64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "ce439b98c0e655ba1821fc529442f677135a9a44206a85cf2081a445bde02554",
    )
    http_archive(
        name = "com_github_bufbuild_buf_linux_arm64",
        strip_prefix = "buf/bin",
        build_file_content = """exports_files(["buf"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/fb980e62f55a0d705f555124b5b250a6234deae67b3a493b1844fed8cb48f012",
            "https://github.com/bufbuild/buf/releases/download/v1.52.1/buf-Linux-aarch64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "fb980e62f55a0d705f555124b5b250a6234deae67b3a493b1844fed8cb48f012",
    )
    http_archive(
        name = "com_github_bufbuild_buf_darwin_amd64",
        strip_prefix = "buf/bin",
        build_file_content = """exports_files(["buf"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/ff3aedaa8bf53c3ae53bf863bb6141a23ef9e2ef750016ab217b0466ab3e89be",
            "https://github.com/bufbuild/buf/releases/download/v1.52.1/buf-Darwin-x86_64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "ff3aedaa8bf53c3ae53bf863bb6141a23ef9e2ef750016ab217b0466ab3e89be",
    )
    http_archive(
        name = "com_github_bufbuild_buf_darwin_arm64",
        strip_prefix = "buf/bin",
        build_file_content = """exports_files(["buf"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/b11d5afda7f58a44a6520b798902b1263aba921b0989cb156ecf2acc8eb1ce31",
            "https://github.com/bufbuild/buf/releases/download/v1.52.1/buf-Darwin-arm64.tar.gz",
        ],
        type = "tar.gz",
        sha256 = "b11d5afda7f58a44a6520b798902b1263aba921b0989cb156ecf2acc8eb1ce31",
    )

def _talos_docgen_deps():
    http_file(
        name = "com_github_siderolabs_talos_hack_docgen_linux_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/bd1059b49a6db7473b4f991a244e338da887d8017ab556739000abd2cc367c13",
        ],
        executable = True,
        downloaded_file_path = "docgen",
        sha256 = "bd1059b49a6db7473b4f991a244e338da887d8017ab556739000abd2cc367c13",
    )
    http_file(
        name = "com_github_siderolabs_talos_hack_docgen_darwin_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/d06adae41a975a94abaa39cd809464d7a8f7648903d321332c12c73002cc622a",
        ],
        executable = True,
        downloaded_file_path = "docgen",
        sha256 = "d06adae41a975a94abaa39cd809464d7a8f7648903d321332c12c73002cc622a",
    )
    http_file(
        name = "com_github_siderolabs_talos_hack_docgen_linux_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/a87c52a3e947fe90396427a5cd92e6864f46b5db103f84c1cad449e97ca54cec",
        ],
        executable = True,
        downloaded_file_path = "docgen",
        sha256 = "a87c52a3e947fe90396427a5cd92e6864f46b5db103f84c1cad449e97ca54cec",
    )
    http_file(
        name = "com_github_siderolabs_talos_hack_docgen_darwin_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/4aa7ed0de31932d541aa11c9b75ed214ffc28dbd618f489fb5a598407aca072e",
        ],
        executable = True,
        downloaded_file_path = "docgen",
        sha256 = "4aa7ed0de31932d541aa11c9b75ed214ffc28dbd618f489fb5a598407aca072e",
    )

def _helm_deps():
    http_archive(
        name = "com_github_helm_helm_linux_amd64",
        sha256 = "a5844ef2c38ef6ddf3b5a8f7d91e7e0e8ebc39a38bb3fc8013d629c1ef29c259",
        strip_prefix = "linux-amd64",
        build_file_content = """exports_files(["helm"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/a5844ef2c38ef6ddf3b5a8f7d91e7e0e8ebc39a38bb3fc8013d629c1ef29c259",
            "https://get.helm.sh/helm-v3.14.4-linux-amd64.tar.gz",
        ],
        type = "tar.gz",
    )
    http_archive(
        name = "com_github_helm_helm_linux_arm64",
        sha256 = "113ccc53b7c57c2aba0cd0aa560b5500841b18b5210d78641acfddc53dac8ab2",
        strip_prefix = "linux-arm64",
        build_file_content = """exports_files(["helm"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/113ccc53b7c57c2aba0cd0aa560b5500841b18b5210d78641acfddc53dac8ab2",
            "https://get.helm.sh/helm-v3.14.4-linux-arm64.tar.gz",
        ],
        type = "tar.gz",
    )
    http_archive(
        name = "com_github_helm_helm_darwin_amd64",
        sha256 = "73434aeac36ad068ce2e5582b8851a286dc628eae16494a26e2ad0b24a7199f9",
        strip_prefix = "darwin-amd64",
        build_file_content = """exports_files(["helm"], visibility = ["//visibility:public"])""",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/73434aeac36ad068ce2e5582b8851a286dc628eae16494a26e2ad0b24a7199f9",
            "https://get.helm.sh/helm-v3.14.4-darwin-amd64.tar.gz",
        ],
        type = "tar.gz",
    )
    http_archive(
        name = "com_github_helm_helm_darwin_arm64",
        sha256 = "61e9c5455f06b2ad0a1280975bf65892e707adc19d766b0cf4e9006e3b7b4b6c",
        strip_prefix = "darwin-arm64",
        build_file_content = """exports_files(["helm"], visibility = ["//visibility:public"])""",
        type = "tar.gz",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/61e9c5455f06b2ad0a1280975bf65892e707adc19d766b0cf4e9006e3b7b4b6c",
            "https://get.helm.sh/helm-v3.14.4-darwin-arm64.tar.gz",
        ],
    )

def _ghh_deps():
    http_archive(
        name = "com_github_katexochen_ghh_linux_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/32f8a4110d88d80e163212a89a3538a13326494840ac97183d1b20bcc9eac7ba",
            "https://github.com/katexochen/ghh/releases/download/v0.3.5/ghh_0.3.5_linux_amd64.tar.gz",
        ],
        type = "tar.gz",
        build_file_content = """exports_files(["ghh"], visibility = ["//visibility:public"])""",
        sha256 = "32f8a4110d88d80e163212a89a3538a13326494840ac97183d1b20bcc9eac7ba",
    )
    http_archive(
        name = "com_github_katexochen_ghh_linux_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/b43ef1dd2f851eed7c69c87f4f73dd923bd1170cefbde247933d5b398a3319d1",
            "https://github.com/katexochen/ghh/releases/download/v0.3.5/ghh_0.3.5_linux_arm64.tar.gz",
        ],
        type = "tar.gz",
        build_file_content = """exports_files(["ghh"], visibility = ["//visibility:public"])""",
        sha256 = "b43ef1dd2f851eed7c69c87f4f73dd923bd1170cefbde247933d5b398a3319d1",
    )
    http_archive(
        name = "com_github_katexochen_ghh_darwin_amd64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/7db9ebb62faf2a31f56e7a8994a971adddec98b3238880ae58b01eb549b8bba3",
            "https://github.com/katexochen/ghh/releases/download/v0.3.5/ghh_0.3.5_darwin_amd64.tar.gz",
        ],
        type = "tar.gz",
        build_file_content = """exports_files(["ghh"], visibility = ["//visibility:public"])""",
        sha256 = "7db9ebb62faf2a31f56e7a8994a971adddec98b3238880ae58b01eb549b8bba3",
    )
    http_archive(
        name = "com_github_katexochen_ghh_darwin_arm64",
        urls = [
            "https://cdn.confidential.cloud/constellation/cas/sha256/78a2c8b321893736a2b5de59898794a72b878db9329157f348489c73d4592c6f",
            "https://github.com/katexochen/ghh/releases/download/v0.3.5/ghh_0.3.5_darwin_arm64.tar.gz",
        ],
        type = "tar.gz",
        build_file_content = """exports_files(["ghh"], visibility = ["//visibility:public"])""",
        sha256 = "78a2c8b321893736a2b5de59898794a72b878db9329157f348489c73d4592c6f",
    )
