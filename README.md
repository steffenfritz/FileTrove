# FileTrove
STATUS: Development

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


FileTrove also checks if the file is in the NSRL. 

For this check a 5.2GB BoltDB is needed and can be downloaded from https://archive.org/details/nsrl_20230918. You can also create your own database for the NSRL check. You just need a text file with SHA1 hashes, one per line and the tool admftrove from this repository.

All results are written into a SQLite database and can be exportd to TSV files.
