package filetrove

import (
	"github.com/rwcarlsen/goexif/exif"
	"os"
)

type ExifParsed struct {
	ExifVersion  string
	DateTime     string
	DateTimeOrig string
	Artist       string
	Copyright    string
	Make         string
	Software     string
	XPTitle      string
	XPComment    string
	XPAuthor     string
	XPKeywords   string
	XPSubject    string
}

func ExifDecode(fileName string) (ExifParsed, error) {
	var ep ExifParsed

	fd, err := os.Open(fileName)
	if err != nil {
		return ep, err
	}

	imgdecoded, err := exif.Decode(fd)
	if err != nil {
		return ep, err
	}

	exVersion, err := imgdecoded.Get(exif.ExifVersion)
	if err != nil {
		ep.ExifVersion = "not found"
	} else {
		ep.ExifVersion = exVersion.String()
	}

	exDate, err := imgdecoded.Get(exif.DateTime)
	if err != nil {
		ep.DateTime = "not found"
	} else {
		ep.DateTime = exDate.String()
	}

	exDateOrig, err := imgdecoded.Get(exif.DateTimeOriginal)
	if err != nil {
		ep.DateTimeOrig = "not found"
	} else {
		ep.DateTimeOrig = exDateOrig.String()
	}

	exArtist, err := imgdecoded.Get(exif.Artist)
	if err != nil {
		ep.Artist = "not found"
	} else {
		ep.Artist = exArtist.String()
	}

	exCopyright, err := imgdecoded.Get(exif.Copyright)
	if err != nil {
		ep.Copyright = "not found"
	} else {
		ep.Copyright = exCopyright.String()
	}

	exMake, err := imgdecoded.Get(exif.Make)
	if err != nil {
		ep.Make = "not found"
	} else {
		ep.Make = exMake.String()
	}

	exSoftware, err := imgdecoded.Get(exif.Software)
	if err != nil {
		ep.Software = "not found"
	} else {
		ep.Software = exSoftware.String()
	}

	// Microsoft part
	exxptitle, err := imgdecoded.Get(exif.XPTitle)
	if err != nil {
		ep.XPTitle = "not found"
	} else {
		ep.XPTitle = exxptitle.String()
	}

	exxpcomment, err := imgdecoded.Get(exif.XPComment)
	if err != nil {
		ep.XPComment = "not found"
	} else {
		ep.XPComment = exxpcomment.String()
	}

	exxpauthor, err := imgdecoded.Get(exif.XPAuthor)
	if err != nil {
		ep.XPAuthor = "not found"
	} else {
		ep.XPAuthor = exxpauthor.String()
	}

	exxpkeywords, err := imgdecoded.Get(exif.XPKeywords)
	if err != nil {
		ep.XPKeywords = "not found"
	} else {
		ep.XPKeywords = exxpkeywords.String()
	}

	exxpsubject, err := imgdecoded.Get(exif.XPSubject)
	if err != nil {
		ep.XPSubject = "not found"
	} else {
		ep.XPSubject = exxpsubject.String()
	}

	return ep, nil
}
