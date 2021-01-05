package urlbuilder

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/xo/dburl"
	"makeit.imfr.cgi.com/lino/internal/app/localstorage"
	"makeit.imfr.cgi.com/lino/pkg/dataconnector"
)

func BuildURL(dc *dataconnector.DataConnector, out io.Writer) *dburl.URL {
	u, e2 := dburl.Parse(dc.URL)
	if e2 != nil {
		fmt.Fprintln(out, e2.Error())
		os.Exit(3)
	}
	// get user from env
	if dc.User.ValueFromEnv != "" {
		userFromEnv := os.Getenv(dc.User.ValueFromEnv)
		if userFromEnv == "" {
			if out != nil {
				fmt.Fprintf(out, "warn: missing environment variable %s", dc.User.ValueFromEnv)
				fmt.Fprintln(out)
			}
		} else {
			u.User = url.User(userFromEnv)
		}
	} else if dc.User.Value != "" {
		// set user from dc
		u.User = url.User(dc.User.Value)
	}
	// get password from env
	if dc.Password.ValueFromEnv != "" {
		passwordFromEnv := os.Getenv(dc.Password.ValueFromEnv)
		if passwordFromEnv == "" {
			if out != nil {
				fmt.Fprintf(out, "warn: missing environment variable %s", dc.Password.ValueFromEnv)
				fmt.Fprintln(out)
			}
		} else {
			u.User = url.UserPassword(u.User.Username(), passwordFromEnv)
		}
	}
	// if credentials still missing, check default store
	username := u.User.Username()
	_, passwordIsSet := u.User.Password()
	if username == "" || !passwordIsSet {
		store := defaultCredentialsStore()
		creds, err := client.Get(store, u.String())
		if err != nil {
			// failed to use credential store backend, fallback to local storage
			creds, err = localstorage.Read(u.String())
		}
		if err == nil {
			u.User = url.UserPassword(creds.Username, creds.Secret)
		}
	}
	return u
}

func StorePassword(u *dburl.URL, password string, out io.Writer) error {
	store := defaultCredentialsStore()
	creds := &credentials.Credentials{ServerURL: u.URL.String(), Username: u.URL.User.Username(), Secret: password}
	err := client.Store(store, creds)
	if err != nil {
		// failed to use credential store backend
		if credentials.IsCredentialsMissingServerURL(err) || credentials.IsCredentialsMissingUsername(err) || credentials.IsErrCredentialsNotFound(err) {
			return err
		}
		// fall back to local storage
		fmt.Fprintf(out, "warn: password will be stored unencrypted in %s, configure a credential helper to remove this warning. See https://github.com/docker/docker-credential-helpers", localstorage.GetFileLocation())
		return localstorage.Store(creds)
	}
	return nil
}
