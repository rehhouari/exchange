# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2022-10-20

### Added

- Added `func (exchange *Exchange) SetContext(context context.Context)` to allow passing of context to
  `net/http` package.
- Some tests added.

### Fixed

- Fixed validation of HTTP request query parameters.

## [1.0.0] - 2021-10-02

- Initial release.
