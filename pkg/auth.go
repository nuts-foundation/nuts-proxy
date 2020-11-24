package pkg

import (
	"sync"
	"time"

	nutscrypto "github.com/nuts-foundation/nuts-crypto/pkg"
	core "github.com/nuts-foundation/nuts-go-core"
	registry "github.com/nuts-foundation/nuts-registry/pkg"

	contract2 "github.com/nuts-foundation/nuts-auth/pkg/contract"
	"github.com/nuts-foundation/nuts-auth/pkg/services"
	"github.com/nuts-foundation/nuts-auth/pkg/services/contract"
	"github.com/nuts-foundation/nuts-auth/pkg/services/dummy"
	"github.com/nuts-foundation/nuts-auth/pkg/services/oauth"
	"github.com/nuts-foundation/nuts-auth/pkg/services/validator"
)

// ConfAddress is the config key for the address the http server listens on
const ConfAddress = "address"

// PublicURL is the config key for the public URL the http/irma server can be discovered
const PublicURL = "publicUrl"

// ConfMode is the config name for the engine mode
const ConfMode = "mode"

const ConfEnableCORS = "enableCORS"

// ConfActingPartyCN is the config key to provide the Acting party common name
const ConfActingPartyCN = "actingPartyCn"

// ConfContractValidators is the config key for defining which contract validators to use
const ConfContractValidators = "contractValidators"

// AuthClient is the interface which should be implemented for clients or mocks
type AuthClient interface {
	// OAuthClient returns an instance of OAuthClient
	OAuthClient() services.OAuthClient
	// ContractClient returns an instance of ContractClient
	ContractClient() services.ContractClient
	// ContractNotary returns an instance of ContractNotary
	ContractNotary() services.ContractNotary
	// Signer returns a signer for a given signerID
	Signer(signerID string) contract2.Signer
}

// Auth is the main struct of the Auth service
type Auth struct {
	Config              AuthConfig
	configOnce          sync.Once
	configDone          bool
	OAuth               services.OAuthClient
	oneOauthInstance    sync.Once
	Contract            services.ContractClient
	oneContractInstance sync.Once
	Crypto              nutscrypto.Client
	Registry            registry.RegistryClient
	contractNotary      services.ContractNotary
	signerBackends      map[string]contract2.Signer
}

func (auth *Auth) Signer(signerID string) contract2.Signer {
	return auth.signerBackends[signerID]
}

func (auth *Auth) ContractNotary() services.ContractNotary {
	return auth.contractNotary
}

func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		Address:            "localhost:1323",
		IrmaSchemeManager:  "pbdf",
		ContractValidators: []string{"irma", "uzi"},
	}
}

var instance *Auth
var oneBackend sync.Once

// AuthInstance create an returns a singleton of the Auth struct
func AuthInstance() *Auth {
	if instance != nil {
		return instance
	}
	oneBackend.Do(func() {
		instance = NewAuthInstance(DefaultAuthConfig(), nutscrypto.CryptoInstance(), registry.RegistryInstance())
	})
	return instance
}

func NewAuthInstance(config AuthConfig, cryptoClient nutscrypto.Client, registryClient registry.RegistryClient) *Auth {
	return &Auth{
		Config:         config,
		Crypto:         cryptoClient,
		Registry:       registryClient,
		contractNotary: contract.NewContractNotary(registryClient, 60*time.Minute),
		// todo: put in separate init function
		signerBackends: map[string]contract2.Signer{"dummy": dummy.Dummy{InStrictMode: false, Sessions: make(map[string]string), Status: make(map[string]string)}},
	}
}

// OAuthClient returns an instance of OAuthClient
func (auth *Auth) OAuthClient() services.OAuthClient {
	if auth.OAuth != nil {
		return auth.OAuth
	}
	auth.oneOauthInstance.Do(func() {
		auth.OAuth = oauth.NewOAuthService(core.NutsConfig().VendorID(), auth.Crypto, auth.Registry, auth.Contract)
	})
	return auth.OAuth
}

// ContractClient returns an instance of ContractClient
func (auth *Auth) ContractClient() services.ContractClient {
	if auth.Contract != nil {
		return auth.Contract
	}
	auth.oneContractInstance.Do(func() {
		cfg := validator.Config{
			Mode:                      auth.Config.Mode,
			Address:                   auth.Config.Address,
			PublicUrl:                 auth.Config.PublicUrl,
			IrmaConfigPath:            auth.Config.IrmaConfigPath,
			IrmaSchemeManager:         auth.Config.IrmaSchemeManager,
			SkipAutoUpdateIrmaSchemas: auth.Config.SkipAutoUpdateIrmaSchemas,
			ActingPartyCn:             auth.Config.ActingPartyCn,
		}
		auth.Contract = validator.NewContractInstance(cfg, auth.Crypto, auth.Registry)
	})
	return auth.Contract
}

// Configure the Auth struct by creating a validator and create an Irma server
func (auth *Auth) Configure() (err error) {
	auth.configOnce.Do(func() {
		auth.Config.Mode = core.NutsConfig().GetEngineMode(auth.Config.Mode)
		if auth.Config.Mode == core.ServerEngineMode {

			auth.ContractClient()
			if err = auth.Contract.Configure(); err != nil {
				return
			}

			auth.OAuthClient()
			if err = auth.OAuth.Configure(); err != nil {
				return
			}
			auth.configDone = true
		}
	})

	return err
}
