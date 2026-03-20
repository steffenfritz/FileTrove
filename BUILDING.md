# Building FileTrove from Source

FileTrove requires CGO because it links against the [YARA-X](https://virustotal.github.io/yara-x/) C library. This guide walks you through setting up a complete development environment on **macOS** and **Debian/Ubuntu**.

---

## macOS

### 1. Install Homebrew

If you don't have Homebrew yet:

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install the Task build tool

```sh
brew install go-task/tap/go-task
```

### 3. Install remaining dependencies and build YARA-X

From the repository root, run:

```sh
cd cmd/ftrove
task setup_mac
```

This will:
- Install Go, Rust, Zig, SQLite, and pkg-config via Homebrew
- Install `cargo-c`
- Clone YARA-X v1.14.0 into `$HOME/yara-x` and build its C library into `$HOME/yara_install`
- Write the required CGO environment variables to `cmd/ftrove/.yara_env` (auto-loaded by subsequent `task` invocations)

### 4. Build

```sh
task build
```

---

## Debian / Ubuntu (latest stable or LTS)

### 1. Install the Task build tool

```sh
curl -fsSL https://taskfile.dev/install.sh | sh -s -- -b ~/.local/bin
```

Ensure `~/.local/bin` is on your `PATH`, then verify with `task --version`.

### 2. Install remaining dependencies and build YARA-X

From the repository root, run:

```sh
cd cmd/ftrove
task setup_linux
```

This will:
- Install build tools, pkg-config, sqlite3, and Go via apt (Go 1.21+ supports `GOTOOLCHAIN=auto`, so the version in `go.mod` is fetched automatically)
- Install Rust via rustup and `cargo-c`
- Clone YARA-X v1.14.0 into `$HOME/yara-x` and build its C library into `$HOME/yara_install`
- Write the required CGO environment variables to `cmd/ftrove/.yara_env` (auto-loaded by subsequent `task` invocations)

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

This builds both binaries, downloads `siegfried.sig`, copies `nsrl.bloom`, and packages everything into `build/<os>_<arch>/`. The resulting folder can be used immediately:

```sh
cd build/darwin_arm64         # or linux_amd64, etc.
./ftrove --install .          # creates filetrove.db and logs/
./ftrove -i /path/to/files
```

To create a `.tar.gz` archive for distribution:

```sh
task dist:archive
```

> **Prerequisite:** `db/nsrl.bloom` must exist in the repository root before running `task dist:bundle`. Build it with `task nsrl:build-all` if it doesn't exist yet. See below.

---

## Building the NSRL Bloom filter

FileTrove uses a Bloom filter for NSRL lookups. The pre-built `db/nsrl.bloom` is included in the repository. To rebuild it from upstream NIST data (requires `sqlite3`, `curl`, `unzip`, and a built `admftrove`):

```sh
task nsrl:build-all       # All subsets including legacy (~80-110 MB, recommended)
task nsrl:build-mobile    # Modern + Android + iOS (~50-65 MB)
task nsrl:build-modern    # Modern OS software only (~30-45 MB)
```

For archival and digital preservation work, `build-all` is recommended since legacy software is commonly found on older media and disk images.

> **Disk space warning:** The build tasks download and extract the NSRL RDS SQLite databases temporarily. The temporary files are stored in `tmp/nsrl/` and can be removed after the build with `task nsrl:clean`. The resulting `nsrl.bloom` file is only 30-110 MB.

| Build target | NSRL subsets | Download | Extracted | Total disk needed | Distinct hashes |
|---|---|---|---|---|---|
| `build-modern` | Modern | ~18 GB | ~169 GB | **~190 GB** | ~31M |
| `build-mobile` | Modern + Android + iOS | ~29 GB | ~242 GB | **~275 GB** | ~81M |
| `build-all` | Modern + Android + iOS + Legacy | ~40 GB | ~305 GB | **~350 GB** | ~87M |

Check whether your local bloom file matches the configured upstream version:

```sh
task nsrl:check
```

To update, bump `NSRL_VERSION` in `Taskfile.nsrl.yml` to the latest [NIST RDS release](https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl/nsrl-download/current-rds), run `task nsrl:clean`, and rebuild.

---

## Running tests

From the repository root:

```sh
go test -v ./...
```

For a full integration test (build + install + scan + DB inspection):

```sh
cd cmd/ftrove
task changetest
```
