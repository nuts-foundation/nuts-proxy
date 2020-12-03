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

package x509

import (
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"

	"github.com/nuts-foundation/nuts-auth/assets"
	"github.com/nuts-foundation/nuts-auth/pkg/contract"
	"github.com/nuts-foundation/nuts-auth/pkg/services"
)

// UziSignedToken implements a SignedToken interface for contracts signed by the UZI means in the JwtX509Token form.
type UziSignedToken struct {
	jwtX509Token *JwtX509Token
	contract     *contract.Contract
}

// UziValidator can check Uzi signed JWTs.
// It can parse and validate a UziSignedToken which implements the SignedToken interface
type UziValidator struct {
	validator         *JwtX509Validator
	contractTemplates *contract.TemplateStore
}

// UziEnv is used to indicate which Uzi environment should be used.
type UziEnv string

// UziProduction uses the production certificate tree:
// https://www.zorgcsp.nl/ca-certificaten
const UziProduction UziEnv = "production"

// UziAcceptation uses the acceptation certificate tree:
// https://acceptatie.zorgcsp.nl/ca-certificaten
const UziAcceptation UziEnv = "acceptation"

func getUziAttributeNames() []string {
	// A list of Uzi attribute names used to sign the message
	// See table 12 on page 62 of the Certification Practice Statement (CPS) UZI-register v10.x
	// https://zorgcsp.nl/Media/Default/documenten/2020-05-06_RK1%20CPS%20UZI-register%20V10.0.pdf
	return []string{
		"oidCa",
		"version",
		"uziNr",
		"cardType",
		"orgID",
		"rollCode",
		"agbCode",
	}
}

// SignerAttributes returns the attributes from the Uzi card used in the signature.
// For more information on these attributes, see table 12 on page 62 of the Certification Practice Statement (CPS) UZI-register v10.x
// https://zorgcsp.nl/Media/Default/documenten/2020-05-06_RK1%20CPS%20UZI-register%20V10.0.pdf
func (t UziSignedToken) SignerAttributes() (map[string]string, error) {
	res := map[string]string{}
	otherNames, err := t.jwtX509Token.SubjectAltNameOtherNames()
	if err != nil {
		return nil, fmt.Errorf("could not extract SAN from certificate: %w", err)
	}

	for _, otherNameStr := range otherNames {
		parts := strings.Split(otherNameStr, "-")
		attrNames := getUziAttributeNames()
		if len(parts) != len(attrNames) {
			continue
		}

		for idx, name := range getUziAttributeNames() {
			res[name] = parts[idx]
		}
	}
	return res, nil
}

// Contract returns the Contract signed by the Uzi means
func (t UziSignedToken) Contract() contract.Contract {
	return *t.contract
}

// certsFromAssets allows for easy loading of the used UziCertificates.
// These certs are compiled into the binary for easy distribution using go-bindata.
func certsFromAssets(paths []string) (certs []*x509.Certificate, err error) {
	for _, path := range paths {
		var (
			rawCert []byte
			cert    *x509.Certificate
		)
		rawCert, err = assets.Asset(path)
		if err != nil {
			return
		}

		cert, err = x509.ParseCertificate(rawCert)
		if err != nil {
			return
		}
		certs = append(certs, cert)
	}
	return
}

func validUziSigningAlgs() []jwa.SignatureAlgorithm {
	return []jwa.SignatureAlgorithm{jwa.RS256, jwa.RS512}
}

// NewUziValidator creates a new UziValidator.
// It accepts a UziEnv and preloads corresponding certificate tree.
// It accepts a contract template store which is used to check if the signed contract exists and is valid.
// It accepts an optional CrlGetter. If non is given, the CachedHttpCrlService is used as default
func NewUziValidator(env UziEnv, contractTemplates *contract.TemplateStore, crls CrlGetter) (validator *UziValidator, err error) {
	var roots []*x509.Certificate
	var intermediates []*x509.Certificate

	if env == UziProduction {
		roots, err = certsFromAssets([]string{
			"certs/uzi-prod/RootCA-G3.cer",
		})
		if err != nil {
			return
		}

		intermediates, err = certsFromAssets([]string{
			"certs/uzi-prod/20190418_UZI-register_Medewerker_op_naam_CA_G3.cer",
			"certs/uzi-prod/20190418_UZI-register_Zorgverlener_CA_G3.cer",
			"certs/uzi-prod/DomOrganisatiePersoonCA-G3.cer",
			"certs/uzi-prod/UZI-register_Medewerker_op_naam_CA_G3.cer",
			"certs/uzi-prod/UZI-register_Zorgverlener_CA_G3.cer",
		})
		if err != nil {
			return
		}
	} else if env == UziAcceptation {
		roots, err = certsFromAssets([]string{
			"certs/uzi-acc/test_zorg_csp_root_ca_g3.cer",
		})
		if err != nil {
			return
		}

		intermediates, err = certsFromAssets([]string{
			"certs/uzi-acc/test_uzi-register_medewerker_op_naam_ca_g3.cer",
			"certs/uzi-acc/test_zorg_csp_level_2_persoon_ca_g3.cer",
		})
		if err != nil {
			return
		}
	} else {
		return nil, fmt.Errorf("unknown uzi environment: %s", env)
	}

	if crls == nil {
		crls = NewCachedHttpCrlService()
	}

	validator = &UziValidator{
		validator:         NewJwtX509Validator(roots, intermediates, validUziSigningAlgs(), crls),
		contractTemplates: contractTemplates,
	}
	return
}

// Parse tries to parse a string as a base64 encoded JWT.
// The jwt should contain at leas one certificate in the x509 header
// It tries to find the contract in the given contractStore.
// No other verifications are performed.
// Make sure to call Verify to perform the actual crypto verifications
func (u UziValidator) Parse(rawAuthToken string) (services.SignedToken, error) {
	x509Token, err := u.validator.Parse(rawAuthToken)
	if err != nil {
		return nil, err
	}
	tokenField, ok := x509Token.token.Get("message")
	if !ok {
		return nil, fmt.Errorf("jwt did not contain token field")
	}
	contractText, ok := tokenField.(string)
	if !ok {
		return nil, fmt.Errorf("token field should contain a string")
	}

	c, err := contract.ParseContractString(contractText, *u.contractTemplates)
	if err != nil {
		return nil, err
	}

	return UziSignedToken{jwtX509Token: x509Token, contract: c}, nil
}

// extKeyUsageDocumentSigning is required for signing documents, according to the UZI spec.
var extKeyUsageDocumentSigning = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 10, 3, 12}

// Verify performs all the crypto verifications like:
// Correct hashing algorithm
// Correct certificate tree
// Certificates are not revoked
// Verifies all the extra jwt fields like exp, iat and nbf.
// Verifies if the signer attributes are valid
func (u UziValidator) Verify(token services.SignedToken) error {
	x509SignedToken, ok := token.(UziSignedToken)
	if !ok {
		return fmt.Errorf("wrong token type")
	}
	_, err := token.SignerAttributes()
	if err != nil {
		return fmt.Errorf("invalid signer attributes in uzi certificate: %w", err)
	}
	err = u.validator.Verify(x509SignedToken.jwtX509Token)
	if err != nil {
		return err
	}

	// check if the certificate has the Extended Key usage for document signing
	keyUsageFound := false
	for _, keyUsage := range x509SignedToken.jwtX509Token.chain[0].UnknownExtKeyUsage {
		if keyUsage.Equal(extKeyUsageDocumentSigning) {
			keyUsageFound = true
			break
		}
	}
	if !keyUsageFound {
		return fmt.Errorf("certificate is missing the extended key usage for document signing (%s)", extKeyUsageDocumentSigning.String())
	}
	return nil
}
