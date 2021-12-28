package config

import "github.com/urfave/cli/v2"

var OutputDirFlag = cli.StringFlag{
	Name:     "output-dir",
	Aliases:  []string{"o"},
	Usage:    "output directory to create",
	Required: true,
}

var DeletePrevDirFlag = cli.BoolFlag{
	Name:  "delete-previous-dir",
	Usage: "if to delete previous output directory",
	Value: false,
}

var EncryptKeyFlag = cli.StringFlag{
	Name:  "secrets-encryption-key",
	Usage: "symmetric secrets encryption hex key for aes GCM (lower case 64 chars)",
}

var CliFlags = []cli.Flag{
	&OutputDirFlag,
	&DeletePrevDirFlag,
	&EncryptKeyFlag,
}
