# FileTrove
STATUS: Development

FileTrove indexes files and creates metadata from them.

The single binary application walks a directory tree and identifies all regular files by type with siegfried. Furthermore it creates UUIDv4s, hash sums (md5, sha1, sha256, sha512 or blake2b-512) and FileTrove checks if the file is in the NSRL.
