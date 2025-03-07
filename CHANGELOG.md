# Change log

All notable changes to the ld-find-code-refs program will be documented in this file. This project adheres to [Semantic Versioning](http://semver.org).

## [0.0.1](https://github.com/bucketeer-io/code-refs/compare/v0.0.1...v0.0.1) (2025-03-07)


### Features

* bucketeer integration ([#1](https://github.com/bucketeer-io/code-refs/issues/1)) ([3926f23](https://github.com/bucketeer-io/code-refs/commit/3926f23ea27ecc8c63ed95e7da7cc14e6a01a507))
* don't skip .github dir ([#441](https://github.com/bucketeer-io/code-refs/issues/441)) ([afe83e5](https://github.com/bucketeer-io/code-refs/commit/afe83e5566b4136293ad064b44fd8a2bc9fb2ddf))
* expose alias generator functions ([#406](https://github.com/bucketeer-io/code-refs/issues/406)) ([68be382](https://github.com/bucketeer-io/code-refs/commit/68be382c08706ab1c4dc3480894a6222ae22da26))
* respect X-Ratelimit-Reset when retrying API requests ([#434](https://github.com/bucketeer-io/code-refs/issues/434)) ([650c6a7](https://github.com/bucketeer-io/code-refs/commit/650c6a715b5dc5f074ca4286095d41c46c63e276))


### Bug Fixes

* **#429:** Multi-project scans now ignore projects without valid flags ([#430](https://github.com/bucketeer-io/code-refs/issues/430)) ([6dca539](https://github.com/bucketeer-io/code-refs/commit/6dca53907a8cdfb9793df7f8cfd2ee64828d024b))
* extinctions scans take hours ([#331](https://github.com/bucketeer-io/code-refs/issues/331)) ([3af7d0c](https://github.com/bucketeer-io/code-refs/commit/3af7d0cc329b293d5af3e7c662e32ce9618f3b03))
* failed to find tag name for annotated tags ([#334](https://github.com/bucketeer-io/code-refs/issues/334)) ([3cb948a](https://github.com/bucketeer-io/code-refs/commit/3cb948a28f7b62cca278a315eb853e5e584d7e0e))
* incorrect names for some release assets ([#385](https://github.com/bucketeer-io/code-refs/issues/385)) ([091c549](https://github.com/bucketeer-io/code-refs/commit/091c5491c8a747a1905e691180d739df9ee5b91b))
* use RateLimitBackoff client in LD API ([#437](https://github.com/bucketeer-io/code-refs/issues/437)) ([04ccef1](https://github.com/bucketeer-io/code-refs/commit/04ccef115b4a62e808a0b4ddc84dea7464f18290))


### Miscellaneous

* Add release version to commit message ([#428](https://github.com/bucketeer-io/code-refs/issues/428)) ([b8dd0ad](https://github.com/bucketeer-io/code-refs/commit/b8dd0adfc29de22ae3545e64a6cf4b9a4ac99b97))
* bump actions/checkout version in example ([#394](https://github.com/bucketeer-io/code-refs/issues/394)) ([e2663b1](https://github.com/bucketeer-io/code-refs/commit/e2663b1b5a7f0519b799bb954e2c5e8f4f60a7c8))
* **deps:** bump alpine from 3.18.4 to 3.18.5 ([#414](https://github.com/bucketeer-io/code-refs/issues/414)) ([8916082](https://github.com/bucketeer-io/code-refs/commit/8916082a42015690d38490f0fb809e38b04cfbbf))
* **deps:** bump alpine from 3.18.5 to 3.19.0 ([#418](https://github.com/bucketeer-io/code-refs/issues/418)) ([634c272](https://github.com/bucketeer-io/code-refs/commit/634c2722413a13f037d6eca700d5457ab4cd1018))
* **deps:** bump alpine from 3.19.0 to 3.19.1 ([#426](https://github.com/bucketeer-io/code-refs/issues/426)) ([60e141f](https://github.com/bucketeer-io/code-refs/commit/60e141f1353cd3513c010888e792c5318f72cdf3))
* **deps:** bump alpine from 3.19.1 to 3.20.0 ([#448](https://github.com/bucketeer-io/code-refs/issues/448)) ([ef25e95](https://github.com/bucketeer-io/code-refs/commit/ef25e955f6106b1f371f863c5534ba4feaf82cc6))
* **deps:** bump alpine from 3.20.0 to 3.20.1 ([#453](https://github.com/bucketeer-io/code-refs/issues/453)) ([ba766ca](https://github.com/bucketeer-io/code-refs/commit/ba766caa740c7e2fd47f008a437918b3fabc9607))
* **deps:** bump alpine from 3.20.1 to 3.20.2 ([#456](https://github.com/bucketeer-io/code-refs/issues/456)) ([d98e82a](https://github.com/bucketeer-io/code-refs/commit/d98e82a71f06f9c5cd5b5388a999933753f2dae8))
* **deps:** bump alpine from 3.20.2 to 3.20.3 ([#458](https://github.com/bucketeer-io/code-refs/issues/458)) ([42ccee7](https://github.com/bucketeer-io/code-refs/commit/42ccee78e0ea3ba9c893487e66976c106c3c0579))
* **deps:** bump alpine from 3.20.3 to 3.21.0 ([#470](https://github.com/bucketeer-io/code-refs/issues/470)) ([8907421](https://github.com/bucketeer-io/code-refs/commit/89074211afd91d44b863a658b593180a985d9c8d))
* **deps:** bump github.com/bmatcuk/doublestar/v4 from 4.6.1 to 4.7.1 ([#463](https://github.com/bucketeer-io/code-refs/issues/463)) ([7fcd5be](https://github.com/bucketeer-io/code-refs/commit/7fcd5befcccb1c46de69d508eb1b54330f2dfc42))
* **deps:** bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 ([#422](https://github.com/bucketeer-io/code-refs/issues/422)) ([d73aa5d](https://github.com/bucketeer-io/code-refs/commit/d73aa5dc9816acbe6d890b522257bd958b156894))
* **deps:** bump github.com/go-git/go-git/v5 from 5.10.0 to 5.10.1 ([#413](https://github.com/bucketeer-io/code-refs/issues/413)) ([3d56771](https://github.com/bucketeer-io/code-refs/commit/3d5677142bdc06f2d5313e3a95bdc4f4c3be12ce))
* **deps:** bump github.com/go-git/go-git/v5 from 5.10.1 to 5.11.0 ([#417](https://github.com/bucketeer-io/code-refs/issues/417)) ([ad45639](https://github.com/bucketeer-io/code-refs/commit/ad45639d8ee25d9a4ba37cb17f2e56a133d94dd7))
* **deps:** bump github.com/go-git/go-git/v5 from 5.11.0 to 5.12.0 ([#443](https://github.com/bucketeer-io/code-refs/issues/443)) ([263a831](https://github.com/bucketeer-io/code-refs/commit/263a8313a41a2c13df258784c6c5212915506401))
* **deps:** bump github.com/go-git/go-git/v5 from 5.12.0 to 5.13.1 ([#473](https://github.com/bucketeer-io/code-refs/issues/473)) ([4413d04](https://github.com/bucketeer-io/code-refs/commit/4413d049c0a5bfc7774b158dfb90c0c4bcca8efa))
* **deps:** bump github.com/hashicorp/go-retryablehttp from 0.7.4 to 0.7.5 ([#410](https://github.com/bucketeer-io/code-refs/issues/410)) ([b78df54](https://github.com/bucketeer-io/code-refs/commit/b78df5411a82e164b3c9604d7b971e0c36347833))
* **deps:** bump github.com/hashicorp/go-retryablehttp from 0.7.5 to 0.7.6 ([#447](https://github.com/bucketeer-io/code-refs/issues/447)) ([579f7c1](https://github.com/bucketeer-io/code-refs/commit/579f7c16fba1e53cbd3262bed52e73c2e88a7fc2))
* **deps:** bump github.com/hashicorp/go-retryablehttp from 0.7.6 to 0.7.7 ([#449](https://github.com/bucketeer-io/code-refs/issues/449)) ([04272e4](https://github.com/bucketeer-io/code-refs/commit/04272e4846ead5113edd1337bc3e144cc0dd9b00))
* **deps:** bump github.com/launchdarkly/api-client-go/v15 from 15.0.0 to 15.1.0 ([#444](https://github.com/bucketeer-io/code-refs/issues/444)) ([cf3299b](https://github.com/bucketeer-io/code-refs/commit/cf3299bff852b309d7e1f6ddc6311142b872bfdb))
* **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([#409](https://github.com/bucketeer-io/code-refs/issues/409)) ([f8a724a](https://github.com/bucketeer-io/code-refs/commit/f8a724a29af7f79d62d77b2931617414484d6c52))
* **deps:** bump github.com/spf13/cobra from 1.8.0 to 1.8.1 ([#452](https://github.com/bucketeer-io/code-refs/issues/452)) ([c74ab9a](https://github.com/bucketeer-io/code-refs/commit/c74ab9a9638f15c7f95a0f514be2ac604a178dda))
* **deps:** bump github.com/spf13/viper from 1.17.0 to 1.18.0 ([#415](https://github.com/bucketeer-io/code-refs/issues/415)) ([c251dc3](https://github.com/bucketeer-io/code-refs/commit/c251dc3cfea4158718f39550a0cb9afb9c8123ed))
* **deps:** bump github.com/spf13/viper from 1.18.0 to 1.18.1 ([#416](https://github.com/bucketeer-io/code-refs/issues/416)) ([88dfaba](https://github.com/bucketeer-io/code-refs/commit/88dfaba9f285aa37b228f56e3229ee92fec5e861))
* **deps:** bump github.com/spf13/viper from 1.18.1 to 1.18.2 ([#420](https://github.com/bucketeer-io/code-refs/issues/420)) ([dd5b7e8](https://github.com/bucketeer-io/code-refs/commit/dd5b7e8834bb40f9b18a36726a60d0b6bcee99d2))
* **deps:** bump github.com/spf13/viper from 1.18.2 to 1.19.0 ([#450](https://github.com/bucketeer-io/code-refs/issues/450)) ([c08cc73](https://github.com/bucketeer-io/code-refs/commit/c08cc73fe72d551a8903492bc45f23af13efe876))
* **deps:** bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#435](https://github.com/bucketeer-io/code-refs/issues/435)) ([df5ffdf](https://github.com/bucketeer-io/code-refs/commit/df5ffdf08b0a387e0e052667787014f30c63aff6))
* **deps:** bump github.com/stretchr/testify from 1.9.0 to 1.10.0 ([#466](https://github.com/bucketeer-io/code-refs/issues/466)) ([c5d2f8b](https://github.com/bucketeer-io/code-refs/commit/c5d2f8b9a6fc4e596f2de4c51a2046a521036257))
* **deps:** bump golang.org/x/crypto from 0.16.0 to 0.17.0 ([#421](https://github.com/bucketeer-io/code-refs/issues/421)) ([e38da79](https://github.com/bucketeer-io/code-refs/commit/e38da79c92a352faa5ca098db32ff9dcf7461be5))
* **deps:** bump golang.org/x/tools from 0.14.0 to 0.15.0 ([#411](https://github.com/bucketeer-io/code-refs/issues/411)) ([1bce9ce](https://github.com/bucketeer-io/code-refs/commit/1bce9ce2e69d46d30c44088ea024b7557a738eab))
* **deps:** bump golang.org/x/tools from 0.15.0 to 0.16.0 ([#412](https://github.com/bucketeer-io/code-refs/issues/412)) ([0873096](https://github.com/bucketeer-io/code-refs/commit/0873096b947159c4bcafe4ac707f80e730e35920))
* **deps:** bump golang.org/x/tools from 0.16.0 to 0.16.1 ([#419](https://github.com/bucketeer-io/code-refs/issues/419)) ([2a75be3](https://github.com/bucketeer-io/code-refs/commit/2a75be30e77ff1059aa8a9187ebf557dc3002afa))
* **deps:** bump golang.org/x/tools from 0.16.1 to 0.17.0 ([#424](https://github.com/bucketeer-io/code-refs/issues/424)) ([38b8388](https://github.com/bucketeer-io/code-refs/commit/38b83881538a3f20b355aabeadd3e30d2a1acc26))
* **deps:** bump golang.org/x/tools from 0.17.0 to 0.18.0 ([#431](https://github.com/bucketeer-io/code-refs/issues/431)) ([3900e92](https://github.com/bucketeer-io/code-refs/commit/3900e92912739a49293dd90e82edc82534f61c71))
* **deps:** bump golang.org/x/tools from 0.18.0 to 0.19.0 ([#436](https://github.com/bucketeer-io/code-refs/issues/436)) ([29d1c5f](https://github.com/bucketeer-io/code-refs/commit/29d1c5feaabfb31dab3d78518e86f58c344757e4))
* **deps:** bump golang.org/x/tools from 0.19.0 to 0.20.0 ([#445](https://github.com/bucketeer-io/code-refs/issues/445)) ([09e5ed3](https://github.com/bucketeer-io/code-refs/commit/09e5ed3e6016960df89b6f3526799369d96838fa))
* **deps:** bump golang.org/x/tools from 0.20.0 to 0.21.0 ([#446](https://github.com/bucketeer-io/code-refs/issues/446)) ([ae06277](https://github.com/bucketeer-io/code-refs/commit/ae062779bf342cd244f7363b6ff49b8af389242f))
* **deps:** bump golang.org/x/tools from 0.21.0 to 0.22.0 ([#451](https://github.com/bucketeer-io/code-refs/issues/451)) ([e3db039](https://github.com/bucketeer-io/code-refs/commit/e3db03986495e9bc1ed6699e41caabe22a18d630))
* **deps:** bump golang.org/x/tools from 0.22.0 to 0.23.0 ([#454](https://github.com/bucketeer-io/code-refs/issues/454)) ([e6505a3](https://github.com/bucketeer-io/code-refs/commit/e6505a32aa5cc6a54fa8a4c4080c6428933dce3a))
* **deps:** bump golang.org/x/tools from 0.23.0 to 0.26.0 ([#462](https://github.com/bucketeer-io/code-refs/issues/462)) ([4aa30a8](https://github.com/bucketeer-io/code-refs/commit/4aa30a8942c98d5a956df82c041d3e68617e8631))
* **deps:** bump golang.org/x/tools from 0.26.0 to 0.27.0 ([#465](https://github.com/bucketeer-io/code-refs/issues/465)) ([7834a7f](https://github.com/bucketeer-io/code-refs/commit/7834a7f19a14a59870df7500f27366e8ba189566))
* **deps:** bump golang.org/x/tools from 0.27.0 to 0.28.0 ([#469](https://github.com/bucketeer-io/code-refs/issues/469)) ([00ea07e](https://github.com/bucketeer-io/code-refs/commit/00ea07eb1ba71487e44b20c45a3eacbdb1d0022f))
* **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 ([#439](https://github.com/bucketeer-io/code-refs/issues/439)) ([fa34367](https://github.com/bucketeer-io/code-refs/commit/fa34367010c79cfcd4e88d2d12cebd32364d74ca))
* reorganize search package ([#407](https://github.com/bucketeer-io/code-refs/issues/407)) ([5133ae5](https://github.com/bucketeer-io/code-refs/commit/5133ae5a54292869221309f33a6a655239aee7a0))
* update commit author for homebrew tap ([#380](https://github.com/bucketeer-io/code-refs/issues/380)) ([18d4e3f](https://github.com/bucketeer-io/code-refs/commit/18d4e3f281c4345205bc5d5ef5711150bf830a5f))
* update dockerfile to use latest alpine ([#389](https://github.com/bucketeer-io/code-refs/issues/389)) ([0a1af5a](https://github.com/bucketeer-io/code-refs/commit/0a1af5ae3ba636e757299d5f809fbb5610c6b5d2))
* update path to orb token ([#384](https://github.com/bucketeer-io/code-refs/issues/384)) ([4391894](https://github.com/bucketeer-io/code-refs/commit/43918949cc57176cd79c1f1e8c027966526e8469))
* update requests made to use 'state' filter; request deprecated flags ([#400](https://github.com/bucketeer-io/code-refs/issues/400)) ([a1e6e82](https://github.com/bucketeer-io/code-refs/commit/a1e6e82dd0f50f04050100fe8020afbcabbe659e))


### Build System

* bump to use go1.17 ([30d1948](https://github.com/bucketeer-io/code-refs/commit/30d1948c4e0f213e9b84380ad9d8c81ac830a41d))
* bump to use go1.17 ([c48c583](https://github.com/bucketeer-io/code-refs/commit/c48c583883316abf90245c892e2be3d4992d2477))
* **dockerfile:** update alpine to 3.14.1 ([#160](https://github.com/bucketeer-io/code-refs/issues/160)) ([31dd9c5](https://github.com/bucketeer-io/code-refs/commit/31dd9c5289ed51aeba919d71465ef1b6b154f0c7))
* update x/sys to support go 1.17 ([a5ee198](https://github.com/bucketeer-io/code-refs/commit/a5ee1981ae724836bd65934e4f6c5a2a235b2f2f))

## [0.0.1](https://github.com/bucketeer-io/code-refs/compare/v0.0.1...v0.0.1) (2025-03-06)

### Features

* bucketeer integration ([#1](https://github.com/bucketeer-io/code-refs/issues/1)) ([3926f23](https://github.com/bucketeer-io/code-refs/commit/3926f23ea27ecc8c63ed95e7da7cc14e6a01a507))


## [2.13.0] - 2024-12-18
### Added:
- `subdirectory` option to set path to `.launchdarkly/coderefs.yaml` config file if not located in the root, in order to support monorepo subdirectories.

### Changed:
- Updated dependencies

## [2.12.0] - 2024-03-28
### Added:
- Enable scanning of github workflow files [#441](https://github.com/launchdarkly/ld-find-code-refs/pull/441)

### Changed:
- Streamline HTTP client used in requests [#438](https://github.com/launchdarkly/ld-find-code-refs/pull/438)

## [2.11.10] - 2024-03-14
### Changed:
- Use same http client with rate-limited retries for ld api client calls [#437](https://github.com/launchdarkly/ld-find-code-refs/pull/437)

## [2.11.9] - 2024-03-04
### Changed:
- Respect rate-limit headers during retries [#434](https://github.com/launchdarkly/ld-find-code-refs/pull/434)

## [2.11.8] - 2024-02-13
### Changed:
- Updated dependencies

### Fixed:
- Multi-project scans now ignore projects without valid flags [#430](https://github.com/launchdarkly/ld-find-code-refs/pull/430)

## [2.11.7] - 2024-01-31
### Changed:
- Updated docker images to use [alpine 3.19.1](https://www.alpinelinux.org/posts/Alpine-3.19.1-released.html) ([#426](https://github.com/launchdarkly/ld-find-code-refs/pull/426))

## [2.11.6] - 2024-01-24
### Changed:
- Dependencies updated

### Fixed:
- Fixes index out of range error during extinction scanning [#425](https://github.com/launchdarkly/ld-find-code-refs/pull/425)

## [2.11.5] - 2024-01-08
### Added:
- Allow prune stage to be disabled [#405](https://github.com/launchdarkly/ld-find-code-refs/pull/405)

### Changed:
- Disable branch pruning by default in GitHub action [#405](https://github.com/launchdarkly/ld-find-code-refs/pull/405)
- Dependency updates

## [2.11.4] - 2023-10-16
### Changed:
- Dependencies updated

## [2.11.3] - 2023-09-20
### Changed:
- Dependencies updated

### Fixed:
- Docker image was not being built with correct alpine version (does not affect github or bitbucket integrations)

## [2.11.2] - 2023-08-21
### Fixed:
- Updates to release process

## [2.11.1] - 2023-08-18
### Changed:
- Update docker images

## [2.11.0] - 2023-08-10
### Changed:
- Update app to go 1.20
- Update dependencies

### Fixed:
- Update module-path to v2 #362

## [2.10.0] - 2023-02-21
### Changed:
- Performance improvements around searching for flag extinctions in commit diffs, including changing search to use Aho-Corasick algorithm to find flags.
- Updated dependencies

### Fixed:
- Error parsing git tag name when tag is annotated #329

## [2.9.2] - 2023-02-16
### Fixed:
- Bug introduced in 2.5.0 caused extinction scanning to run for hours and timeout.
- Typo in CircleCI orb caused failures

## [2.9.1] - 2023-01-27
### Fixed:
- CircleCI - pipeline exits with error when `$BASH_ENV` doesn't exist.

## [2.9.0] - 2023-01-26
Added:
* CircleCI orb - manually source `BASH_ENV` since the orb uses `sh` entrypoint instead of `bash`

Changed:
* Update dependencies

Fixed:
* Monorepo configuration should not require that a `dir` be specified for each project key, as described in the documentation

## [2.8.0] - 2022-11-10
Added:
* Support doublestar glob patterns in `filepattern` aliases

Changed:
* Update dependencies

## [2.7.0] - 2022-09-14
### Changed:
- Log a warning when a `filepattern` alias configuration does not match any files instead of failing
- Bumped dependencies

### Fixed:
- Error running code refs when using a newer API token with version `20220603` and later.

## [2.6.3] - 2022-09-09
### Changed:
- update documentation
- bumped dependencies

## [2.6.2] - 2022-09-06
### Fixed:
- error running CircleCI orb

## [2.6.1] - 2022-09-01
### Fixed:
- git repo permissions in GitHub Action

## [2.6.0] - 2022-08-31
### Changed:
- `gitlab` is a supported `repoType`
- Optional `defaultBranch` will fallback to `main` when not provided, instead of `master`
- Bumped dependencies
- Added debug logging

## [2.5.7] - 2022-03-01
### Changed:
- Update release configuration

## [2.5.6] - 2022-03-01
### Fixed:
- Change in release process lead to build with incorrect docker image tag

## [2.5.5] - 2022-03-01
### Fixed:
- Slice bounds out of range error when saving hunks (#224)

### Added:
- Enable builds for arm64 (#221)

## [2.5.4] - 2022-02-16
### Fixed
- Only a single flag per run was being searched for extinctions

### Added 
- `extinctions` command that will only generate and send extinctions using the `lookback` parameter

### Changed 
- Added additional examples for GitHub Action repo on how to configure the action

## [2.5.3] - 2022-02-16
### Fixed
- Only a single flag per run was being searched for extinctions

### Added 
- `extinctions` command that will only generate and send extinctions using the `lookback` parameter

### Changed 
- Added additional examples for GitHub Action repo on how to configure the action

## [2.5.0] - 2022-02-04
### Fixed:
- Snake case aliases we not being correctly generated due to bug in dependency.

### Changed:
- If new `projects` block is used with CSV output, the first project key is used in the output file name. If still using `projKey` there is no change.

### Added:
- Monorepo with starting directory support. More info can be read at [Projects](https://github.com/launchdarkly/ld-find-code-refs/blob/main/docs/CONFIGURATION.md#projects).

## [2.4.1] - 2021-12-17
### Fixed:
- Relative paths were not being expanded to an absolute path when used.

### Changed:
- Find Code References GitHub Action is moving to semver versioning. Previously it was a major version that was incremented on every release of the underlying command line tool. Now the GitHub Action version will mirror the command line tooling version. This is moving it from `v14` to `v2.4.1`

## [2.4.0] - 2021-11-22
### Changed:
- Performance improvements around searching for flags and aliases in the code base. Including changing search to use Aho-Corasick algorithm to find flags in file.
- If `--dryRun` is set, extinctions will not be attempted.

### Added:
- `--allowTags` which allows Code Refs to run against a tag instead of the branch.

### Fixed:
- Bug where alias filepattern's were not being validated. This meant `ld-find-code-refs` would run but if that file did not exist no aliases were generated.

## [2.3.0] - 2021-11-02
### Changed:
- Performance improvements around searching for flags and aliases in the code base.

## [2.2.4] - 2021-06-14
### Changed
- Matching flags with delimiters has been implemented in a more performant way.

## [2.2.3] - 2021-05-26
### Fixed
- File globbing for FilePattern alias support.

## [2.2.2] - 2021-04-27

### Added

- `repoName` is now a supported configuration option for [github action](https://docs.launchdarkly.com/home/code/github-actions#additional-configuration-options) and [bitbucket pipes](https://docs.launchdarkly.com/home/code/bitbucket#pipeline-configuration) 🎉. This is especially useful for a monorepo where multiple yaml configurations exist, each mapping to its own LD project key.

## [2.1.0]

### Added

- `ld-find-code-refs` will now scan for archived flags.

## [2.0.1] - 2020-10-05

### Fixed

- Fixes a bug causing `ld-find-code-refs` to scan non-regular files, like symlinks. Thanks @d3d-z7n!

## [2.0.0] - 2020-08-13

ℹ️ This release includes breaking changes to the command line tool. If you experience errors or unexpected behavior after upgrading, be sure to read these changelog notes carefully to make adjustments for any breaking changes.

### Added

- Most command line flags can now be [specified in a YAML file](https://github.com/launchdarkly/ld-find-code-refs/blob/main/docs/CONFIGURATION.md#yaml) located in the `.launchdarkly/coderefs.yaml` subdirectory of your repository. [docs](https://github.com/launchdarkly/ld-find-code-refs/blob/main/docs/CONFIGURATION.md#yaml)
  - The following options cannot be specified in YAML, and must be set using the command line or as environment variables:
    - `--dir` / `LD_DIR`
    - `--accessToken` / `LD_ACCESS_TOKEN`
- All command line flags can now be specified as environment variables. [docs](https://github.com/launchdarkly/ld-find-code-refs/blob/main/docs/CONFIGURATION.md#environment-variables)
- When flags with no code references are detected, `ld-find-code-refs` will search Git commit history to detect when the last reference to a feature flag was removed. Use the `--lookback` command line flag to configure the number of commits you would like to search. The lookback will start at the current commit and will review up to the last n commits to find the last reference of the flag. The default is 10 commits.
- Added support for scanning non-git repositories. Use the `--revision` flag to specify your repository version number.
- Added the `prune` sub-command to delete stale code reference data from LaunchDarkly manually by providing a list of branch names as arguments. example: `ld-find-code-refs prune [flags] "branch1" "branch2"`
- The GitHub actions wrapper now supports the `pull_request` event

### Fixed

- Exclude negations in `.ldignore` (lines beginning with an exclamation mark) now correctly include files.

### Changed

- Command line arguments names have been improved. Now, a flag specified with a single dash indicates a shorthand name, while 2 dashes indicate the longform name. Some existing configurations may be invalid, see `ld-find-code-refs --help` for details.
- The default delimiters (single quotes, double quotes and backticks) can now be disabled in the `coderefs.yaml` configuration. [docs](https://github.com/launchdarkly/ld-find-code-refs/blob/main/docs/CONFIGURATION.md#delimiters). Delimiters can no longer be specified using command line flags or environment variables. If you use additional delimiters, or would like to disable delimiters completely, use YAML configuration instead.

### Removed

- The `exclude` command-line option has been removed. Use the `.ldignore` file instead.
- `ld-find-code-refs` no longer requires the silver searcher (ag) as a runtime dependency.

## [1.5.1] - 2020-05-22

### Added

- Added support for specifying a custom default branch for the GitHub actions and Bitbucket pipes wrappers.

## [1.5.0] - 2020-05-11

### Added

- Added the ability to configure flag alias detection using a YAML configuration. See [the README](https://github.com/launchdarkly/ld-find-code-refs#configuring-aliases) for instructions.

### Fixed

- Improved logging around limitations.
- Fixed an edge case where false positives might be picked up for flag keys containing regular expression characters.

## [1.4.0] - 2020-03-16

### Added

- Added a `--ignoreServiceErrors` option to the CLI. If enabled, the scanner will terminate with exit code 0 when the LaunchDarkly API is unreachable or returns an unexpected response.

### Changed

- ld-find-code-refs now requires go1.13 to build.

## [1.3.1] - 2019-09-24

### Fixed

- Fixed a regression causing no references to be found when a relative path is supplied to `dir`

## [1.3.0] - 2019-09-19

### Added

- Added a `--outDir` option to the CLI. If provided, code references will be written to a csv file in `outDir`.
- Added a `--dryRun` option to the CLI. If provided, `ld-find-code-refs` will scan for code references without sending them to LaunchDarkly. May be used in conjunction with `--outDir` to output code references data to a csv file instead of sending data to LaunchDarkly.

### Fixed

- `ld-find-code-refs` now supports scanning repositories with a large number of flags using a pagination strategy. Thanks @cuzzasoft!
- Delimiters will now always be respected when searching for flags referenced in code. This fixes a bug causing references for certain flag keys to match against other flag keys that are substrings of the matched reference.

## [1.2.0] - 2019-08-13

### Added

- Added a `--branch` option to the CLI. This lets a branch name be manually specified when the repo is in a detached head state.
- GitHub actions v2 support: the github actions wrapper reads the branch name from `GITHUB_REF` and populates the `branch` option with it.

## [1.1.1] - 2019-04-11

### Fixed

- `ld-find-code-refs` will no longer exit with a fatal error when Git credentials have not been configured (required for branch cleanup). Instead, a warning will be logged.

## [1.1.0] - 2019-04-11

### Added

- `ld-find-code-refs` will now remove branches that no longer exist in the git remote from LaunchDarkly.

## [1.0.1] - 2019-03-12

### Changed

- Fixed a potential bug causing `.ldignore` paths to not be detected in some environments.
- When `.ldignore` is found, a debug message is logged.

## [1.0.0] - 2019-02-21

Official release

## [0.7.0] - 2019-02-15

### Added

- Added support for Windows. `ld-find-code-refs` releases will now contain a windows executable.
- Added a new option `-delimiters` (`-D` for short), which may be specified multiple times to specify delimiters used to match flag keys.

### Fixed

- The `dir` command line option was marked as optional, but is actually required. `ld-find-code-refs` will now recognize this option as required.
- `ld-find-code-refs` was performing extra steps to ignore directories for files in directories matched by patterns in `.ldignore`. This ignore process has been streamlined directly into the search so files in `.ldignore` are never scanned.

### Changed

- The command-line [docker image](https://hub.docker.com/r/launchdarkly/ld-find-code-refs) now specifies `ld-find-code-refs` as the entrypoint. See our [documentation](https://github.com/launchdarkly/ld-find-code-refs#docker) for instructions on running `ld-find-code-refs` via docker.
- `ld-find-code-refs` will now only match flag keys delimited by single-quotes, double-quotes, or backticks by default. To add more delimiters, use the `delimiters` command line option.

## [0.6.0] - 2019-02-11

### Added

- Added a new command line argument, `version`. If provided, the current `ld-find-code-refs` version number will be logged, and the scanner will exit with a return code of 0.
- The `debug` option is now available to the CircleCI orb.
- Added support for parsing `.ldignore` files specified in the root directory of the scanned repository. `.ldignore` may be used to specify a pattern (compatible with the `.gitignore` spec: https://git-scm.com/docs/gitignore#_pattern_format) for files to exclude from scanning.

### Changed

- The internal API for specifying the default git branch (`defaultBranch`) has been changed. The `defaultBranch` argument on earlier versions of `ld-find-code-refs` will no longer do anything.

### Fixed

- `ld-find-code-refs` will no longer error out if an unknown error occurs when scanning for code reference hunks within a file. Instead, an error will be logged.

## [0.5.0] - 2019-02-01

### Master

- Added support for parsing `.ldignore` files specified in the root directory of the scanned repository. `.ldignore` may be used to specify a pattern (compatible with the `.gitignore` spec: https://git-scm.com/docs/gitignore#_pattern_format) for files to exclude from scanning.

### Added

- Generate deb and rpm packages when releasing artifacts.

### Changed

- Automate Homebrew releases
- Added word boundaries to flag key regexes.
  - This should reduce false positives. E.g. for flag key `cool-feature` we will no longer match `verycool-features`.

## [0.4.0] - 2019-01-30

### Added

- Added support for relative paths to CLI `-dir` parameter.
- Added a new command line argument, `debug`, which enables verbose debug logging.
- `ld-find-code-refs` will now exit early if required dependencies are not installed on the system PATH.

### Changed

- Renamed `parse` package to `coderefs`. The `Parse()` method in the aformentioned package is now `Scan()`.

### Fixed

- `ld-find-code-refs` will no longer erroneously make PATCH API requests to LaunchDarkly when url template parameters have not been configured.

## [0.3.0] - 2019-01-23

### Added

- Added openssh as a dependency for the command-line docker image.

### Changed

- The default for `contextLines` is now 2. To disable sending source code to LaunchDarkly, set the `contextLines` argument to `-1`.
- Improved logging to provide more detailed summaries of actions performed by the scanner.

### Fixed

- Fixed a bug in the CircleCI orb config causing `contextLines` to be a string parameter, instead of an integer.

### Removed

- Removed the `repoHead` parameter. `ld-find-code-refs` now only supports scanning repositories already checked out to the desired branch.
- Removed an unnecessary dependency on openssh in Dockerfiles.

## [0.2.1] - 2019-01-17

### Fixed

- Fix a bug causing an error to be returned when a repository connection to LaunchDarkly does not initially exist on execution.

### Removed

- Removed the `cloneEndpoint` command line argument. `ld-find-code-refs` now only supports scanning existing repository clones.

## [0.2.0] - 2019-01-16

### Fixed

- Use case-sensitive `ag` search so we don't get false positives that look like flag keys but have different casing.

### Changed

- This project has been renamed to `ld-find-code-refs`.
- Logging has been overhauled.
- Project layout has been updated to comply with https://github.com/golang-standards/project-layout.
- `updateSequenceId` is now an optional parameter. If not provided, data will always be updated. If provided, data will only be updated if the existing `updateSequenceId` is less than the new `updateSequenceId`.
- Payload limits have been implemented
  - Flags with keys shorter than 3 characters are no longer supported.
  - Lines are truncated after 500 characters.
  - Search is terminated after 5,000 files are matched.
  - Search is terminated after 5,000 hunks are generated.
  - Number of hunks per file is limited to 1,000.
  - A file can only have 500 hunked lines per flag.
- Use `launchdarkly` docker hub namespace instead of `ldactions`.

## [0.1.0] - 2019-01-02

### Changed

- `pushTime` CLI arg renamed to `updateSequenceId`. Its type has been changed from timestamp to integer.
  - Note: this is not considered a breaking change as the CLI args are still in flux. After the 1.0 release arg changes will be considered breaking.

### Fixed

- Upserting repos no longer fails on non-existent repos

## [0.0.1] - 2018-12-14

### Added

- Automated release pipeline for github releases and docker images
- Changelog
