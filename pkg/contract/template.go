/*
 * Nuts auth
 * Copyright (C) 2020. Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package contract

import (
	"errors"
	"fmt"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/goodsign/monday"
)

const timeLayout = "Monday, 2 January 2006 15:04:05"
const ValidFromAttr = "valid_from"
const ValidToAttr = "valid_to"

// deprecated
const ActingPartyAttr = "acting_party"
const LegalEntityAttr = "legal_entity"

// Template stores al properties of a contract template which can result in a signed contract
type Template struct {
	Type                 Type     `json:"type"`
	Version              Version  `json:"version"`
	Language             Language `json:"language"`
	SignerAttributes     []string `json:"signer_attributes"`
	SignerDemoAttributes []string `json:"-"`
	Template             string   `json:"Template"`
	TemplateAttributes   []string `json:"template_attributes"`
	Regexp               string   `json:"-"`
}

// Language of the contract in all caps. example: "NL"
type Language string

// Type contains type of the contract to sign. Example: "BehandelaarLogin"
type Type string

// Version of the contract. example: "v1"
type Version string

// NowFunc is used to store a function that returns the current time. This can be changed when you want to mock the current time.
var NowFunc = time.Now

// StandardSignerAttributes defines the standard list of attributes used for a contract.
// If SignerAttribute name starts with a dot '.', it uses the configured scheme manager
var StandardSignerAttributes = []string{
	".gemeente.personalData.firstnames",
	"pbdf.sidn-pbdf.email.email",
}

func (c Template) timeLocation() *time.Location {
	loc, _ := time.LoadLocation(AmsterdamTimeZone)
	return loc
}

// Render a template using the given templates variables. The combination of validFrom and the duration configure the validFrom and validTo template attributes.
// The ValidFrom or ValidTo provided in the vars map will be overwritten.
// Note: For date calculation the Amsterdam timezone and Dutch locale is used.
func (c Template) Render(vars map[string]string, validFrom time.Time, validDuration time.Duration) (*Contract, error) {
	vars[ValidFromAttr] = monday.Format(validFrom.In(c.timeLocation()), timeLayout, monday.LocaleNlNL)
	vars[ValidToAttr] = monday.Format(validFrom.Add(validDuration).In(c.timeLocation()), timeLayout, monday.LocaleNlNL)

	rawContractText, err := mustache.Render(c.Template, vars)
	if err != nil {
		return nil, fmt.Errorf("could not render contract template: %w", err)
	}

	contract := &Contract{
		RawContractText: rawContractText,
		Template:        &c,
	}
	if err := contract.initParams(); err != nil {
		return nil, err
	}

	return contract, nil
}

// ErrUnknownContractFormat is returned when the contract format is unknown
var ErrUnknownContractFormat = errors.New("unknown contract format")

// ErrInvalidContractFormat indicates tha a contract format is unknown.
var ErrInvalidContractFormat = errors.New("unknown contract type")

// ErrContractNotFound is used when a certain combination of type, language and version cannot resolve to a contract
var ErrContractNotFound = errors.New("contract not found")

// ErrInvalidContractText is used when contract texts cannot be parsed or contain invalid values
var ErrInvalidContractText = errors.New("invalid contract text")
