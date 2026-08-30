package filetrove

import (
	"database/sql"
	"fmt"
	"strings"
)

// SessionSummary holds a session with aggregated counts for web display
type SessionSummary struct {
	Session   SessionMD
	FileCount int
	DirCount  int
}

// WebFileMD holds file metadata with UUIDs needed for web display
type WebFileMD struct {
	FileUUID            string
	SessionUUID         string
	Filename            string
	Filepath            string
	Ext                 string
	Filesize            int64
	Filemd5             string
	Filesha1            string
	Filesha256          string
	Filesha512          string
	Fileblake2b         string
	Filesffmt           string
	Filesfmime          string
	Filesfformatname    string
	Filesfformatversion string
	Filesfidentnote     string
	Filesfidentproof    string
	Filectime           string
	Filemtime           string
	Fileatime           string
	Filebtime           string
	Filensrl            string
	Fileentropy         float64
	HasYara             bool
}

// FileFilters holds filter and pagination parameters for file queries
// NSRL accepts: "" (all), "only" (NSRL known), "exclude" (non-NSRL only)
type FileFilters struct {
	Query       string
	QueryNegate bool
	Ext         string
	Mimes       []string
	NSRL        string
	YaraOnly    bool
	SortBy   string
	Order    string
	Limit    int
	Offset   int
}

// WebExifRow holds EXIF metadata for web display
type WebExifRow struct {
	ExifVersion  string
	DateTime     string
	DateTimeOrig string
	Artist       string
	Copyright    string
	Make         string
	XPTitle      string
	XPComment    string
	XPAuthor     string
	XPKeywords   string
	XPSubject    string
}

// WebYaraRow holds a single YARA rule match
type WebYaraRow struct {
	RuleName string
}

// WebXattrRow holds a single extended attribute entry
type WebXattrRow struct {
	Name  string
	Value string
}

// WebNtfsAdsRow holds a single NTFS alternate data stream entry
type WebNtfsAdsRow struct {
	AdsName  string
	AdsValue string
}

// WebDirMD holds directory metadata for web display
type WebDirMD struct {
	DirUUID   string
	Dirname   string
	Dirpath   string
	Dirctime  string
	Dirmtime  string
	Diratime  string
	Hierarchy int
}

// FileDetail holds a file and all its related table data
type FileDetail struct {
	File    WebFileMD
	Exif    *WebExifRow
	Yara    []WebYaraRow
	Xattr   []WebXattrRow
	NtfsAds []WebNtfsAdsRow
}

// GetSessionSummaries returns all sessions with file and directory counts
func GetSessionSummaries(db *sql.DB) ([]SessionSummary, error) {
	rows, err := db.Query(`
		SELECT s.uuid,
		       COALESCE(s.starttime,''), COALESCE(s.endtime,''),
		       COALESCE(s.project,''), COALESCE(s.archivistname,''),
		       COALESCE(s.mountpoint,''), COALESCE(s.pathseparator,''),
		       COALESCE(s.exifflag,''), COALESCE(s.dublincoreflag,''),
		       COALESCE(s.yaraflag,''), COALESCE(s.yarasource,''),
		       COALESCE(s.xattrflag,''), COALESCE(s.ntfsadsflag,''),
		       COALESCE(s.filetroveversion,''), COALESCE(s.filetrovedbversion,''),
		       COALESCE(s.nsrlversion,''), COALESCE(s.siegfriedversion,''),
		       COALESCE(s.goversion,''),
		       (SELECT COUNT(*) FROM files f WHERE f.sessionuuid = s.uuid) AS filecount,
		       (SELECT COUNT(*) FROM directories d WHERE d.sessionuuid = s.uuid) AS dircount
		FROM sessionsmd s
		ORDER BY s.starttime DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []SessionSummary
	for rows.Next() {
		var ss SessionSummary
		s := &ss.Session
		if err := rows.Scan(
			&s.UUID, &s.Starttime, &s.Endtime, &s.Project,
			&s.Archivistname, &s.Mountpoint, &s.Pathseparator,
			&s.ExifFlag, &s.Dublincoreflag, &s.Yaraflag, &s.Yarasource,
			&s.XattrFlag, &s.NtfsadsFlag, &s.Filetroveversion,
			&s.Filetrovedbversion, &s.Nsrlversion, &s.Sfversion, &s.Goversion,
			&ss.FileCount, &ss.DirCount,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, ss)
	}
	return summaries, rows.Err()
}

// GetSessionByUUID returns a single session's metadata
func GetSessionByUUID(db *sql.DB, uuid string) (SessionMD, error) {
	var s SessionMD
	row := db.QueryRow(`
		SELECT uuid,
		       COALESCE(starttime,''), COALESCE(endtime,''),
		       COALESCE(project,''), COALESCE(archivistname,''),
		       COALESCE(mountpoint,''), COALESCE(pathseparator,''),
		       COALESCE(exifflag,''), COALESCE(dublincoreflag,''),
		       COALESCE(yaraflag,''), COALESCE(yarasource,''),
		       COALESCE(xattrflag,''), COALESCE(ntfsadsflag,''),
		       COALESCE(filetroveversion,''), COALESCE(filetrovedbversion,''),
		       COALESCE(nsrlversion,''), COALESCE(siegfriedversion,''),
		       COALESCE(goversion,'')
		FROM sessionsmd WHERE uuid = ?`, uuid)
	err := row.Scan(
		&s.UUID, &s.Starttime, &s.Endtime, &s.Project,
		&s.Archivistname, &s.Mountpoint, &s.Pathseparator,
		&s.ExifFlag, &s.Dublincoreflag, &s.Yaraflag, &s.Yarasource,
		&s.XattrFlag, &s.NtfsadsFlag, &s.Filetroveversion,
		&s.Filetrovedbversion, &s.Nsrlversion, &s.Sfversion, &s.Goversion,
	)
	return s, err
}

// QueryFiles returns files for a session with optional filters and pagination.
// Returns the matching files, total count of matches, and any error.
func QueryFiles(db *sql.DB, sessionUUID string, f FileFilters) ([]WebFileMD, int, error) {
	conds := []string{"f.sessionuuid = ?"}
	args := []interface{}{sessionUUID}

	if f.Query != "" {
		op := "LIKE"
		if f.QueryNegate {
			op = "NOT LIKE"
		}
		conds = append(conds, "(f.filename "+op+" ? OR f.filepath "+op+" ?)")
		like := "%" + f.Query + "%"
		args = append(args, like, like)
	}
	if f.Ext != "" {
		conds = append(conds, "f.filenameextension = ?")
		args = append(args, f.Ext)
	}
	if len(f.Mimes) > 0 {
		placeholders := strings.Repeat("?,", len(f.Mimes))
		placeholders = placeholders[:len(placeholders)-1]
		conds = append(conds, "f.filesfmime IN ("+placeholders+")")
		for _, m := range f.Mimes {
			args = append(args, m)
		}
	}
	switch f.NSRL {
	case "only":
		conds = append(conds, "f.filensrl = 'TRUE'")
	case "exclude":
		conds = append(conds, "(f.filensrl != 'TRUE' OR f.filensrl IS NULL)")
	}
	if f.YaraOnly {
		conds = append(conds, "EXISTS (SELECT 1 FROM yara y WHERE y.fileuuid = f.fileuuid AND y.sessionuuid = f.sessionuuid)")
	}

	where := strings.Join(conds, " AND ")

	var sortCol string
	switch f.SortBy {
	case "filesize":
		sortCol = "filesize"
	case "filemtime":
		sortCol = "filemtime"
	case "fileentropy":
		sortCol = "fileentropy"
	case "filenameextension":
		sortCol = "filenameextension"
	case "filesfmime":
		sortCol = "filesfmime"
	default:
		sortCol = "filename"
	}
	order := "ASC"
	if strings.ToUpper(f.Order) == "DESC" {
		order = "DESC"
	}

	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM files f WHERE "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 100
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	query := fmt.Sprintf(`
		SELECT f.fileuuid, f.sessionuuid,
		       COALESCE(f.filename,''), COALESCE(f.filepath,''),
		       COALESCE(f.filenameextension,''), COALESCE(f.filesize,0),
		       COALESCE(f.filemd5,''), COALESCE(f.filesha1,''), COALESCE(f.filesha256,''),
		       COALESCE(f.filesha512,''), COALESCE(f.fileblake2b,''),
		       COALESCE(f.filesffmt,''), COALESCE(f.filesfmime,''),
		       COALESCE(f.filesfformatname,''), COALESCE(f.filesfformatversion,''),
		       COALESCE(f.filesfidentnote,''), COALESCE(f.filesfidentproof,''),
		       COALESCE(f.filectime,''), COALESCE(f.filemtime,''), COALESCE(f.fileatime,''),
		       COALESCE(f.filebtime,''),
		       COALESCE(f.filensrl,''), COALESCE(f.fileentropy,0),
		       EXISTS(SELECT 1 FROM yara y WHERE y.fileuuid = f.fileuuid) AS has_yara
		FROM files f
		WHERE %s
		ORDER BY f.%s %s
		LIMIT ? OFFSET ?`, where, sortCol, order)

	args = append(args, limit, f.Offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []WebFileMD
	for rows.Next() {
		var wf WebFileMD
		if err := rows.Scan(
			&wf.FileUUID, &wf.SessionUUID, &wf.Filename, &wf.Filepath,
			&wf.Ext, &wf.Filesize, &wf.Filemd5, &wf.Filesha1, &wf.Filesha256,
			&wf.Filesha512, &wf.Fileblake2b, &wf.Filesffmt, &wf.Filesfmime,
			&wf.Filesfformatname, &wf.Filesfformatversion, &wf.Filesfidentnote,
			&wf.Filesfidentproof, &wf.Filectime, &wf.Filemtime, &wf.Fileatime,
			&wf.Filebtime,
			&wf.Filensrl, &wf.Fileentropy, &wf.HasYara,
		); err != nil {
			return nil, 0, err
		}
		files = append(files, wf)
	}
	return files, total, rows.Err()
}

// GetFileDetail returns a file with all related table data
func GetFileDetail(db *sql.DB, fileUUID string) (FileDetail, error) {
	var d FileDetail

	row := db.QueryRow(`
		SELECT fileuuid, sessionuuid,
		       COALESCE(filename,''), COALESCE(filepath,''),
		       COALESCE(filenameextension,''), COALESCE(filesize,0),
		       COALESCE(filemd5,''), COALESCE(filesha1,''), COALESCE(filesha256,''),
		       COALESCE(filesha512,''), COALESCE(fileblake2b,''),
		       COALESCE(filesffmt,''), COALESCE(filesfmime,''),
		       COALESCE(filesfformatname,''), COALESCE(filesfformatversion,''),
		       COALESCE(filesfidentnote,''), COALESCE(filesfidentproof,''),
		       COALESCE(filectime,''), COALESCE(filemtime,''), COALESCE(fileatime,''),
		       COALESCE(filebtime,''),
		       COALESCE(filensrl,''), COALESCE(fileentropy,0)
		FROM files WHERE fileuuid = ?`, fileUUID)
	if err := row.Scan(
		&d.File.FileUUID, &d.File.SessionUUID,
		&d.File.Filename, &d.File.Filepath, &d.File.Ext, &d.File.Filesize,
		&d.File.Filemd5, &d.File.Filesha1, &d.File.Filesha256,
		&d.File.Filesha512, &d.File.Fileblake2b, &d.File.Filesffmt,
		&d.File.Filesfmime, &d.File.Filesfformatname, &d.File.Filesfformatversion,
		&d.File.Filesfidentnote, &d.File.Filesfidentproof,
		&d.File.Filectime, &d.File.Filemtime, &d.File.Fileatime,
		&d.File.Filebtime,
		&d.File.Filensrl, &d.File.Fileentropy,
	); err != nil {
		return d, err
	}

	var e WebExifRow
	exifErr := db.QueryRow(`
		SELECT COALESCE(exifversion,''), COALESCE(datetime,''),
		       COALESCE(datetimeorig,''), COALESCE(artist,''), COALESCE(copyright,''),
		       COALESCE(make,''), COALESCE(xptitle,''), COALESCE(xpcomment,''),
		       COALESCE(xpauthor,''), COALESCE(xpkeywords,''), COALESCE(xpsubject,'')
		FROM exif WHERE fileuuid = ?`, fileUUID).Scan(
		&e.ExifVersion, &e.DateTime, &e.DateTimeOrig, &e.Artist, &e.Copyright,
		&e.Make, &e.XPTitle, &e.XPComment, &e.XPAuthor, &e.XPKeywords, &e.XPSubject,
	)
	if exifErr == nil {
		d.Exif = &e
	}

	yaraRows, err := db.Query("SELECT COALESCE(rulename,'') FROM yara WHERE fileuuid = ?", fileUUID)
	if err != nil {
		return d, err
	}
	defer yaraRows.Close()
	for yaraRows.Next() {
		var y WebYaraRow
		if err := yaraRows.Scan(&y.RuleName); err != nil {
			return d, err
		}
		d.Yara = append(d.Yara, y)
	}

	xattrRows, err := db.Query("SELECT COALESCE(xattrname,''), COALESCE(xattrvalue,'') FROM xattr WHERE fileuuid = ?", fileUUID)
	if err != nil {
		return d, err
	}
	defer xattrRows.Close()
	for xattrRows.Next() {
		var x WebXattrRow
		if err := xattrRows.Scan(&x.Name, &x.Value); err != nil {
			return d, err
		}
		d.Xattr = append(d.Xattr, x)
	}

	adsRows, err := db.Query("SELECT COALESCE(adsname,''), COALESCE(adsvalue,'') FROM ntfsads WHERE fileuuid = ?", fileUUID)
	if err != nil {
		return d, err
	}
	defer adsRows.Close()
	for adsRows.Next() {
		var a WebNtfsAdsRow
		if err := adsRows.Scan(&a.AdsName, &a.AdsValue); err != nil {
			return d, err
		}
		d.NtfsAds = append(d.NtfsAds, a)
	}

	return d, nil
}

// GetSessionDirs returns directories for a session, optionally filtered by a search string.
func GetSessionDirs(db *sql.DB, sessionUUID string, query string) ([]WebDirMD, error) {
	q := `SELECT diruuid,
		       COALESCE(dirname,''), COALESCE(dirpath,''),
		       COALESCE(dircttime,''), COALESCE(dirmtime,''), COALESCE(diratime,''),
		       COALESCE(hierarchy,0)
		FROM directories WHERE sessionuuid = ?`
	args := []interface{}{sessionUUID}
	if query != "" {
		q += ` AND (dirname LIKE ? OR dirpath LIKE ?)`
		like := "%" + query + "%"
		args = append(args, like, like)
	}
	q += ` ORDER BY dirpath`
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dirs []WebDirMD
	for rows.Next() {
		var dir WebDirMD
		if err := rows.Scan(&dir.DirUUID, &dir.Dirname, &dir.Dirpath,
			&dir.Dirctime, &dir.Dirmtime, &dir.Diratime, &dir.Hierarchy); err != nil {
			return nil, err
		}
		dirs = append(dirs, dir)
	}
	return dirs, rows.Err()
}

// GetDistinctMimes returns distinct MIME types present in a session
func GetDistinctMimes(db *sql.DB, sessionUUID string) ([]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT filesfmime FROM files
		WHERE sessionuuid = ? AND filesfmime != '' AND filesfmime IS NOT NULL
		ORDER BY filesfmime`, sessionUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mimes []string
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}
		mimes = append(mimes, m)
	}
	return mimes, rows.Err()
}

// GetDistinctExtensions returns distinct file extensions present in a session
func GetDistinctExtensions(db *sql.DB, sessionUUID string) ([]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT filenameextension FROM files
		WHERE sessionuuid = ? AND filenameextension != '' AND filenameextension IS NOT NULL
		ORDER BY filenameextension`, sessionUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exts []string
	for rows.Next() {
		var e string
		if err := rows.Scan(&e); err != nil {
			return nil, err
		}
		exts = append(exts, e)
	}
	return exts, rows.Err()
}
