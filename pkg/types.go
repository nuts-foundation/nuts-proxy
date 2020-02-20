package pkg

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	irma "github.com/privacybydesign/irmago"
	"github.com/privacybydesign/irmago/server"
	"github.com/privacybydesign/irmago/server/irmaserver"
)

// SsoEndpointType defines the string to use in the endpoint for the AuthenticationServer
const SsoEndpointType string = "urn:ietf:rfc:3986:urn:oid:1.3.6.1.4.1.54851.1:nuts-oauth-authentication-server"

// JwtBearerGrantType defines the grant-type to use in the access token request
const JwtBearerGrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"

// SessionID contains a number to uniquely identify a contract signing session
type SessionID string

// ContractFormat describes the format of a signed contract. Based on the format an appropriate validator can be selected.
type ContractFormat string

// ValidationState contains the outcome of the validation. It van be VALID or INVALID. This makes it human readable.
type ValidationState string

// ErrUnknownContractFormat is returned when the contract format is unknown
var ErrUnknownContractFormat = errors.New("unknown contract format")

// ErrSessionNotFound is returned when there is no contract signing session found for a certain SessionID
var ErrSessionNotFound = errors.New("session not found")

// ContractValidator interface must be implemented by contract validators
type ContractValidator interface {
	ValidateContract(contract string, format ContractFormat, actingPartyCN string) (*ContractValidationResult, error)
	ValidateJwt(contract string, actingPartyCN string) (*ContractValidationResult, error)
	IsInitialized() bool
}

// ContractSessionHandler interface must be implemented by ContractSessionHandlers
type ContractSessionHandler interface {
	SessionStatus(session SessionID) (*SessionStatusResult, error)
	StartSession(request interface{}, handler irmaserver.SessionHandler) (*irma.Qr, string, error)
}

type AccessTokenHandler interface {
	// CreateJwtBearerToken from a JwtBearerTokenRequest. Returns a signed JWT string.
	CreateJwtBearerToken(request *CreateJwtBearerTokenRequest) (token *JwtBearerAccessTokenResponse, err error)

	ParseAndValidateJwtBearerToken(token string) (*NutsJwtBearerToken, error)

	BuildAccessToken(jwtClaims *NutsJwtBearerToken, identityValidationResult *ContractValidationResult) (token string, err error)

	ParseAndValidateAccessToken(accessToken string) (*NutsAccessToken, error)
}

const (
	// IrmaFormat is used to indicate a contract is in he form of a base64 encoded IRMA signature
	IrmaFormat ContractFormat = "irma"
	// JwtFormat is used to indicate a contract in in the form of a Jwt encoded signature
	JwtFormat ContractFormat = "JWT"
	// Valid is used to indicate a contract was valid on the time of testing
	Valid ValidationState = "VALID"
	// Invalid is used to indicate a contract was invalid on the time of testing
	Invalid ValidationState = "INVALID"
)

// CreateSessionRequest is used to create a contract signing session.
type CreateSessionRequest struct {
	// ContractType such as "BehandelaarLogin"
	Type ContractType
	// Version of the contract such as "v1"
	Version Version
	// Language of the contact such as "NL"
	Language Language
	// LegalEntity denotes the organization of the user
	LegalEntity string
	// ValidFrom describes the time from which this contract should be considered valid
	ValidFrom time.Time
	// ValidFrom describes the time until this contract should be considered valid
	ValidTo time.Time
	// TemplateAttributes is an object containing extra template values. example: {"reason":"providing care"}
	TemplateAttributes map[string]string
}

// CreateSessionResult contains the results needed to setup an irma flow
type CreateSessionResult struct {
	QrCodeInfo irma.Qr
	SessionID  string
}

// SessionStatusResult contains the current state of a session. If the session is DONE it also contains a JWT in the NutsAuthToken
type SessionStatusResult struct {
	server.SessionResult
	// NutsAuthToken contains the JWT if the sessionStatus is DONE
	NutsAuthToken string `json:"nuts_auth_token"`
}

// ValidationRequest is used to pass all information to ValidateContract
type ValidationRequest struct {
	// ContractFormat specifies the type of format used for the contract, e.g. 'irma'
	ContractFormat ContractFormat

	// The actual contract in string format to validate
	ContractString string

	// ActingPartyCN is the common name of the Acting party extracted from the client cert
	ActingPartyCN string
}

type CreateAccessTokenRequest struct {
	// JwtString the string containing the unmarshalled jwt.
	JwtString string
}

type CreateJwtBearerTokenRequest struct {
	Actor         string
	Custodian     string
	IdentityToken string
	Subject       string
	Scope         string
}

// NutsIdentityToken contains the signed identity of the user performing the request
type NutsIdentityToken struct {
	jwt.StandardClaims
	//Identifier of the legalEntity who issued and signed the token
	//Issuer string
	// What kind of signature? Currently only IRMA is supported
	Type ContractFormat `json:"type"`
	// The base64 encoded signature
	Signature string `json:"sig"`
	//Contract SignedIrmaContract `json:"nuts_signature"`
}

// NutsJwtBearerToken contains the deserialized Jwt Bearer Token as defined in rfc7523. It contains a NutsIdentity token which can be
// verified by the authorization server.
type NutsJwtBearerToken struct {
	jwt.StandardClaims
	Custodian     string `json:"custodian"`
	IdentityToken string `json:"usi"`
	SubjectId     string `json:"sid"`
	Scope         string `json:"scope"`
}

// NutsAccessToken is a OAuth 2.0 access token which provides context to a request.
// Its contents are derived from a Jwt Bearer token. The Jwt Bearer token is verified by the authorization server and
// stripped from the proof to make it compact.
type NutsAccessToken struct {
	jwt.StandardClaims
	SubjectId  string `json:"sid"`
	Scope      string `json:"scope"`
	FamilyName string `json:"family_name"`
	Prefix     string `json:"prefix"`
	Initials   string `json:"initials"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type RequestContext struct {
	Actor         string
	Custodian     string
	IdentityToken string
	Subject       string
}

type AccessTokenResponse struct {
	AccessToken string
}

type JwtBearerAccessTokenResponse struct {
	BearerToken string
}

// ContractValidationResult contains the result of a contract validation
type ContractValidationResult struct {
	ValidationResult ValidationState `json:"validation_result"`
	ContractFormat   ContractFormat  `json:"contract_format"`
	// DisclosedAttributes contain the attributes used to sign this contract
	DisclosedAttributes map[string]string `json:"disclosed_attributes"`
}

// IrmaServerClient is an abstraction for the Irma Server, mainly for enabling better testing
type IrmaServerClient interface {
	GetSessionResult(token string) *server.SessionResult
	StartSession(request interface{}, handler irmaserver.SessionHandler) (*irma.Qr, string, error)
}

// DefaultIrmaClient is a wrapper for the Irma Server
type DefaultIrmaClient struct {
	I *irmaserver.Server
}

// GetSessionResult forwards to Irma Server instance
func (d *DefaultIrmaClient) GetSessionResult(token string) *server.SessionResult {
	return d.I.GetSessionResult(token)
}

// StartSession forwards to Irma Server instance
func (d *DefaultIrmaClient) StartSession(request interface{}, handler irmaserver.SessionHandler) (*irma.Qr, string, error) {
	return d.I.StartSession(request, handler)
}
