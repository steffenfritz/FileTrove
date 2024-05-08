<p align="center">
<img src="https://github.com/steffenfritz/FileTrove/assets/16431534/b8c1456d-08bb-48bb-afcf-5e99db8466b9" width="350">
</p>


![Build Status](https://github.com/steffenfritz/FileTrove/actions/workflows/buildstatus.yml/badge.svg)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/steffenfritz/FileTrove.svg)](https://pkg.go.dev/github.com/steffenfritz/FileTrove)

STATUS:  Development

VERSION: v1.0.0-DEV-16

NOTE: Release DEV-16 is the last DEV release. The next one will be a feature freeze BETA version. (2024-04-23)

## About

FileTrove indexes files and creates metadata from them.

The single binary application walks a directory tree and identifies all regular files by type with [siegfried](https://github.com/richardlehane/siegfried), giving you the 

* MIME type
* [PRONOM](https://www.nationalarchives.gov.uk/PRONOM/) identifier
* Format version
* Identification proof and note
* filename extension

os.Stat() is giving you the

* File size
* File creation time
* File modification time
* File access time

* and the same for directories


Furthermore it creates and calculates

* UUIDv4s as unique identifiers (not stable across sessions)
* hash sums (md5, sha1, sha256, sha512 and blake2b-512)
* the entropy of each file (up to 1GB)

and it extracts some EXIF metadata and you can add your own [DublinCore Elements](https://www.dublincore.org/specifications/dublin-core/usageguide/elements/) metadata to scans.

FileTrove also checks if the file is in the NSRL (https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl).

For this check a 4.0GB BoltDB is needed and can be downloaded with FileTrove during the installation. 

You can also create your own database for the NSRL check. You just need a text file with SHA1 hashes, one per line and the tool admftrove from this repository. With this tool you can also add your own hashes to an existing database.

All results are written into a SQLite database and can be exported to TSV files.


## How to install
1. Download a release from https://github.com/steffenfritz/FileTrove/releases or compile from source (using _task build_ in cmd/ftrove (https://taskfile.dev)).
2. Copy the file where you want to install ftrove (the downloaded file has a suffix, omitted in the following documentation)
3. Run `./ftrove --install .`  (Mind the period)
   
	a) If you don't have already a NSRL database, you have to download it. Please be patient.
    
	b) If you have a NSRL database copy/move it do the "db" directory that ftrove just created.

4. You are ready to go!

## How to run
`./ftrove -h` gives you all flags ftrove understands.

A run only with necessary flags looks like this:

`./ftrove -i $DIRECTORY`

where $DIRECTORY is a directory you want to use as a starting point. FileTrove will walk this directory recursively down.

## How to see the results
You can export the results via `./ftrove -t $UUID` where $UUID is the session id. 
Every indexing run gets its own session id. You get a list of all sessions using `./ftrove -l`. 

Example:

1. `./ftrove -l`
2. `./ftrove -t 926be141-ab75-4106-8236-34edfcf102f2`

This will create several TSV files that can be read with Excel, Numbers and your preferred text editor. 


You can also work with SQL on the database, using sqlite on the console or a GUI like sqlitebrowser (https://sqlitebrowser.org/). Sqliteviz is also a neat tool to visualize the data (https://sqliteviz.com/app/#/).

## Background
FileTrove is the successor of [filedriller](https://github.com/steffenfritz/filedriller) and based on my iPres 2021 paper [Marrying siegfried and the National Software Reference Library](https://phaidra.univie.ac.at/detail/o:1424904)
