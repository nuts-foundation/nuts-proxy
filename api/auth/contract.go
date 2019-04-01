package auth

import "github.com/sirupsen/logrus"

type Contract struct {
	Type             string   `json: "type"`
	Version          string   `json: "version"`
	Language         string   `json: "language"`
	SignerAttributes []string `json: "signer_attributes"`
	Template         string   `json: "template"`
}

func ContractByType(contractType string, language string) *Contract {

	switch contractType {
	case "PractitionerLogin", "BehandelaarLogin":
		return &Contract{
			Type:             contractType,
			Version:          "v1",
			Language:         "NL",
			SignerAttributes: []string{"irma-demo.nuts.agb.agbcode"},
			Template: `NL:BehandelaarLogin:v1 Ondergetekende geeft toestemming aan {{acting_party}} om uit zijn/haar naam het nuts netwerk te bevragen. Deze toestemming is geldig van {{valid_from}} tot {{valid_to}}.`, }
	}

	logrus.Warnf("Contract with type %v and language %v not found", contractType, language)

	return nil

}
