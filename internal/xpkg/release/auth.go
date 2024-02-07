package release

import (
	"io"
	"os"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/pkg/errors"

	configv1 "github.com/mistermx/xpreleaser/config/v1"
)

var (
	amazonKeychain authn.Keychain = authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithLogger(io.Discard)))
	azureKeychain  authn.Keychain = authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper())
)

// fixedKeychain is a key chain that uses static credentials (i.e. basic auth)
type fixedKeychain map[authn.Resource]authn.Authenticator

func (fk fixedKeychain) Resolve(target authn.Resource) (authn.Authenticator, error) {
	if auth, ok := fk[target]; ok {
		return auth, nil
	}
	return authn.Anonymous, nil
}

// BuildKeyChainFromConfig creates a new [authn.Keychain] from a list of docker
// login configurations. It always includes [authn.DefaultKeychain].
func BuildKeyChainFromConfig(logins []configv1.DockerConfigLogin) (authn.Keychain, error) {
	fixed := fixedKeychain{}
	keychains := []authn.Keychain{
		authn.DefaultKeychain,
		fixed,
	}
	for _, login := range logins {
		switch login.Type {
		case configv1.DockerConfigLoginTypeAWS:
			keychains = append(keychains, amazonKeychain)
		case configv1.DockerConfigLoginTypeAzure:
			keychains = append(keychains, azureKeychain)
		case configv1.DockerConfigLoginTypeGoogle:
			keychains = append(keychains, google.Keychain)
		case configv1.DockerConfigLoginTypeBasic:
			if login.Registry == "" {
				return nil, errors.New("basic auth requires registry")
			}
			registry, err := name.NewRegistry(login.Registry)
			if err != nil {
				return nil, errors.Wrap(err, "invalid registry")
			}
			if login.Basic.UsernameFromEnv == "" || login.Basic.PasswordFromEnv == "" {
				return nil, errors.New("basic auth requires username and password")
			}
			fixed[registry] = &authn.Basic{
				Username: os.Getenv(login.Basic.UsernameFromEnv),
				Password: os.Getenv(login.Basic.PasswordFromEnv),
			}
		default:
			return nil, errors.Errorf("unkown login type %q", login.Type)
		}
	}
	return authn.NewMultiKeychain(keychains...), nil
}
