# Changelog

## [0.5.0](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.4.0...ohpygossh-v0.5.0) (2026-07-09)


### Features

* support testing with LXD ([#68](https://github.com/b-long/ohpygossh/issues/68)) ([4995461](https://github.com/b-long/ohpygossh/commit/4995461e691df5971f8fb5dc4df7b92c6a17e651))
* update to go 1.25 ([#73](https://github.com/b-long/ohpygossh/issues/73)) ([055932b](https://github.com/b-long/ohpygossh/commit/055932b878995ad066bcce68bba6148e17742f41))
* update to latest gopy main (2026 June 02) ([#59](https://github.com/b-long/ohpygossh/issues/59)) ([d6c8cee](https://github.com/b-long/ohpygossh/commit/d6c8cee5643e575d11894f891453ba656194bbf3))

## [0.4.0](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.3.2...ohpygossh-v0.4.0) (2026-06-10)


### Features

* improve README for end-users ([#54](https://github.com/b-long/ohpygossh/issues/54)) ([120e815](https://github.com/b-long/ohpygossh/commit/120e815c9af6e42509ede0a85716f4a8d2f9abec))

## [0.3.2](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.3.1...ohpygossh-v0.3.2) (2026-05-29)


### Bug Fixes

* create new func for VM naming ([#51](https://github.com/b-long/ohpygossh/issues/51)) ([a1e8b86](https://github.com/b-long/ohpygossh/commit/a1e8b86c5923f4ea3ef8e0b7471bd76220437dc8))
* sync pyproject.toml to release-please manifest ([#50](https://github.com/b-long/ohpygossh/issues/50)) ([b7d6164](https://github.com/b-long/ohpygossh/commit/b7d61645f132338fb15ed3c616f59eae935e45e2))
* use explicit toml type in extra-files for pyproject.toml ([#52](https://github.com/b-long/ohpygossh/issues/52)) ([42b7084](https://github.com/b-long/ohpygossh/commit/42b70842a88a0d56627f44d5df720034b1813fb0))
* various cibuildwheel fixes ([#48](https://github.com/b-long/ohpygossh/issues/48)) ([338ee68](https://github.com/b-long/ohpygossh/commit/338ee68871eb6d834653c45ee2f367137be7d68b))

## [0.3.1](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.3.0...ohpygossh-v0.3.1) (2026-05-22)


### Bug Fixes

* add actions: write permission to dispatch publish workflow ([#31](https://github.com/b-long/ohpygossh/issues/31)) ([731e93a](https://github.com/b-long/ohpygossh/commit/731e93af4ac28770f5afbb67b7e763ee5778ce82))
* pass --repo to gh workflow run to avoid missing git context ([#30](https://github.com/b-long/ohpygossh/issues/30)) ([ccfb0c6](https://github.com/b-long/ohpygossh/commit/ccfb0c67f0e4a761909d88b20fdde39590c5858e))
* remove colima, docker, and vagrant steps from publish workflow ([#32](https://github.com/b-long/ohpygossh/issues/32)) ([fd6dee2](https://github.com/b-long/ohpygossh/commit/fd6dee25b348e2911f69ab94f5eb9278eb457e58))
* trigger publish via workflow_dispatch from release-please ([#28](https://github.com/b-long/ohpygossh/issues/28)) ([3b3e254](https://github.com/b-long/ohpygossh/commit/3b3e254a92c013c9624c540dc11226d701c4edfa))

## [0.3.0](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.2.0...ohpygossh-v0.3.0) (2026-05-22)


### Features

* add SFTP Upload and Download functions ([#26](https://github.com/b-long/ohpygossh/issues/26)) ([0cf669b](https://github.com/b-long/ohpygossh/commit/0cf669b0acef3424180e4c68a552c7591de4caaf))


### Bug Fixes

* PyPI publishing & tidy go-multipass tests ([#24](https://github.com/b-long/ohpygossh/issues/24)) ([5df9c8b](https://github.com/b-long/ohpygossh/commit/5df9c8bf47c83dad88b25d751cad5842de2c2bb5))
* sync pyproject.toml to 0.2.0, fix extra-files ([#27](https://github.com/b-long/ohpygossh/issues/27)) ([40ca820](https://github.com/b-long/ohpygossh/commit/40ca820adb0af974eabdedb8f56a63fbecc90bb7))

## [0.2.0](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.1.0...ohpygossh-v0.2.0) (2026-05-21)


### Features

* add multipass e2e test and CI workflows ([#23](https://github.com/b-long/ohpygossh/issues/23)) ([69b21d1](https://github.com/b-long/ohpygossh/commit/69b21d1c16b8edbd8285a008dff56a1f99cbc84c))
* expand go-lang tests ([#21](https://github.com/b-long/ohpygossh/issues/21)) ([7538ae5](https://github.com/b-long/ohpygossh/commit/7538ae5a9131ab9b82feacfef517eda9cbc10d63))

## [0.1.0](https://github.com/b-long/ohpygossh/compare/ohpygossh-v0.0.16...ohpygossh-v0.1.0) (2026-05-20)


### Bug Fixes

* drop TOML jsonpath & add version.txt  ([#19](https://github.com/b-long/ohpygossh/issues/19)) ([745513e](https://github.com/b-long/ohpygossh/commit/745513ebdbaf75f558a5bb67df03c1490a5de902))
* drop TOML jsonpath for pyproject.toml ([#18](https://github.com/b-long/ohpygossh/issues/18)) ([b1e8448](https://github.com/b-long/ohpygossh/commit/b1e844871131cd80e354e62c80fed6d50dee856c))
