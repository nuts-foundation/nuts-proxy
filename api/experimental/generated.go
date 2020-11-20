// Package experimental provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package experimental

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ContractLanguage defines model for ContractLanguage.
type ContractLanguage string

// ContractResponse defines model for ContractResponse.
type ContractResponse struct {

	// Language of the contract in all caps
	Language ContractLanguage `json:"language"`

	// The contract message
	Message string `json:"message"`

	// Type of which contract to sign
	Type ContractType `json:"type"`

	// Version of the contract
	Version ContractVersion `json:"version"`
}

// ContractTemplateResponse defines model for ContractTemplateResponse.
type ContractTemplateResponse struct {

	// Language of the contract in all caps
	Language ContractLanguage `json:"language"`
	Template string           `json:"template"`

	// Type of which contract to sign
	Type ContractType `json:"type"`

	// Version of the contract
	Version ContractVersion `json:"version"`
}

// ContractType defines model for ContractType.
type ContractType string

// ContractVersion defines model for ContractVersion.
type ContractVersion string

// CreateSignSessionRequest defines model for CreateSignSessionRequest.
type CreateSignSessionRequest struct {
	Means string `json:"means"`

	// Params are passed to the means. Should be documented in the means documentation.
	Params map[string]interface{} `json:"params"`

	// base64 encoded payload what needs to be signed
	Payload string `json:"payload"`
}

// CreateSignSessionResult defines model for CreateSignSessionResult.
type CreateSignSessionResult struct {

	// The means this session uses to sign.
	Means string `json:"means"`

	// A pointer to a signature session. This is an opaque value which only has meaning in the context of the signing means. Can be an URL, base64 encoded image of a QRCode etc.
	SessionPtr string `json:"sessionPtr"`
}

// DrawUpContractRequest defines model for DrawUpContractRequest.
type DrawUpContractRequest struct {

	// Language of the contract in all caps
	Language ContractLanguage `json:"language"`

	// Identifier of the legalEntity as registered in the Nuts registry
	LegalEntity LegalEntity `json:"legalEntity"`

	// Type of which contract to sign
	Type ContractType `json:"type"`

	// Version of the contract
	Version ContractVersion `json:"version"`
}

// LegalEntity defines model for LegalEntity.
type LegalEntity string

// DrawUpContractJSONBody defines parameters for DrawUpContract.
type DrawUpContractJSONBody DrawUpContractRequest

// GetContractTemplateParams defines parameters for GetContractTemplate.
type GetContractTemplateParams struct {

	// The version of this contract. If omitted, the most recent version will be returned
	Version *string `json:"version,omitempty"`
}

// CreateSignSessionJSONBody defines parameters for CreateSignSession.
type CreateSignSessionJSONBody CreateSignSessionRequest

// DrawUpContractRequestBody defines body for DrawUpContract for application/json ContentType.
type DrawUpContractJSONRequestBody DrawUpContractJSONBody

// CreateSignSessionRequestBody defines body for CreateSignSession for application/json ContentType.
type CreateSignSessionJSONRequestBody CreateSignSessionJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Draw up a contract using a specified contract template, language and version
	// (PUT /auth/internal/experimental/contract/drawup)
	DrawUpContract(ctx echo.Context) error
	// Get the contract template by version, and type
	// (GET /auth/internal/experimental/contract/template/{language}/{contractType})
	GetContractTemplate(ctx echo.Context, language string, contractType string, params GetContractTemplateParams) error
	// Create a signing session for a supported means.
	// (POST /auth/internal/experimental/sign)
	CreateSignSession(ctx echo.Context) error
	// Get the current status of a signing session
	// (GET /auth/internal/experimental/sign/{sessionPtr})
	GetSignSessionStatus(ctx echo.Context, sessionPtr string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// DrawUpContract converts echo context to params.
func (w *ServerInterfaceWrapper) DrawUpContract(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DrawUpContract(ctx)
	return err
}

// GetContractTemplate converts echo context to params.
func (w *ServerInterfaceWrapper) GetContractTemplate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "language" -------------
	var language string

	err = runtime.BindStyledParameter("simple", false, "language", ctx.Param("language"), &language)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter language: %s", err))
	}

	// ------------- Path parameter "contractType" -------------
	var contractType string

	err = runtime.BindStyledParameter("simple", false, "contractType", ctx.Param("contractType"), &contractType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter contractType: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetContractTemplateParams
	// ------------- Optional query parameter "version" -------------

	err = runtime.BindQueryParameter("form", true, false, "version", ctx.QueryParams(), &params.Version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter version: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetContractTemplate(ctx, language, contractType, params)
	return err
}

// CreateSignSession converts echo context to params.
func (w *ServerInterfaceWrapper) CreateSignSession(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateSignSession(ctx)
	return err
}

// GetSignSessionStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetSignSessionStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "sessionPtr" -------------
	var sessionPtr string

	err = runtime.BindStyledParameter("simple", false, "sessionPtr", ctx.Param("sessionPtr"), &sessionPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sessionPtr: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSignSessionStatus(ctx, sessionPtr)
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

	router.PUT("/auth/internal/experimental/contract/drawup", wrapper.DrawUpContract)
	router.GET("/auth/internal/experimental/contract/template/:language/:contractType", wrapper.GetContractTemplate)
	router.POST("/auth/internal/experimental/sign", wrapper.CreateSignSession)
	router.GET("/auth/internal/experimental/sign/:sessionPtr", wrapper.GetSignSessionStatus)

}

