package config

import (
	"github.com/gigurra/kdump/internal/cliUtil"
	"github.com/urfave/cli/v2"
)

var CliFlag = struct {
	OutputDir     cli.StringFlag
	DeletePrevDir cli.BoolFlag
	EncryptKey    cli.StringFlag
}{
	OutputDir: cli.StringFlag{
		Name:     "output-dir",
		Aliases:  []string{"o"},
		Usage:    "output directory to create",
		Required: true,
	},
	DeletePrevDir: cli.BoolFlag{
		Name:  "delete-previous-dir",
		Usage: "if to delete previous output directory",
		Value: false,
	},
	EncryptKey: cli.StringFlag{
		Name:  "secrets-encryption-key",
		Usage: "symmetric secrets encryption hex key for aes GCM (lower case 64 chars)",
	},
}

var CliFlags = cliUtil.FindAllFlags(CliFlag)
