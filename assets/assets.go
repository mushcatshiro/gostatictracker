package assets

import "embed"

//go:embed templates/*
var TemplateFS embed.FS

//go:embed sql/*
var SqlFS embed.FS
