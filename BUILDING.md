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
