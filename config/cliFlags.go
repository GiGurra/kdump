package config

import (
	"github.com/gigurra/kdump/internal/refl"
	"github.com/urfave/cli/v2"
	"reflect"
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

var CliFlags = findAllFlags()

////////////////////////////////////////////////////////////////////
// Private helpers below...  prob should exist a better solution :D

func findAllFlags() []cli.Flag {

	out := make([]cli.Flag, 0, 20)

	object := reflect.ValueOf(CliFlag)
	count := object.NumField()

	for i := 0; i < count; i++ {
		out = append(out, refl.GetPtrToFieldCopy(object, i).Interface().(cli.Flag))
	}

	return out
}
