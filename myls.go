package main

import "time"

// Windows os name
const Windows = "windows"


// file types
const (
	fileRegular int = iota
	fileDirectory
	fileExecutable
	fileCompress
	fileImage
	fileLink
)

// file extensions
const (
	exe = ".exe"
	deb = ".deb"
	zip = ".zip"
	gz = ".gz"
	tar = ".tar"
	rar = ".rar"
	png = ".png"
	jpg = ".jpg"
	jpeg = ".jpeg"
	gif = ".gif"
)

// Primera estructura de nuestro proyecto: para todas las caracteristicas de un archivo/directorio
type file struct {
	name string
	fileType int
	isDir bool
	isHidden bool
	userName string
	groupName string
	size int64
	modificationTime time.Time
	mode string
}

// Segunda estructura de nuestro proyecto: para diferenciar los tipos de archivos
type styleFileType struct {
	icon string
	color string
	symbol string
}

var mapStyleByFileType = map[int]styleFileType {
	fileRegular:    {icon: "📄"},
	fileDirectory:  {icon: "📂", color: "BLUE", symbol: "/"},
	fileExecutable: {icon: "🚀", color: "GREEN", symbol: "*"},
	fileCompress:   {icon: "📦", color: "RED"},
	fileImage:      {icon: "📸", color: "MAGENTA"},
	fileLink:       {icon: "🔗", color: "CYAN"},
}
