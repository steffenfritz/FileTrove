# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

FileTrove is a digital forensics and archival metadata indexing tool (Go, AGPL v3). It walks directory trees, identifies files by type using PRONOM/siegfried, computes cryptographic hashes, extracts EXIF/xattr metadata, scans with YARA-X rules, and writes all results into a SQLite database with TSV export support.

## Build & Test Commands

**Requires CGO** due to YARA-X C library integration.

Run from `cmd/ftrove/`:
```sh
task build          # Build for current platform (PIE mode, version injection)
task changetest     # Full integration test: unit tests + build + install + scan + inspect DB
task clean          # Remove build artifacts and test DB
```

Unit tests only (from repo root):
```sh
go test -v ./...
```

First-time dev setup (installs deps and builds YARA-X C library):
```sh
task setup_mac          # macOS (uses Homebrew)
task setup_linux        # Debian/Ubuntu (uses apt + rustup)
```

Cross-compile (requires Zig):
```sh
task build-linux-amd    # Linux amd64 via Zig CC
task build-windows-amd  # Windows amd64 via Zig CC
```

Admin tool (`cmd/admftrove/`):
```sh
make build      # Single platform
make build_all  # darwin + linux + windows
```

## Architecture

### Package Layout

- **Root package** (`github.com/steffenfritz/FileTrove`) — core library used by both CLI tools
- **`cmd/ftrove/`** — main CLI; orchestrates scanning sessions
- **`cmd/admftrove/`** — admin utility for NSRL hash database management

### Core Library Modules

| File | Responsibility |
|------|---------------|
| `db.go` | SQLite schema, session/file/dir CRUD, prepared statements |
| `hash.go` | Single-pass multi-algorithm hashing (MD5, SHA1, SHA256, SHA512, BLAKE2B-512) |
| `siegfried.go` | File type identification via PRONOM registry |
| `yara.go` | YARA-X rule compilation and scanning (CGO) |
| `nsrl.go` | BoltDB-backed NSRL SHA1 lookup and import |
| `exif.go` | EXIF metadata extraction from images |
| `xattr.go` | Filesystem extended attribute reading |
| `entropy.go` | Shannon entropy calculation (files up to 1GB) |
| `times.go` | File access/change/birth timestamps |
| `filewalk.go` | Directory traversal via `filepath.WalkDir` |
| `dublincore.go` | Dublin Core Elements session metadata from JSON |

### Data Flow

1. **Install phase**: `ftrove --install <dir>` creates the SQLite DB, downloads siegfried signatures and optionally the 4GB NSRL BoltDB.
2. **Scan phase**: `ftrove -i <input-dir>` — walks the tree, runs each file through identification → hashing → entropy → EXIF/xattr/YARA → NSRL lookup → DB insert.
3. **Export phase**: `ftrove -e` writes TSV files per table for the active session.

### Key Data Structures (`db.go`)

- `SessionMD` — UUID-based session with timestamps, flags for enabled modules, version strings
- `FileMD` — per-file record: all hashes, PRONOM/MIME/format info, timestamps, entropy, NSRL hit flag
- `DirMD` — per-directory record with timestamps

### Database Schema

SQLite database (`db/filetrove.db`) with tables: `filetrove` (version), `sessionsmd`, `files`, `directories`, `exif`, `yara`, `xattr`, `ntfsads`, `dublincore`. Full schema in `database_schema.dbml`.

### YARA-X Integration

YARA-X requires a C library built separately from source. The CI workflow (`.github/workflows/buildstatus.yml`) builds it via `cargo-c`. For local development, see `BUILDING.md` for platform-specific setup instructions.

## CLI Usage Reference

```sh
ftrove --install <dir>        # Initialize install directory
ftrove -i <input-dir>         # Run a scan
ftrove -l                     # List sessions
ftrove -e                     # Export active session to TSV
ftrove -d <dublincore.json>   # Attach Dublin Core metadata to scan
ftrove -p "<project note>"    # Add project description
ftrove -z <timezone>          # Set timezone (e.g. Europe/Berlin)
ftrove -y <yara-rules-dir>    # Enable YARA scanning
```

Example YARA rule file: `testdata/yara/`.
Example Dublin Core JSON: `testdata/dublincore_ex.json`.
