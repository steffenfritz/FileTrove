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

## Cross-Compilation

All cross-compile tasks use [Zig](https://ziglang.org/) as the C cross-compiler (`brew install zig` / `apt install zig`). They produce a single dynamic binary and are **not** used by the CI release runners — intended for local developer use only.

Run from `cmd/ftrove/` (or `cmd/admftrove/` for the admin tool):

| Target | Task | Output |
|--------|------|--------|
| Linux x86\_64 | `task build-linux-amd` | `ftrove_linux` |
| Linux arm64 | `task build-linux-arm` | `ftrove_arm64_linux` |
| macOS x86\_64 | `task build-mac-amd` | `ftrove_amd64_darwin` |
| macOS arm64 | `task build-mac-arm` | `ftrove_arm64_darwin` |
| Windows x86\_64 | `task build-windows-amd` | `ftrove.exe` |

The same task names are available in `cmd/admftrove/` and produce equivalently named `admftrove_*` binaries.

---

## Running Tests

From the repository root:

```sh
go test -v ./...
```

Full integration test (build + install + scan + DB inspection):

```sh
cd cmd/ftrove
task changetest
```
