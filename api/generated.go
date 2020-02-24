// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

// AccessTokenRequestFailedResponse defines model for AccessTokenRequestFailedResponse.
type AccessTokenRequestFailedResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// AccessTokenRequestJWT defines model for AccessTokenRequestJWT.
type AccessTokenRequestJWT struct {
	Aud string  `json:"aud"`
	Con *string `json:"con,omitempty"`
	Exp float32 `json:"exp"`
	Iat float32 `json:"iat"`
	Iss string  `json:"iss"`
	Jti string  `json:"jti"`
	Osi *string `json:"osi,omitempty"`
	Sid string  `json:"sid"`
	Sub string  `json:"sub"`
	Uid string  `json:"uid"`
}

// AccessTokenResponse defines model for AccessTokenResponse.
type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float32 `json:"expires_in"`
	TokenType   string  `json:"token_type"`
}

// Contract defines model for Contract.
type Contract struct {
	Language           Language  `json:"language"`
	SignerAttributes   *[]string `json:"signer_attributes,omitempty"`
	Template           *string   `json:"template,omitempty"`
	TemplateAttributes *[]string `json:"template_attributes,omitempty"`
	Type               Type      `json:"type"`
	Version            Version   `json:"version"`
}

// ContractSigningRequest defines model for ContractSigningRequest.
type ContractSigningRequest struct {
	Language    Language    `json:"language"`
	LegalEntity LegalEntity `json:"legalEntity"`
	Type        Type        `json:"type"`
	ValidFrom   *string     `json:"valid_from,omitempty"`
	ValidTo     *string     `json:"valid_to,omitempty"`
	Version     Version     `json:"version"`
}

// CreateAccessTokenRequest defines model for CreateAccessTokenRequest.
type CreateAccessTokenRequest struct {
	Assertion string `json:"assertion"`
	GrantType string `json:"grant_type"`
}

// CreateJwtBearerTokenRequest defines model for CreateJwtBearerTokenRequest.
type CreateJwtBearerTokenRequest struct {
	Actor     string `json:"actor"`
	Custodian string `json:"custodian"`
	Identity  string `json:"identity"`
	Scope     string `json:"scope"`
	Subject   string `json:"subject"`
}

// CreateSessionResult defines model for CreateSessionResult.
type CreateSessionResult struct {
	QrCodeInfo IrmaQR `json:"qr_code_info"`
	SessionId  string `json:"session_id"`
}

// DisclosedAttribute defines model for DisclosedAttribute.
type DisclosedAttribute struct {
	Identifier string                   `json:"identifier"`
	Rawvalue   *string                  `json:"rawvalue,omitempty"`
	Status     string                   `json:"status"`
	Value      DisclosedAttribute_Value `json:"value"`
}

// DisclosedAttribute_Value defines model for DisclosedAttribute.Value.
type DisclosedAttribute_Value struct {
	AdditionalProperties map[string]string `json:"-"`
}

// DisclosedAttributeIndex defines model for DisclosedAttributeIndex.
type DisclosedAttributeIndex struct {
	Attr *int `json:"attr,omitempty"`
	Cred *int `json:"cred,omitempty"`
}

// ErrorString defines model for ErrorString.
type ErrorString string

// IrmaQR defines model for IrmaQR.
type IrmaQR struct {
	Irmaqr string `json:"irmaqr"`
	U      string `json:"u"`
}

// JwtBearerTokenResponse defines model for JwtBearerTokenResponse.
type JwtBearerTokenResponse struct {
	BearerToken string `json:"bearer_token"`
}

// Language defines model for Language.
type Language string

// LegalEntity defines model for LegalEntity.
type LegalEntity string

// Proof defines model for Proof.
type Proof interface{}

// ProofD defines model for ProofD.
type ProofD struct {
	A          *float32           `json:"A,omitempty"`
	ADisclosed *ProofD_ADisclosed `json:"a_disclosed,omitempty"`
	AResponses *ProofD_AResponses `json:"a_responses,omitempty"`
	C          *float32           `json:"c,omitempty"`
	EResponse  *float32           `json:"e_response,omitempty"`
	VResponse  *float32           `json:"v_response,omitempty"`
}

// ProofD_ADisclosed defines model for ProofD.ADisclosed.
type ProofD_ADisclosed struct {
	AdditionalProperties map[string]float32 `json:"-"`
}

// ProofD_AResponses defines model for ProofD.AResponses.
type ProofD_AResponses struct {
	AdditionalProperties map[string]float32 `json:"-"`
}

// ProofP defines model for ProofP.
type ProofP struct {
	P         *float32 `json:"P,omitempty"`
	C         *float32 `json:"c,omitempty"`
	SResponse *float32 `json:"s_response,omitempty"`
}

// ProofS defines model for ProofS.
type ProofS struct {
	C         *float32 `json:"c,omitempty"`
	EResponse *float32 `json:"e_response,omitempty"`
}

// ProofU defines model for ProofU.
type ProofU struct {
	U              *float32 `json:"U,omitempty"`
	C              *float32 `json:"c,omitempty"`
	SResponse      *float32 `json:"s_response,omitempty"`
	VPrimeResponse *float32 `json:"v_prime_response,omitempty"`
}

// RemoteError defines model for RemoteError.
type RemoteError struct {
	Description *string `json:"description,omitempty"`
	Error       *string `json:"error,omitempty"`
	Message     *string `json:"message,omitempty"`
	Stacktrace  *string `json:"stacktrace,omitempty"`
	Status      *int    `json:"status,omitempty"`
}

// SessionResult defines model for SessionResult.
type SessionResult struct {
	Disclosed     *[]DisclosedAttribute `json:"disclosed,omitempty"`
	Error         *RemoteError          `json:"error,omitempty"`
	NutsAuthToken *string               `json:"nuts_auth_token,omitempty"`
	ProofStatus   *string               `json:"proofStatus,omitempty"`
	Signature     *SignedMessage        `json:"signature,omitempty"`
	Status        string                `json:"status"`
	Token         string                `json:"token"`
	Type          string                `json:"type"`
}

// SignedMessage defines model for SignedMessage.
type SignedMessage struct {
	Context   *float32                     `json:"context,omitempty"`
	Indices   *[][]DisclosedAttributeIndex `json:"indices,omitempty"`
	Message   *string                      `json:"message,omitempty"`
	Nonce     *float32                     `json:"nonce,omitempty"`
	Signature *[]Proof                     `json:"signature,omitempty"`
	Timestamp *Timestamp                   `json:"timestamp,omitempty"`
}

// Timestamp defines model for Timestamp.
type Timestamp struct {
	Time *int64 `json:"time,omitempty"`
}

// TokenIntrospectionRequest defines model for TokenIntrospectionRequest.
type TokenIntrospectionRequest struct {
	Token string `json:"token"`
}

// TokenIntrospectionResponse defines model for TokenIntrospectionResponse.
type TokenIntrospectionResponse struct {
	Active     bool    `json:"active"`
	Aud        *string `json:"aud,omitempty"`
	Email      *string `json:"email,omitempty"`
	Exp        *int    `json:"exp,omitempty"`
	FamilyName *string `json:"family_name,omitempty"`
	GivenName  *string `json:"given_name,omitempty"`
	Iat        *int    `json:"iat,omitempty"`
	Iss        *string `json:"iss,omitempty"`
	Name       *string `json:"name,omitempty"`
	Prefix     *string `json:"prefix,omitempty"`
	Scope      *string `json:"scope,omitempty"`
	Sid        *string `json:"sid,omitempty"`
	Sub        *string `json:"sub,omitempty"`
	Uid        *string `json:"uid,omitempty"`
}

// Type defines model for Type.
type Type string

// ValidationRequest defines model for ValidationRequest.
type ValidationRequest struct {
	ActingPartyCn  string `json:"acting_party_cn"`
	ContractFormat string `json:"contract_format"`
	ContractString string `json:"contract_string"`
}

// ValidationResult defines model for ValidationResult.
type ValidationResult struct {
	ContractFormat   string                            `json:"contract_format"`
	SignerAttributes ValidationResult_SignerAttributes `json:"signer_attributes"`
	ValidationResult string                            `json:"validation_result"`
}

// ValidationResult_SignerAttributes defines model for ValidationResult.SignerAttributes.
type ValidationResult_SignerAttributes struct {
	AdditionalProperties map[string]string `json:"-"`
}

// Version defines model for Version.
type Version string

// createSessionJSONBody defines parameters for CreateSession.
type createSessionJSONBody ContractSigningRequest

// validateContractJSONBody defines parameters for ValidateContract.
type validateContractJSONBody ValidationRequest

// GetContractByTypeParams defines parameters for GetContractByType.
type GetContractByTypeParams struct {

	// The version of this contract. If omitted, the most recent version will be returned
	Version  *string `json:"version,omitempty"`
	Language *string `json:"language,omitempty"`
}

// createJwtBearerTokenJSONBody defines parameters for CreateJwtBearerToken.
type createJwtBearerTokenJSONBody CreateJwtBearerTokenRequest

// CreateSessionRequestBody defines body for CreateSession for application/json ContentType.
type CreateSessionJSONRequestBody createSessionJSONBody

// ValidateContractRequestBody defines body for ValidateContract for application/json ContentType.
type ValidateContractJSONRequestBody validateContractJSONBody

// CreateJwtBearerTokenRequestBody defines body for CreateJwtBearerToken for application/json ContentType.
type CreateJwtBearerTokenJSONRequestBody createJwtBearerTokenJSONBody

// Getter for additional properties for DisclosedAttribute_Value. Returns the specified
// element and whether it was found
func (a DisclosedAttribute_Value) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for DisclosedAttribute_Value
func (a *DisclosedAttribute_Value) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for DisclosedAttribute_Value to handle AdditionalProperties
func (a *DisclosedAttribute_Value) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for DisclosedAttribute_Value to handle AdditionalProperties
func (a DisclosedAttribute_Value) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ProofD_ADisclosed. Returns the specified
// element and whether it was found
func (a ProofD_ADisclosed) Get(fieldName string) (value float32, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ProofD_ADisclosed
func (a *ProofD_ADisclosed) Set(fieldName string, value float32) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]float32)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ProofD_ADisclosed to handle AdditionalProperties
func (a *ProofD_ADisclosed) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]float32)
		for fieldName, fieldBuf := range object {
			var fieldVal float32
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ProofD_ADisclosed to handle AdditionalProperties
func (a ProofD_ADisclosed) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ProofD_AResponses. Returns the specified
// element and whether it was found
func (a ProofD_AResponses) Get(fieldName string) (value float32, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ProofD_AResponses
func (a *ProofD_AResponses) Set(fieldName string, value float32) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]float32)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ProofD_AResponses to handle AdditionalProperties
func (a *ProofD_AResponses) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]float32)
		for fieldName, fieldBuf := range object {
			var fieldVal float32
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ProofD_AResponses to handle AdditionalProperties
func (a ProofD_AResponses) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ValidationResult_SignerAttributes. Returns the specified
// element and whether it was found
func (a ValidationResult_SignerAttributes) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ValidationResult_SignerAttributes
func (a *ValidationResult_SignerAttributes) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ValidationResult_SignerAttributes to handle AdditionalProperties
func (a *ValidationResult_SignerAttributes) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ValidationResult_SignerAttributes to handle AdditionalProperties
func (a ValidationResult_SignerAttributes) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create an access token based on the OAuth JWT Bearer flow.// (POST /auth/accesstoken)
	CreateAccessToken(ctx echo.Context) error
	// CreateSessionHandler Initiates an IRMA signing session with the correct contract.// (POST /auth/contract/session)
	CreateSession(ctx echo.Context) error
	// returns the result of the contract request// (GET /auth/contract/session/{id})
	SessionRequestStatus(ctx echo.Context, id string) error
	// Validate a Nuts Security Contract// (POST /auth/contract/validate)
	ValidateContract(ctx echo.Context) error
	// Get a contract by type and version// (GET /auth/contract/{contractType})
	GetContractByType(ctx echo.Context, contractType string, params GetContractByTypeParams) error
	// Create a JWT Bearer Token which can be used in the createAccessToken request in the assertion field// (POST /auth/jwtbearertoken)
	CreateJwtBearerToken(ctx echo.Context) error
	// Introspection endpoint to retrieve information from an Access Token as described by RFC7662// (POST /auth/token_introspection)
	IntrospectAccessToken(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateAccessToken converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAccessToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateAccessToken(ctx)
	return err
}

// CreateSession converts echo context to params.
func (w *ServerInterfaceWrapper) CreateSession(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateSession(ctx)
	return err
}

// SessionRequestStatus converts echo context to params.
func (w *ServerInterfaceWrapper) SessionRequestStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SessionRequestStatus(ctx, id)
	return err
}

// ValidateContract converts echo context to params.
func (w *ServerInterfaceWrapper) ValidateContract(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ValidateContract(ctx)
	return err
}

// GetContractByType converts echo context to params.
func (w *ServerInterfaceWrapper) GetContractByType(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "contractType" -------------
	var contractType string

	err = runtime.BindStyledParameter("simple", false, "contractType", ctx.Param("contractType"), &contractType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter contractType: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetContractByTypeParams
	// ------------- Optional query parameter "version" -------------
	if paramValue := ctx.QueryParam("version"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "version", ctx.QueryParams(), &params.Version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter version: %s", err))
	}

	// ------------- Optional query parameter "language" -------------
	if paramValue := ctx.QueryParam("language"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "language", ctx.QueryParams(), &params.Language)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter language: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetContractByType(ctx, contractType, params)
	return err
}

// CreateJwtBearerToken converts echo context to params.
func (w *ServerInterfaceWrapper) CreateJwtBearerToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateJwtBearerToken(ctx)
	return err
}

// IntrospectAccessToken converts echo context to params.
func (w *ServerInterfaceWrapper) IntrospectAccessToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.IntrospectAccessToken(ctx)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/auth/accesstoken", wrapper.CreateAccessToken)
	router.POST("/auth/contract/session", wrapper.CreateSession)
	router.GET("/auth/contract/session/:id", wrapper.SessionRequestStatus)
	router.POST("/auth/contract/validate", wrapper.ValidateContract)
	router.GET("/auth/contract/:contractType", wrapper.GetContractByType)
	router.POST("/auth/jwtbearertoken", wrapper.CreateJwtBearerToken)
	router.POST("/auth/token_introspection", wrapper.IntrospectAccessToken)

}

