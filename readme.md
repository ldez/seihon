# Seihon

[![GitHub release](https://img.shields.io/github/release/ldez/seihon.svg)](https://github.com/ldez/seihon/releases/latest)
[![Build Status](https://github.com/ldez/seihon/workflows/Main/badge.svg?branch=master)](https://github.com/ldez/seihon/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/seihon)](https://goreportcard.com/report/github.com/ldez/seihon)

A simple tool to publish multi-arch images on the Docker Hub.

If you appreciate this project:

[![Sponsor](https://img.shields.io/badge/Sponsor%20me-%E2%9D%A4%EF%B8%8F-pink)](https://github.com/sponsors/ldez)

![image](docs/img.png)

## Usage

- [seihon](docs/seihon.md)
- [seihon publish](docs/seihon_publish.md)

## Installation

### Download / CI Integration

```bash
curl -sfL https://raw.githubusercontent.com/ldez/seihon/master/godownloader.sh | bash -s -- -b $GOPATH/bin v0.5.1
```

<!--
To generate the script:

```bash
godownloader --repo=ldez/seihon -o godownloader.sh

# or

godownloader --repo=ldez/seihon > godownloader.sh
```
-->

### From a package manager

- [ArchLinux (AUR)](https://aur.archlinux.org/packages/seihon/):
```bash
yay -S seihon
```

- [Homebrew Taps](https://github.com/ldez/homebrew-tap)
```bash
brew tap ldez/tap
brew update
brew install seihon
```

### From Binaries

You can use pre-compiled binaries:

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/seihon/releases/)
* Unzip the archive.
* Add `seihon` in your `PATH`.

## Tips

- GitHub Actions:

```yaml
name: Example

# ...

jobs:

  main:
    # ...
    env:
      # ...
      SEIHON_VERSION: v0.7.1

    steps:
      # ...
      
      # Install Docker image multi-arch builder
      - name: Install Seihon ${{ env.SEIHON_VERSION }}
        #if: startsWith(github.ref, 'refs/tags/v')
        run: |
          curl -sSfL https://raw.githubusercontent.com/ldez/seihon/master/godownloader.sh | sh -s -- -b $(go env GOPATH)/bin ${SEIHON_VERSION}
          seihon --version

      - name: Publish Docker Images (Seihon)
        #if: startsWith(github.ref, 'refs/tags/v')
        run: make publish-images
```

- Travis CI:

```yaml
before_deploy:
  # Install Docker image multi-arch builder
  - curl -sfL https://raw.githubusercontent.com/ldez/seihon/master/godownloader.sh | bash -s -- -b $(go env GOPATH)/bin ${SEIHON_VERSION}
  - seihon --version
  # Add QEMU only for some specific cases.
  - docker run --rm --privileged hypriot/qemu-register

deploy:
  - provider: script
    skip_cleanup: true
    script: seihon publish <your configuration>
    on:
      tags: true
```
