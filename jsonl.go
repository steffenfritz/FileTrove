package filetrove

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// JSONLRecord is the top-level envelope for every JSONL line.
// The "type" field lets consumers filter records with: jq 'select(.type == "file")'
type JSONLRecord struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// FileRecord mirrors the full files table row for JSONL export.
type FileRecord struct {
	Fileuuid    string `json:"fileuuid"`
	Sessionuuid string `json:"sessionuuid"`
	FileMD
	Hierarchy int64 `json:"hierarchy"`
}

// DirRecord mirrors the full directories table row for JSONL export.
type DirRecord struct {
	Diruuid     string `json:"diruuid"`
	Sessionuuid string `json:"sessionuuid"`
	DirMD
	Hierarchy int64 `json:"hierarchy"`
}

// ExifRecord mirrors the full exif table row for JSONL export.
type ExifRecord struct {
	Exifuuid    string `json:"exifuuid"`
	Sessionuuid string `json:"sessionuuid"`
	Fileuuid    string `json:"fileuuid"`
	ExifParsed
}

// YaraRecord mirrors the full yara table row for JSONL export.
type YaraRecord struct {
	Yaraentryuuid string `json:"yaraentryuuid"`
	Sessionuuid   string `json:"sessionuuid"`
	Fileuuid      string `json:"fileuuid"`
	Rulename      string `json:"rulename"`
}

// XattrRecord mirrors the full xattr table row for JSONL export.
type XattrRecord struct {
	Xattruuid   string `json:"xattruuid"`
	Sessionuuid string `json:"sessionuuid"`
	Fileuuid    string `json:"fileuuid"`
	Xattrname   string `json:"xattrname"`
	Xattrvalue  string `json:"xattrvalue"`
}

// NtfsadsRecord mirrors the full ntfsads table row for JSONL export.
type NtfsadsRecord struct {
	Ntfsadsuuid string `json:"ntfsadsuuid"`
	Sessionuuid string `json:"sessionuuid"`
	Fileuuid    string `json:"fileuuid"`
	Adsname     string `json:"adsname"`
	Adsvalue    string `json:"adsvalue"`
}

// DCRecord mirrors the full dublincore table row for JSONL export.
type DCRecord struct {
	UUID        string `json:"uuid"`
	Sessionuuid string `json:"sessionuuid"`
	DublinCore
}

// ExportSessionJSONL writes all records for the given session as JSONL to w.
// Each line is a self-contained JSON object with a "type" discriminator field.
// Tables exported (in order): session, files, directories, exif, dublincore, yara, xattr, ntfsads.
// Optional tables are silently skipped when they contain no rows for the session.
func ExportSessionJSONL(sessionuuid string, w io.Writer) error {
	db, err := sql.Open("sqlite3", filepath.Join("db", "filetrove.db"))
	if err != nil {
		return err
	}
	defer db.Close()

	enc := json.NewEncoder(w)

	if err := exportSessionRowJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("session: %w", err)
	}
	if err := exportFilesJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("files: %w", err)
	}
	if err := exportDirectoriesJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("directories: %w", err)
	}
	if err := exportExifJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("exif: %w", err)
	}
	if err := exportDCJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("dublincore: %w", err)
	}
	if err := exportYaraJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("yara: %w", err)
	}
	if err := exportXattrJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("xattr: %w", err)
	}
	if err := exportNtfsadsJSONL(db, sessionuuid, enc); err != nil {
		return fmt.Errorf("ntfsads: %w", err)
	}

	return nil
}

func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func exportSessionRowJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	row := db.QueryRow(
		"SELECT uuid, starttime, COALESCE(endtime,''), COALESCE(project,''), "+
			"COALESCE(archivistname,''), COALESCE(mountpoint,''), pathseparator, "+
			"COALESCE(exifflag,''), COALESCE(dublincoreflag,''), COALESCE(yaraflag,''), "+
			"COALESCE(yarasource,''), COALESCE(xattrflag,''), COALESCE(ntfsadsflag,''), "+
			"filetroveversion, filetrovedbversion, nsrlversion, siegfriedversion, goversion "+
			"FROM sessionsmd WHERE uuid=?", sessionuuid)

	var s SessionMD
	if err := row.Scan(
		&s.UUID, &s.Starttime, &s.Endtime, &s.Project, &s.Archivistname, &s.Mountpoint,
		&s.Pathseparator, &s.ExifFlag, &s.Dublincoreflag, &s.Yaraflag, &s.Yarasource,
		&s.XattrFlag, &s.NtfsadsFlag, &s.Filetroveversion, &s.Filetrovedbversion,
		&s.Nsrlversion, &s.Sfversion, &s.Goversion,
	); err != nil {
		return err
	}
	return enc.Encode(JSONLRecord{Type: "session", Payload: s})
}

func exportFilesJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT fileuuid, sessionuuid, filename, filepath, filenameextension, "+
			"filesize, filemd5, filesha1, filesha256, filesha512, fileblake2b, "+
			"filesffmt, filesfmime, filesfformatname, filesfformatversion, "+
			"filesfidentnote, filesfidentproof, filectime, filemtime, fileatime, "+
			"filensrl, fileentropy, hierarchy FROM files WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r FileRecord
		if err := rows.Scan(
			&r.Fileuuid, &r.Sessionuuid,
			&r.FileMD.Filename, &r.FileMD.Filepath, &r.FileMD.Filenameextension,
			&r.FileMD.Filesize, &r.FileMD.Filemd5, &r.FileMD.Filesha1,
			&r.FileMD.Filesha256, &r.FileMD.Filesha512, &r.FileMD.Fileblake2b,
			&r.FileMD.Filesffmt, &r.FileMD.Filesfmime, &r.FileMD.Filesfformatname,
			&r.FileMD.Filesfformatversion, &r.FileMD.Filesfidentnote,
			&r.FileMD.Filesfidentproof, &r.FileMD.Filectime, &r.FileMD.Filemtime,
			&r.FileMD.Fileatime, &r.FileMD.Filensrl, &r.FileMD.Fileentropy,
			&r.Hierarchy,
		); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "file", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func exportDirectoriesJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT diruuid, sessionuuid, dirname, dirpath, dircttime, dirmtime, diratime, hierarchy "+
			"FROM directories WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r DirRecord
		if err := rows.Scan(
			&r.Diruuid, &r.Sessionuuid,
			&r.DirMD.Dirname, &r.DirMD.Dirpath,
			&r.DirMD.Dirctime, &r.DirMD.Dirmtime, &r.DirMD.Diratime,
			&r.Hierarchy,
		); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "directory", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func exportExifJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT exifuuid, sessionuuid, fileuuid, exifversion, datetime, datetimeorig, "+
			"artist, copyright, make, xptitle, xpcomment, xpauthor, xpkeywords, xpsubject "+
			"FROM exif WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r ExifRecord
		if err := rows.Scan(
			&r.Exifuuid, &r.Sessionuuid, &r.Fileuuid,
			&r.ExifParsed.ExifVersion, &r.ExifParsed.DateTime, &r.ExifParsed.DateTimeOrig,
			&r.ExifParsed.Artist, &r.ExifParsed.Copyright, &r.ExifParsed.Make,
			&r.ExifParsed.XPTitle, &r.ExifParsed.XPComment, &r.ExifParsed.XPAuthor,
			&r.ExifParsed.XPKeywords, &r.ExifParsed.XPSubject,
		); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "exif", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func exportDCJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT uuid, sessionuuid, title, creator, contributor, publisher, subject, "+
			"description, date, language, type, format, identifier, source, relation, rights, coverage "+
			"FROM dublincore WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r DCRecord
		if err := rows.Scan(
			&r.UUID, &r.Sessionuuid,
			&r.DublinCore.Title, &r.DublinCore.Creator, &r.DublinCore.Contributor,
			&r.DublinCore.Publisher, &r.DublinCore.Subject, &r.DublinCore.Description,
			&r.DublinCore.Date, &r.DublinCore.Language, &r.DublinCore.Type,
			&r.DublinCore.Format, &r.DublinCore.Identifier, &r.DublinCore.Source,
			&r.DublinCore.Relation, &r.DublinCore.Rights, &r.DublinCore.Coverage,
		); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "dublincore", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func exportYaraJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT yaraentryuuid, sessionuuid, fileuuid, rulename FROM yara WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r YaraRecord
		if err := rows.Scan(&r.Yaraentryuuid, &r.Sessionuuid, &r.Fileuuid, &r.Rulename); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "yara", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func exportXattrJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT xattruuid, sessionuuid, fileuuid, xattrname, xattrvalue FROM xattr WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r XattrRecord
		if err := rows.Scan(&r.Xattruuid, &r.Sessionuuid, &r.Fileuuid, &r.Xattrname, &r.Xattrvalue); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "xattr", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func exportNtfsadsJSONL(db *sql.DB, sessionuuid string, enc *json.Encoder) error {
	rows, err := db.Query(
		"SELECT ntfsadsuuid, sessionuuid, fileuuid, adsname, adsvalue FROM ntfsads WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r NtfsadsRecord
		if err := rows.Scan(&r.Ntfsadsuuid, &r.Sessionuuid, &r.Fileuuid, &r.Adsname, &r.Adsvalue); err != nil {
			return err
		}
		if err := enc.Encode(JSONLRecord{Type: "ntfsads", Payload: r}); err != nil {
			return err
		}
	}
	return rows.Err()
}
