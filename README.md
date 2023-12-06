# FileTrove
STATUS: Development

![Build Status](https://github.com/steffenfritz/FileTrove/actions/workflows/buildstatus.yml/badge.svg)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

FileTrove indexes files and creates metadata from them.

The single binary application walks a directory tree and identifies all regular files by type with siegfried, giving you the 

* MIME type
* PRONOM identifier
* Format version
* Identification proof and note


os.Stat() is giving you the

* File size
* File creation time
* File modification time
* File access time


Furthermore it creates and calculates

* UUIDv4s as unique identifiers (not stable across sessions)
* hash sums (md5, sha1, sha256, sha512 and blake2b-512)
* the entropy of each file (up to 1GB)

and it extracts some EXIF metadata and you can add your own DublinCore metadata to scans.

FileTrove also checks if the file is in the NSRL. 

For this check a 3.2GB BoltDB is needed and can be downloaded with FileTrive during the installation. 

You can also create your own database for the NSRL check. You just need a text file with SHA1 hashes, one per line and the tool admftrove from this repository. With this tool you can also add your own hashes to an existing database.

All results are written into a SQLite database and can be exported to TSV files.


## How to install
1. Download a release from https://github.com/steffenfritz/FileTrove/releases or compile from source.
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
