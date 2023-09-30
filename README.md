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

All results are written into a SQLite database and can be exported to TSV files.


## How to install
1. ~~Download a release from https://github.com/steffenfritz/FileTrove/releases or~~ compile from source
2. Copy the file where you want to install ftrove
3. Run "./ftrove --install ."  (Mind the period)
   
	a) If you don't have already a NSRL database, you have to download it. Please be patient.
    
	b) If you have a NSRL database copy/move it do the "db" directory that ftrove just created.

5. You are ready to go!

## How to run
"./ftrove -h" gives you all flags ftrove understands.

A run only with necessary flags looks like this:

./ftrove -i $DIRECTORY

where $DIRECTORY is a directory you want to use as a starting point. FileTrove will walk this directory recursively down.

## How to see the results
You can export the results via "./ftrove -t $UUID" where $UUID is the session id. 
Every indexing run gets its own session id. You get a list of all sessions using "./ftrove -l". 

Example:

1. ./ftrove -l
2. ./ftrove -t 926be141-ab75-4106-8236-34edfcf102f2 

This will create two TSV files (directories and files) that can be read with Excel, Numbers and your preferred text editor. 

You can also work with SQL on the database, using sqlite on the console or a GUI like sqlitebrowser (https://sqlitebrowser.org/).
