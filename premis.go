package filetrove

import (
	"database/sql"
	"encoding/xml"
	"io"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	premisXmlns    = "http://www.loc.gov/premis/v3"
	premisXsiXmlns = "http://www.w3.org/2001/XMLSchema-instance"
	premisSchemaLoc = "http://www.loc.gov/premis/v3 https://www.loc.gov/standards/premis/premis.xsd"
	premisVersion  = "3.0"
)

// --- PREMIS v3 XML structs (literal prefix approach: xmlns:premis declared on root) ---

type premisObjectIdentifier struct {
	Type  string `xml:"premis:objectIdentifierType"`
	Value string `xml:"premis:objectIdentifierValue"`
}

type premisFixity struct {
	Algorithm  string `xml:"premis:messageDigestAlgorithm"`
	Digest     string `xml:"premis:messageDigest"`
	Originator string `xml:"premis:messageDigestOriginator"`
}

type premisFormatDesignation struct {
	Name    string `xml:"premis:formatName"`
	Version string `xml:"premis:formatVersion,omitempty"`
}

type premisFormatRegistry struct {
	Name string `xml:"premis:formatRegistryName"`
	Key  string `xml:"premis:formatRegistryKey"`
	Role string `xml:"premis:formatRegistryRole,omitempty"`
}

type premisFormat struct {
	Designation *premisFormatDesignation `xml:"premis:formatDesignation"`
	Registry    *premisFormatRegistry    `xml:"premis:formatRegistry,omitempty"`
}

type premisObjectCharacteristics struct {
	CompositionLevel int            `xml:"premis:compositionLevel"`
	Fixities         []premisFixity `xml:"premis:fixity"`
	Size             int64          `xml:"premis:size"`
	Format           premisFormat   `xml:"premis:format"`
}

type premisContentLocation struct {
	Type  string `xml:"premis:contentLocationType"`
	Value string `xml:"premis:contentLocationValue"`
}

type premisStorage struct {
	ContentLocation premisContentLocation `xml:"premis:contentLocation"`
}

type premisLinkingEventIdentifier struct {
	Type  string `xml:"premis:linkingEventIdentifierType"`
	Value string `xml:"premis:linkingEventIdentifierValue"`
}

type premisObject struct {
	XMLName          xml.Name                       `xml:"premis:object"`
	XsiType          string                         `xml:"xsi:type,attr"`
	ObjectIdentifier premisObjectIdentifier         `xml:"premis:objectIdentifier"`
	ObjectChars      premisObjectCharacteristics    `xml:"premis:objectCharacteristics"`
	OriginalName     string                         `xml:"premis:originalName"`
	Storage          premisStorage                  `xml:"premis:storage"`
	LinkingEventIds  []premisLinkingEventIdentifier `xml:"premis:linkingEventIdentifier"`
}

type premisEventIdentifier struct {
	Type  string `xml:"premis:eventIdentifierType"`
	Value string `xml:"premis:eventIdentifierValue"`
}

type premisLinkingAgentIdentifier struct {
	Type  string `xml:"premis:linkingAgentIdentifierType"`
	Value string `xml:"premis:linkingAgentIdentifierValue"`
}

type premisEvent struct {
	XMLName         xml.Name                       `xml:"premis:event"`
	EventIdentifier premisEventIdentifier          `xml:"premis:eventIdentifier"`
	EventType       string                         `xml:"premis:eventType"`
	EventDateTime   string                         `xml:"premis:eventDateTime"`
	LinkingAgents   []premisLinkingAgentIdentifier `xml:"premis:linkingAgentIdentifier"`
}

type premisAgentIdentifier struct {
	Type  string `xml:"premis:agentIdentifierType"`
	Value string `xml:"premis:agentIdentifierValue"`
}

type premisAgent struct {
	XMLName         xml.Name              `xml:"premis:agent"`
	AgentIdentifier premisAgentIdentifier `xml:"premis:agentIdentifier"`
	AgentName       string                `xml:"premis:agentName,omitempty"`
	AgentType       string                `xml:"premis:agentType"`
	AgentVersion    string                `xml:"premis:agentVersion,omitempty"`
}

// ExportSessionPREMIS writes all file objects for a session as a PREMIS v3 XML document to w.
// The document contains one Agent (FileTrove software), one Event (ingestion) per session,
// and one Object per file. Streaming: files are encoded row by row without full in-memory load.
func ExportSessionPREMIS(sessionuuid string, w io.Writer) error {
	db, err := sql.Open("sqlite3", filepath.Join("db", "filetrove.db"))
	if err != nil {
		return err
	}
	defer db.Close()

	smd, err := querySessionForPREMIS(db, sessionuuid)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")

	rootStart := xml.StartElement{
		Name: xml.Name{Local: "premis:premis"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "xmlns:premis"}, Value: premisXmlns},
			{Name: xml.Name{Local: "xmlns:xsi"}, Value: premisXsiXmlns},
			{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: premisSchemaLoc},
			{Name: xml.Name{Local: "version"}, Value: premisVersion},
		},
	}
	if err := enc.EncodeToken(rootStart); err != nil {
		return err
	}

	if err := enc.Encode(buildSoftwareAgent(smd.Filetroveversion)); err != nil {
		return err
	}
	if smd.Archivistname != "" {
		if err := enc.Encode(buildArchivistAgent(smd.Archivistname)); err != nil {
			return err
		}
	}

	if err := enc.Encode(buildIngestionEvent(sessionuuid, smd.Starttime)); err != nil {
		return err
	}

	rows, err := db.Query(
		"SELECT fileuuid, filename, filepath, filesize, "+
			"COALESCE(filemd5,''), COALESCE(filesha1,''), COALESCE(filesha256,''), "+
			"COALESCE(filesha512,''), COALESCE(fileblake2b,''), "+
			"COALESCE(filesfmime,''), COALESCE(filesfformatname,''), "+
			"COALESCE(filesfformatversion,''), COALESCE(filesffmt,'') "+
			"FROM files WHERE sessionuuid=?", sessionuuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			fileuuid, filename, fpath                        string
			filesize                                          int64
			md5, sha1, sha256, sha512, blake2b               string
			mime, formatname, formatversion, fmt             string
		)
		if err := rows.Scan(
			&fileuuid, &filename, &fpath, &filesize,
			&md5, &sha1, &sha256, &sha512, &blake2b,
			&mime, &formatname, &formatversion, &fmt,
		); err != nil {
			return err
		}
		obj := buildFileObject(fileuuid, filename, fpath, filesize,
			md5, sha1, sha256, sha512, blake2b,
			mime, formatname, formatversion, fmt,
			sessionuuid)
		if err := enc.Encode(obj); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if err := enc.EncodeToken(rootStart.End()); err != nil {
		return err
	}
	return enc.Flush()
}

func querySessionForPREMIS(db *sql.DB, sessionuuid string) (SessionMD, error) {
	var s SessionMD
	row := db.QueryRow(
		"SELECT uuid, starttime, COALESCE(endtime,''), COALESCE(project,''), "+
			"COALESCE(archivistname,''), COALESCE(mountpoint,''), pathseparator, "+
			"COALESCE(exifflag,''), COALESCE(dublincoreflag,''), COALESCE(yaraflag,''), "+
			"COALESCE(yarasource,''), COALESCE(xattrflag,''), COALESCE(ntfsadsflag,''), "+
			"filetroveversion, filetrovedbversion, nsrlversion, siegfriedversion, goversion "+
			"FROM sessionsmd WHERE uuid=?", sessionuuid)
	err := row.Scan(
		&s.UUID, &s.Starttime, &s.Endtime, &s.Project, &s.Archivistname, &s.Mountpoint,
		&s.Pathseparator, &s.ExifFlag, &s.Dublincoreflag, &s.Yaraflag, &s.Yarasource,
		&s.XattrFlag, &s.NtfsadsFlag, &s.Filetroveversion, &s.Filetrovedbversion,
		&s.Nsrlversion, &s.Sfversion, &s.Goversion,
	)
	return s, err
}

func buildSoftwareAgent(version string) premisAgent {
	return premisAgent{
		AgentIdentifier: premisAgentIdentifier{
			Type:  "local",
			Value: "FileTrove-software",
		},
		AgentName:    "FileTrove",
		AgentType:    "software",
		AgentVersion: version,
	}
}

func buildArchivistAgent(name string) premisAgent {
	return premisAgent{
		AgentIdentifier: premisAgentIdentifier{
			Type:  "local",
			Value: "FileTrove-archivist",
		},
		AgentName: name,
		AgentType: "person",
	}
}

func buildIngestionEvent(sessionuuid, datetime string) premisEvent {
	return premisEvent{
		EventIdentifier: premisEventIdentifier{
			Type:  "UUID",
			Value: sessionuuid,
		},
		EventType:     "ingestion",
		EventDateTime: datetime,
		LinkingAgents: []premisLinkingAgentIdentifier{{
			Type:  "local",
			Value: "FileTrove-software",
		}},
	}
}

func buildFileObject(fileuuid, filename, fpath string, filesize int64,
	md5, sha1, sha256, sha512, blake2b string,
	mime, formatname, formatversion, fmtID string,
	sessionuuid string) premisObject {

	fixities := buildFixities(md5, sha1, sha256, sha512, blake2b)

	format := buildFormat(mime, formatname, formatversion, fmtID)

	return premisObject{
		XsiType: "premis:file",
		ObjectIdentifier: premisObjectIdentifier{
			Type:  "UUID",
			Value: fileuuid,
		},
		ObjectChars: premisObjectCharacteristics{
			CompositionLevel: 0,
			Fixities:         fixities,
			Size:             filesize,
			Format:           format,
		},
		OriginalName: filename,
		Storage: premisStorage{
			ContentLocation: premisContentLocation{
				Type:  "filepath",
				Value: fpath,
			},
		},
		LinkingEventIds: []premisLinkingEventIdentifier{{
			Type:  "UUID",
			Value: sessionuuid,
		}},
	}
}

func buildFixities(md5, sha1, sha256, sha512, blake2b string) []premisFixity {
	algos := []struct {
		name   string
		digest string
	}{
		{"MD5", md5},
		{"SHA-1", sha1},
		{"SHA-256", sha256},
		{"SHA-512", sha512},
		{"BLAKE2b-512", blake2b},
	}
	var fixities []premisFixity
	for _, a := range algos {
		if a.digest != "" {
			fixities = append(fixities, premisFixity{
				Algorithm:  a.name,
				Digest:     a.digest,
				Originator: "FileTrove",
			})
		}
	}
	return fixities
}

func buildFormat(mime, formatname, formatversion, fmtID string) premisFormat {
	name := formatname
	if name == "" {
		name = mime
	}

	f := premisFormat{
		Designation: &premisFormatDesignation{
			Name:    name,
			Version: formatversion,
		},
	}

	if fmtID != "" {
		f.Registry = &premisFormatRegistry{
			Name: "PRONOM",
			Key:  fmtID,
			Role: "specification",
		}
	}

	return f
}
