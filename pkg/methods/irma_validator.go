package methods

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nuts-foundation/nuts-auth/pkg/contract"

	types "github.com/nuts-foundation/nuts-auth/pkg/types"

	core "github.com/nuts-foundation/nuts-go-core"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
	nutscrypto "github.com/nuts-foundation/nuts-crypto/pkg"
	cryptoTypes "github.com/nuts-foundation/nuts-crypto/pkg/types"
	registry "github.com/nuts-foundation/nuts-registry/pkg"
	irma "github.com/privacybydesign/irmago"
	irmaserver "github.com/privacybydesign/irmago/server"
	"github.com/sirupsen/logrus"
)

// IrmaValidator validates contracts using the irma logic.
type IrmaValidator struct {
	IrmaServer     types.IrmaServerClient
	IrmaConfig     *irma.Configuration
	Registry       registry.RegistryClient
	Crypto         nutscrypto.Client
	ValidContracts contract.ContractMatrix
}

// LegacyIdentityToken is the JWT that was used as Identity token in versions prior to < 0.13
type LegacyIdentityToken struct {
	jwt.StandardClaims
	Contract SignedIrmaContract `json:"nuts_signature"`
}

// IsInitialized is a helper function to determine if the validator has been initialized properly.
func (v IrmaValidator) IsInitialized() bool {
	return v.IrmaConfig != nil
}

// ValidateContract is the entry point for contract validation.
// It decodes the base64 encoded contract, parses the contract string, and validates the contract.
// Returns nil, ErrUnknownContractFormat if the contract used in the message is unknown
func (v IrmaValidator) ValidateContract(b64EncodedContract string, format types.ContractFormat, actingPartyCN string) (*types.ContractValidationResult, error) {
	if format == types.IrmaFormat {
		contract, err := base64.StdEncoding.DecodeString(b64EncodedContract)
		if err != nil {
			return nil, fmt.Errorf("could not base64-decode contract: %w", err)
		}
		// Create the irma contract validator
		contractValidator := IrmaContractVerifier{v.IrmaConfig, v.ValidContracts}
		signedContract, err := contractValidator.ParseSignedIrmaContract(string(contract))
		if err != nil {
			return nil, err
		}
		return contractValidator.VerifyAll(signedContract, actingPartyCN)
	}
	return nil, types.ErrUnknownContractFormat
}

// ValidateJwt validates a JWT formatted identity token
func (v IrmaValidator) ValidateJwt(token string, actingPartyCN string) (*types.ContractValidationResult, error) {
	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Name}}
	parsedToken, err := parser.ParseWithClaims(token, &types.NutsIdentityToken{}, func(token *jwt.Token) (i interface{}, e error) {
		legalEntity, err := parseTokenIssuer(token.Claims.(*types.NutsIdentityToken).Issuer)
		if err != nil {
			return nil, err
		}

		// get public key
		org, err := v.Registry.OrganizationById(legalEntity)
		if err != nil {
			return nil, err
		}

		pk, err := org.CurrentPublicKey()

		if err != nil {
			return nil, err
		}

		return pk.Materialize()
	})

	if err != nil {
		return nil, err
	}

	claims := parsedToken.Claims.(*types.NutsIdentityToken)

	if claims.Type != types.IrmaFormat {
		return nil, fmt.Errorf("%s: %w", claims.Type, types.ErrInvalidContractFormat)
	}

	contractStr, err := base64.StdEncoding.DecodeString(claims.Signature)
	if err != nil {
		return nil, err
	}

	// Create the irma contract validator
	contractValidator := IrmaContractVerifier{v.IrmaConfig, v.ValidContracts}
	signedContract, err := contractValidator.ParseSignedIrmaContract(string(contractStr))
	return contractValidator.VerifyAll(signedContract, actingPartyCN)
}

// SessionStatus returns the current status of a certain session.
// It returns nil if the session is not found
func (v IrmaValidator) SessionStatus(id types.SessionID) (*types.SessionStatusResult, error) {
	if result := v.IrmaServer.GetSessionResult(string(id)); result != nil {
		var (
			token string
		)
		if result.Signature != nil {
			contractTemplate, err := contract.NewContractFromMessageContents(result.Signature.Message, v.ValidContracts)
			sic := &SignedIrmaContract{*result.Signature, contractTemplate}
			if err != nil {
				return nil, err
			}

			le, err := v.legalEntityFromContract(sic)
			if err != nil {
				return nil, fmt.Errorf("could not create JWT for given session: %w", err)
			}

			token, err = v.CreateIdentityTokenFromIrmaContract(sic, le)
			if err != nil {
				return nil, err
			}
		}
		result := &types.SessionStatusResult{*result, token}
		logrus.Info(result.NutsAuthToken)
		return result, nil
	}
	return nil, types.ErrSessionNotFound
}

func (v IrmaValidator) legalEntityFromContract(sic *SignedIrmaContract) (core.PartyID, error) {
	params, err := sic.ContractTemplate.ExtractParams(sic.IrmaContract.Message)
	if err != nil {
		return core.PartyID{}, err
	}

	if _, ok := params["legal_entity"]; !ok {
		return core.PartyID{}, types.ErrLegalEntityNotProvided
	}

	le, err := v.Registry.ReverseLookup(params["legal_entity"])
	if err != nil {
		return core.PartyID{}, err
	}

	return le.Identifier, nil
}

// CreateIdentityTokenFromIrmaContract from a signed irma contract. Returns a JWT signed with the provided legalEntity.
func (v IrmaValidator) CreateIdentityTokenFromIrmaContract(contract *SignedIrmaContract, legalEntity core.PartyID) (string, error) {
	signature, err := json.Marshal(contract.IrmaContract)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	if err != nil {
		return "", err
	}
	payload := types.NutsIdentityToken{
		StandardClaims: jwt.StandardClaims{
			Issuer: legalEntity.String(),
		},
		Signature: encodedSignature,
		Type:      types.IrmaFormat,
	}

	claims, err := convertPayloadToClaims(payload)
	if err != nil {
		return "", fmt.Errorf("could not construct claims: %w", err)
	}

	tokenString, err := v.Crypto.SignJWT(claims, cryptoTypes.KeyForEntity(cryptoTypes.LegalEntity{URI: legalEntity.String()}))
	if err != nil {
		return "", fmt.Errorf("could not sign jwt: %w", err)
	}
	return tokenString, nil
}

// func for creating legacy token
func (v IrmaValidator) createLegacyIdentityToken(contract *SignedIrmaContract, legalEntity string) (string, error) {
	payload := LegacyIdentityToken{
		StandardClaims: jwt.StandardClaims{
			Issuer:  "nuts",
			Subject: legalEntity,
		},
		Contract: *contract,
	}

	claims, err := convertPayloadToClaimsLegacy(payload)
	if err != nil {
		err = fmt.Errorf("could not construct claims: %w", err)
		logrus.Error(err)
		return "", err
	}

	tokenString, err := v.Crypto.SignJWT(claims, cryptoTypes.KeyForEntity(cryptoTypes.LegalEntity{URI: legalEntity}))
	if err != nil {
		err = fmt.Errorf("could not sign jwt: %w", err)
		logrus.Error(err)
		return "", err
	}
	return tokenString, nil
}

func convertPayloadToClaimsLegacy(payload LegacyIdentityToken) (map[string]interface{}, error) {

	var (
		jsonString []byte
		err        error
		claims     map[string]interface{}
	)

	if jsonString, err = json.Marshal(payload); err != nil {
		return nil, fmt.Errorf("could not marshall payload: %w", err)
	}

	if err := json.Unmarshal(jsonString, &claims); err != nil {
		return nil, fmt.Errorf("could not unmarshall string: %w", err)
	}

	return claims, nil
}

// convertPayloadToClaims converts a nutsJwt struct to a map of strings so it can be signed with the crypto module
func convertPayloadToClaims(payload types.NutsIdentityToken) (map[string]interface{}, error) {

	var (
		jsonString []byte
		err        error
		claims     map[string]interface{}
	)

	if jsonString, err = json.Marshal(payload); err != nil {
		return nil, fmt.Errorf("could not marshall payload: %w", err)
	}

	if err := json.Unmarshal(jsonString, &claims); err != nil {
		return nil, fmt.Errorf("could not unmarshall string: %w", err)
	}

	return claims, nil
}

// StartSession starts an irma session.
// This is mainly a wrapper around the irma.IrmaServer.StartSession
func (v IrmaValidator) StartSession(request interface{}, handler irmaserver.SessionHandler) (*irma.Qr, string, error) {
	return v.IrmaServer.StartSession(request, handler)
}

// ParseAndValidateJwtBearerToken validates the jwt signature and returns the containing claims
func (v IrmaValidator) ParseAndValidateJwtBearerToken(acString string) (*types.NutsJwtBearerToken, error) {
	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Name}}
	token, err := parser.ParseWithClaims(acString, &types.NutsJwtBearerToken{}, func(token *jwt.Token) (i interface{}, e error) {
		legalEntity, err := parseTokenIssuer(token.Claims.(*types.NutsJwtBearerToken).Issuer)
		if err != nil {
			return nil, err
		}

		// get public key
		org, err := v.Registry.OrganizationById(legalEntity)
		if err != nil {
			return nil, err
		}

		pk, err := org.CurrentPublicKey()
		if err != nil {
			return nil, err
		}

		return pk.Materialize()
	})

	if token != nil && token.Valid {
		if claims, ok := token.Claims.(*types.NutsJwtBearerToken); ok {
			return claims, nil
		}
	}

	return nil, err
}

// BuildAccessToken builds an access token based on the oauth claims and the identity of the user provided by the identityValidationResult
// The token gets signed with the custodians private key and returned as a string.
func (v IrmaValidator) BuildAccessToken(jwtBearerToken *types.NutsJwtBearerToken, identityValidationResult *types.ContractValidationResult) (string, error) {

	if identityValidationResult.ValidationResult != types.Valid {
		return "", fmt.Errorf("could not build accessToken: %w", errors.New("invalid contract"))
	}

	issuer := jwtBearerToken.Subject
	if issuer == "" {
		return "", fmt.Errorf("could not build accessToken: %w", errors.New("subject is missing"))
	}

	at := types.NutsAccessToken{
		StandardClaims: jwt.StandardClaims{
			// Expires in 15 minutes
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   jwtBearerToken.Issuer,
		},
		SubjectID: jwtBearerToken.SubjectID,
		Scope:     jwtBearerToken.Scope,
		// based on
		// https://privacybydesign.foundation/attribute-index/en/pbdf.gemeente.personalData.html
		// https://privacybydesign.foundation/attribute-index/en/pbdf.pbdf.email.html
		// and
		// https://openid.net/specs/openid-connect-basic-1_0.html#StandardClaims
		FamilyName: identityValidationResult.DisclosedAttributes["gemeente.personalData.familyname"],
		GivenName:  identityValidationResult.DisclosedAttributes["gemeente.personalData.firstnames"],
		Prefix:     identityValidationResult.DisclosedAttributes["gemeente.personalData.prefix"],
		Name:       identityValidationResult.DisclosedAttributes["gemeente.personalData.fullname"],
		Email:      identityValidationResult.DisclosedAttributes["pbdf.email.email"],
	}

	var keyVals map[string]interface{}
	inrec, _ := json.Marshal(at)
	if err := json.Unmarshal(inrec, &keyVals); err != nil {
		return "", err
	}

	// Sign with the private key of the issuer
	token, err := v.Crypto.SignJWT(keyVals, cryptoTypes.KeyForEntity(cryptoTypes.LegalEntity{URI: issuer}))
	if err != nil {
		return token, fmt.Errorf("could not build accessToken: %w", err)
	}

	return token, err
}

// CreateJwtBearerToken creates a JwtBearerTokenResponse containing a jwtBearerToken from a CreateJwtBearerTokenRequest.
func (v IrmaValidator) CreateJwtBearerToken(request *types.CreateJwtBearerTokenRequest) (*types.JwtBearerTokenResponse, error) {
	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	jwtBearerToken := types.NutsJwtBearerToken{
		StandardClaims: jwt.StandardClaims{
			//Audience:  endpoint.Identifier.String(),
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
			Id:        jti.String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    request.Actor,
			NotBefore: 0,
			Subject:   request.Custodian,
		},
		IdentityToken: request.IdentityToken,
		SubjectID:     request.Subject,
		Scope:         request.Scope,
	}

	var keyVals map[string]interface{}
	inrec, _ := json.Marshal(jwtBearerToken)
	if err := json.Unmarshal(inrec, &keyVals); err != nil {
		return nil, err
	}

	signingString, err := v.Crypto.SignJWT(keyVals, cryptoTypes.KeyForEntity(cryptoTypes.LegalEntity{URI: request.Actor}))
	if err != nil {
		return nil, err
	}

	return &types.JwtBearerTokenResponse{BearerToken: signingString}, nil
}

// ParseAndValidateAccessToken parses and validates a accesstoken string and returns a filled in NutsAccessToken.
func (v IrmaValidator) ParseAndValidateAccessToken(accessToken string) (*types.NutsAccessToken, error) {
	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Name}}
	token, err := parser.ParseWithClaims(accessToken, &types.NutsAccessToken{}, func(token *jwt.Token) (i interface{}, e error) {
		legalEntity, err := parseTokenIssuer(token.Claims.(*types.NutsAccessToken).Issuer)
		if err != nil {
			return nil, err
		}

		// Check if the care provider which signed the token is managed by this node
		if !v.Crypto.PrivateKeyExists(cryptoTypes.KeyForEntity(cryptoTypes.LegalEntity{URI: legalEntity.String()})) {
			return nil, fmt.Errorf("invalid token: not signed by a care provider of this node")
		}

		// get public key
		org, err := v.Registry.OrganizationById(legalEntity)
		if err != nil {
			return nil, err
		}

		pk, err := org.CurrentPublicKey()
		if err != nil {
			return nil, err
		}

		return pk.Materialize()
	})

	if token != nil && token.Valid {
		if claims, ok := token.Claims.(*types.NutsAccessToken); ok {
			return claims, nil
		}
	}
	return nil, err
}

func parseTokenIssuer(issuer string) (core.PartyID, error) {
	if issuer == "" {
		return core.PartyID{}, types.ErrLegalEntityNotProvided
	}
	if result, err := core.ParsePartyID(issuer); err != nil {
		return core.PartyID{}, fmt.Errorf("invalid token issuer: %w", err)
	} else {
		return result, nil
	}
}