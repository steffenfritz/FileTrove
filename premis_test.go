package filetrove

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestBuildFixities_allPresent(t *testing.T) {
	f := buildFixities("md5val", "sha1val", "sha256val", "sha512val", "blake2bval")
	if len(f) != 5 {
		t.Fatalf("expected 5 fixity entries, got %d", len(f))
	}
	algos := []string{"MD5", "SHA-1", "SHA-256", "SHA-512", "BLAKE2b-512"}
	for i, a := range algos {
		if f[i].Algorithm != a {
			t.Errorf("fixity[%d].Algorithm = %q, want %q", i, f[i].Algorithm, a)
		}
		if f[i].Originator != "FileTrove" {
			t.Errorf("fixity[%d].Originator = %q, want FileTrove", i, f[i].Originator)
		}
	}
}

func TestBuildFixities_skipEmpty(t *testing.T) {
	f := buildFixities("md5val", "", "sha256val", "", "")
	if len(f) != 2 {
		t.Fatalf("expected 2 fixity entries, got %d", len(f))
	}
	if f[0].Algorithm != "MD5" || f[1].Algorithm != "SHA-256" {
		t.Errorf("unexpected algorithms: %v", f)
	}
}

func TestBuildFormat_withPronom(t *testing.T) {
	f := buildFormat("image/jpeg", "JPEG File Interchange Format", "1.02", "fmt/44")
	if f.Designation == nil {
		t.Fatal("Designation is nil")
	}
	if f.Designation.Name != "JPEG File Interchange Format" {
		t.Errorf("Name = %q", f.Designation.Name)
	}
	if f.Designation.Version != "1.02" {
		t.Errorf("Version = %q", f.Designation.Version)
	}
	if f.Registry == nil {
		t.Fatal("Registry is nil")
	}
	if f.Registry.Name != "PRONOM" || f.Registry.Key != "fmt/44" {
		t.Errorf("Registry = %+v", f.Registry)
	}
}

func TestBuildFormat_fallbackToMime(t *testing.T) {
	f := buildFormat("application/pdf", "", "", "")
	if f.Designation.Name != "application/pdf" {
		t.Errorf("expected MIME fallback, got %q", f.Designation.Name)
	}
	if f.Registry != nil {
		t.Error("Registry should be nil when fmtID is empty")
	}
}

func TestBuildSoftwareAgent(t *testing.T) {
	a := buildSoftwareAgent("1.0.0-TEST")
	if a.AgentType != "software" {
		t.Errorf("AgentType = %q", a.AgentType)
	}
	if a.AgentVersion != "1.0.0-TEST" {
		t.Errorf("AgentVersion = %q", a.AgentVersion)
	}
	if a.AgentIdentifier.Value != "FileTrove-software" {
		t.Errorf("AgentIdentifier.Value = %q", a.AgentIdentifier.Value)
	}
}

func TestBuildIngestionEvent(t *testing.T) {
	e := buildIngestionEvent("test-uuid", "2024-01-01T00:00:00Z")
	if e.EventType != "ingestion" {
		t.Errorf("EventType = %q", e.EventType)
	}
	if e.EventIdentifier.Value != "test-uuid" {
		t.Errorf("EventIdentifier.Value = %q", e.EventIdentifier.Value)
	}
	if len(e.LinkingAgents) != 1 || e.LinkingAgents[0].Value != "FileTrove-software" {
		t.Errorf("LinkingAgents = %v", e.LinkingAgents)
	}
}

func TestBuildFileObject_xmlRoundtrip(t *testing.T) {
	obj := buildFileObject(
		"file-uuid-1", "test.jpg", "/data/test.jpg", 12345,
		"md5abc", "sha1abc", "sha256abc", "sha512abc", "blake2babc",
		"image/jpeg", "JPEG File Interchange Format", "1.02", "fmt/44",
		"session-uuid-1",
	)

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(obj); err != nil {
		t.Fatalf("encode error: %v", err)
	}

	out := buf.String()

	checks := []string{
		"premis:object",
		"xsi:type",
		"premis:file",
		"file-uuid-1",
		"premis:objectIdentifier",
		"premis:fixity",
		"MD5",
		"md5abc",
		"FileTrove",
		"SHA-1",
		"BLAKE2b-512",
		"12345",
		"JPEG File Interchange Format",
		"1.02",
		"PRONOM",
		"fmt/44",
		"test.jpg",
		"/data/test.jpg",
		"session-uuid-1",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("XML output missing %q\n%s", want, out)
		}
	}
}

func TestPREMISRootElement(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(&buf)
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
		t.Fatal(err)
	}
	if err := enc.Encode(buildSoftwareAgent("1.0.0")); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(buildIngestionEvent("evt-uuid", "2024-01-01T00:00:00Z")); err != nil {
		t.Fatal(err)
	}
	if err := enc.EncodeToken(rootStart.End()); err != nil {
		t.Fatal(err)
	}
	if err := enc.Flush(); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{
		`xmlns:premis="http://www.loc.gov/premis/v3"`,
		`xmlns:xsi=`,
		`version="3.0"`,
		"premis:agent",
		"premis:event",
		"ingestion",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("root XML missing %q", want)
		}
	}
}
