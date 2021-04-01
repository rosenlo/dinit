# Changes
## 1.3.0 (2020-12-22)
md5sum: 648624b3a21a6b0a4f12017f6ca1efa4

### Added
- [dinit] Complete lifecycle policy for services

### Changed
- [dinit] Change the log formatter to JSONFormatter from TextFormatter

### Removed

- [dinit] Synchronization time for `consul` ( using `post_stop` phase instead of `DINIT_CONSUL_LEAVE_INTERVAL` )

## 1.2.1 (2020-12-18)
md5sum: c726d41cf001f1b7b42717cbc5ccd561

### Changed
- [dinit] reduce encumbrance logs

### Fixed
- [dinit] avoid repeated shutdown procedures
- [dinit] avold using invalid code to exit when the child process receivers a siganl

## 1.2.0 (2020-12-16)
md5sum: 853dfb29ad42e1bbcf4cdd348e3eaa95

### Added
- [dinit] report information to CMDB
- [dinit] offline before consul leave

## 1.1.3 (2020-12-14)
md5sum: 7e71537a23338ee090792fa7801beefe

### Changed
- [dinit] add retry for flowController
- [dinit] `DINIT_DIRECTOR_ADDR` is requirement

## 1.1.2 (2020-12-08)
md5sum: afa78edac263317521794f4ec5f13828

### Added
- [dinit] support for extra environment variable into log fields

## 1.1.1 (2020-11-19)
md5sum: caac97cb053312177df324458b3473a4

### Fixed
- [dinit] fixed maybe cause panic with nil pointer

## 1.1.0 (2020-11-12)

md5sum: 2a99ef48d6436e2a400e49ee83703f42

## 1.1.0-alpha (2020-10-27)

md5sum: d815ac979a3e2a88b443dc4a981e93ce

### Added
- [dinit] Add version information to log
- [dinit] Support service restart in place

### Changed
- [dinit] Collect logs of child process stderr
- [dinit] Add synchronization time to `consul`

## 1.0.8 (2020-09-17)

### Fixed
- [dinit] fix the missing system environment

## 1.0.7 (2020-09-09)

### Changed
- [dinit] Support for environment

## 1.0.6 (2020-08-26)

### Changed
- [dinit] Change the log output format to json format

### Fixed
- [dinit] fix the `pre_stop` pharse did not complete exit


## 1.0.5 (2020-08-14)

### Changed
- [dinit] add log output file
- [dinit] the main process hanging when the child process return non-zero exit code

## 1.0.4 (2020-07-02)

### Changed
- [dinit] - main process exit code use sub-process exit code
- [dinit] - support main process hang through environment variable of `DINIT_EXIT`
- [dinit] - flow control, stop traffic before stopping services

## 1.0.3 (2020-06-28)

### Fixed
- [watchman] - fix synchronization of new apps


## 1.0.2 (2020-06-19)

### Added
- [watchman] - support running outside of k8s

### Changed
- [dinit] - shutdown the ingress before exiting
- [dinit] - ignore the consul check for JS applications


## 1.0.1 (2020-06-18)

### Added
- [dinit] - flow control
- [watchman] - New component, a service for synchronizing the health of the
  applications between internal and external

## 1.0.0-beta (2020-04-10)

brand new version

### Feature

- Dinit as the main process in the container
- Processes lifecycle
- Service startup dependencies
