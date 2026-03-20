<p align="center">
<img src="https://github.com/steffenfritz/FileTrove/assets/16431534/b8c1456d-08bb-48bb-afcf-5e99db8466b9" width="300">
</p>

<p align="center">
  <img alt="Build Status" src="https://github.com/steffenfritz/FileTrove/actions/workflows/buildstatus.yml/badge.svg">
  <a href="https://www.gnu.org/licenses/agpl-3.0"><img alt="License: AGPL v3" src="https://img.shields.io/badge/License-AGPL_v3-blue.svg"></a>
  <a href="https://pkg.go.dev/github.com/steffenfritz/FileTrove"><img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/steffenfritz/FileTrove.svg"></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/steffenfritz/FileTrove"><img alt="OpenSSF Scorecard" src="https://api.scorecard.dev/projects/github.com/steffenfritz/FileTrove/badge"></a>
  <a href="https://www.bestpractices.dev/projects/8952"><img alt="OpenSSF Best Practices" src="https://www.bestpractices.dev/projects/8952/badge"></a>
</p>

**VERSION: v1.0.0-BETA-4**

---

FileTrove walks a directory tree, identifies every file, computes metadata, and writes all results into a SQLite database with TSV export support.

## What it collects

| Category | Details |
|----------|---------|
| **File type** | MIME type, [PRONOM](https://www.nationalarchives.gov.uk/PRONOM/) identifier, format version, identification proof/note, extension — via [siegfried](https://github.com/richardlehane/siegfried) |
| **File & directory timestamps** | Creation, modification, and access times |
| **Hashes** | MD5, SHA1, SHA256, SHA512, BLAKE2B-512 |
| **Entropy** | Shannon entropy (files up to 1 GB) |
| **Extended attributes** | xattr from ext3/ext4, btrfs, APFS, and others |
| **EXIF metadata** | Extracted from image files |
| **YARA-X** | Match results from your own rule files |
| **NSRL** | Flags known software files via the [National Software Reference Library](https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl) |
| **Dublin Core** | Optional session-level descriptive metadata |

Each file and directory gets a UUIDv4 as a unique identifier. All results land in a SQLite database and can be exported to TSV.

## Installation

1. **Get the binary** — download a release from the [releases page](https://github.com/steffenfritz/FileTrove/releases), or compile from source (see [BUILDING.md](BUILDING.md)). Both a standard dynamic binary (`ftrove`) and a static binary (e.g. `ftrove_amd64_linux_static`) are provided.

2. **Run the installer** from the directory where you want FileTrove to live:
   ```sh
   ./ftrove --install .
   ```
   This creates a `db/` directory and downloads the siegfried signature database.

3. **Set up the NSRL database.** FileTrove uses a Bloom filter (`db/nsrl.bloom`) for NSRL lookups. If the repository was cloned with Git LFS, the file is already in `db/`. Otherwise, build it from upstream NIST data (requires `sqlite3`, `curl`, `unzip`):
   ```sh
   task nsrl:build-all       # All subsets including legacy (~80-110 MB, recommended)
   task nsrl:build-mobile    # Modern + Android + iOS (~50-65 MB)
   task nsrl:build-modern    # Modern OS software only (~30-45 MB)
   ```
   The default (`build-all`) is recommended for archival and digital preservation work, where legacy software from older systems is commonly encountered.

   > **Note:** Building from source requires significant temporary disk space (~190 GB for modern, ~350 GB for all) because the NSRL RDS SQLite databases are large. Run `task nsrl:clean` afterwards to reclaim the space. If you cloned with Git LFS, `db/nsrl.bloom` is already present and no build is needed.

4. **You're ready.**

### YARA-X

YARA-X scanning requires a C library that is not bundled with FileTrove. It is built automatically during `task build` if not already present. See [BUILDING.md](BUILDING.md) for setup instructions.

- Example rule files: `testdata/yara/`
- When a rule matches, the rule name, session UUID, and file UUID are recorded in the `yara` table. The rule file itself is not stored.

### NSRL

FileTrove ships a pre-built NSRL Bloom filter via Git LFS. When NIST publishes a new RDS version, rebuild by updating `NSRL_VERSION` in `Taskfile.nsrl.yml` and running one of the build targets above.

You can also build a custom Bloom filter from any newline-delimited list of SHA1 hashes:

```sh
admftrove --creatensrl hashes.txt --nsrlversion "my-hashset-v1"
```

Optional flags: `--nsrl-estimate` (expected hash count, default 40M) and `--nsrl-fpr` (false positive rate, default 0.0001). Copy the resulting `nsrl.bloom` into `db/`.

## Running a scan

```sh
./ftrove -i $DIRECTORY
```

FileTrove walks `$DIRECTORY` recursively. Run `./ftrove -h` for all available flags.

## Viewing results

List all sessions and export one to TSV:

```sh
./ftrove -l
./ftrove -t 926be141-ab75-4106-8236-34edfcf102f2
```

You can also query the SQLite database directly:

- **CLI:** `sqlite3 db/filetrove.db`
- **GUI:** [sqlitebrowser](https://sqlitebrowser.org/)
- **Visualisation:** [Sqliteviz](https://sqliteviz.com/app/#/)

## Background

FileTrove is the successor of [filedriller](https://github.com/steffenfritz/filedriller), based on the iPres 2021 paper [Marrying siegfried and the National Software Reference Library](https://phaidra.univie.ac.at/detail/o:1424904).
