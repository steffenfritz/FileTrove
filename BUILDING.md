# Building FileTrove from Source

FileTrove requires CGO because it links against the [YARA-X](https://virustotal.github.io/yara-x/) C library.

---

## Supported Platforms

| Platform | Architecture | Task | Minimum OS | Notes |
|----------|-------------|------|-----------|-------|
| Linux | x86\_64 | `build` / `build-linux-amd` | Kernel 3.2+ (musl static) | Static binary; runs on Ubuntu 16.04+, Debian 9+, RHEL 7+ |
| Linux | arm64 | `build-linux-arm` | Kernel 3.2+ (musl static) | ARM servers, Raspberry Pi 4+ (64-bit OS) |
| macOS | x86\_64 (Intel) | `build` / `build-mac-amd` | macOS 12 Monterey (2021) | Intel Macs; cross-compile from Apple Silicon |
| macOS | arm64 (Apple Silicon) | `build` / `build-mac-arm` | macOS 12 Monterey (2021) | M1/M2/M3+ chips; cross-compile from Intel |
| Windows | x86\_64 | `build` / `build-windows-amd` | Windows 10 / Server 2016 | Dynamic build requires `yara_x_capi.dll` at runtime |

**Key notes:**
- Linux and macOS binaries are built as static (musl-linked) by the cross-compile tasks, so they carry no runtime library dependencies and run on the oldest kernels listed above. The `task build` native build produces a dynamic binary whose minimum glibc version matches the build host.
- Windows does not support fully static linking of the YARA-X C library; the DLL ships alongside the binary.
- Go 1.21+ dropped support for Windows 7/8 and macOS < 12. FileTrove inherits these lower bounds.

---

## Prerequisites

All build paths require:

- [Go](https://go.dev/) ≥ 1.26
- [Rust + Cargo](https://www.rust-lang.org/tools/install) (stable toolchain)
- [cargo-c](https://github.com/lu-zero/cargo-c) (`cargo install cargo-c`)
- [Task](https://taskfile.dev/) v3
- [Zig](https://ziglang.org/) — required only for cross-compilation

---

## Native Build (macOS)

### 1. Install Homebrew

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install Task

```sh
brew install go-task/tap/go-task
```

### 3. Install dependencies and build YARA-X

```sh
cd cmd/ftrove
task setup_mac
```

This installs Go, Rust, Zig, SQLite, and pkg-config via Homebrew, builds the YARA-X C library into `$HOME/yara_install`, and writes CGO environment variables to `cmd/ftrove/.yara_env` (auto-loaded by subsequent `task` invocations).

### 4. Build

```sh
task build
```

---

## Native Build (Debian / Ubuntu)

### 1. Install Task

```sh
curl -fsSL https://taskfile.dev/install.sh | sh -s -- -b ~/.local/bin
```

Ensure `~/.local/bin` is on your `PATH`, then verify with `task --version`.

### 2. Install dependencies and build YARA-X

```sh
cd cmd/ftrove
task setup_linux
```

This installs build tools, pkg-config, sqlite3, and Go via apt, installs Rust via rustup and `cargo-c`, builds the YARA-X C library into `$HOME/yara_install`, and writes CGO environment variables to `cmd/ftrove/.yara_env`.

### 3. Build

```sh
task build
```

---

## Building a Distribution Bundle

To create a self-contained bundle ready for end users:

```sh
task dist:bundle
```

This builds both binaries, downloads `siegfried.sig`, and packages everything into `build/<os>_<arch>/`. The resulting folder can be used immediately:

```sh
cd build/darwin_arm64         # or linux_amd64, etc.
./ftrove --install .          # creates filetrove.db, logs/, downloads nsrl-<variant>.bloom
./ftrove -i /path/to/files
```

The NSRL bloom filter is downloaded automatically during `--install` from the GitHub Releases page. Use `--nsrl-variant` to select which subset (default: `all`):

```sh
./ftrove --install . --nsrl-variant modern   # ~150 MB, modern OS software only
./ftrove --install . --nsrl-variant mobile   # ~200 MB, modern + Android + iOS
./ftrove --install . --nsrl-variant all      # ~240 MB, all subsets including legacy
```

To create a `.tar.gz` archive for distribution:

```sh
task dist:archive
```

---

## Building the NSRL Bloom filter

The NSRL bloom filter is not bundled in the repository. For most users it is downloaded automatically during `ftrove --install`. If you need to build it locally (e.g. to publish a new release asset, or to use a different FPR), use the Taskfile targets. Requirements: `sqlite3`, `curl`, `unzip`, and a built `admftrove`.

```sh
task nsrl:build-all       # All subsets including legacy (~240 MB, recommended)
task nsrl:build-mobile    # Modern + Android + iOS (~200 MB)
task nsrl:build-modern    # Modern OS software only (~150 MB)
```

For archival and digital preservation work, `build-all` is recommended since legacy software is commonly found on older media and disk images.

> **Disk space warning:** The build tasks download and extract the NSRL RDS SQLite databases temporarily. The temporary files are stored in `tmp/nsrl/` and can be removed after the build with `task nsrl:clean`.

| Build target | NSRL subsets | Download | Extracted | Total disk needed | Distinct hashes | Bloom size (1% FPR) |
|---|---|---|---|---|---|---|
| `build-modern` | Modern | ~18 GB | ~169 GB | **~190 GB** | ~31M | ~150 MB |
| `build-mobile` | Modern + Android + iOS | ~29 GB | ~242 GB | **~275 GB** | ~81M | ~200 MB |
| `build-all` | Modern + Android + iOS + Legacy | ~40 GB | ~305 GB | **~350 GB** | ~87M | ~240 MB |

Check whether your local bloom file matches the configured upstream version:

```sh
task nsrl:check
```

To update to a new NIST release: bump `NSRL_VERSION` in `Taskfile.nsrl.yml` to the latest [NIST RDS release](https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl/nsrl-download/current-rds), run `task nsrl:clean`, rebuild all three variants, upload them to a new GitHub Release tag, and update the URL constants in `install.go`.

### Publishing bloom files to GitHub Releases (maintainers)

The asset filenames must match `nsrl-modern.bloom`, `nsrl-mobile.bloom`, `nsrl-all.bloom` exactly — these are the filenames `install.go` downloads from `NSRLBloomURL{Modern,Mobile,All}`.

1. Bump `NSRL_VERSION` in `Taskfile.nsrl.yml` and run a fresh build of all three variants:
   ```sh
   task nsrl:clean
   task nsrl:build-modern
   task nsrl:build-mobile
   task nsrl:build-all
   ```
2. Create the release with tag `nsrl-<NSRL_VERSION>` (e.g. `nsrl-2026.03.1`) and attach the three bloom files:
   ```sh
   VERSION=$(grep '^  NSRL_VERSION:' Taskfile.nsrl.yml | awk -F'"' '{print $2}')
   gh release create "nsrl-${VERSION}" \
     --title "NSRL ${VERSION}" \
     --notes "NSRL RDS ${VERSION} bloom filters (FPR 1%)." \
     db/nsrl-modern.bloom db/nsrl-mobile.bloom db/nsrl-all.bloom
   ```
3. Update the three `NSRLBloomURL*` constants in `install.go` to point at the new tag, commit, and open a PR.
4. Verify a fresh install picks up the new files:
   ```sh
   ./ftrove --install /tmp/ft-test --nsrl-variant all
   ```

---

## Running tests

From the repository root:

```sh
go test -v ./...
```

Full integration test (build + install + scan + DB inspection):

```sh
cd cmd/ftrove
task changetest
```
