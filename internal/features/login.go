package features

import (
	"errors"
	"fmt"
	"html"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/recode-sh/cli/internal/config"
	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/exceptions"
	"github.com/recode-sh/cli/internal/interfaces"
	"golang.org/x/oauth2"
)

type LoginInput struct{}

type LoginResponseContent struct{}

type LoginResponse struct {
	Error   error
	Content LoginResponseContent
}

type LoginPresenter interface {
	PresentToView(LoginResponse)
}

type LoginFeature struct {
	presenter  LoginPresenter
	logger     interfaces.Logger
	browser    interfaces.BrowserManager
	userConfig interfaces.UserConfigManager
	sleeper    interfaces.Sleeper
	github     interfaces.GitHubManager
}

func NewLoginFeature(
	presenter LoginPresenter,
	logger interfaces.Logger,
	browser interfaces.BrowserManager,
	config interfaces.UserConfigManager,
	sleeper interfaces.Sleeper,
	github interfaces.GitHubManager,
) LoginFeature {

	return LoginFeature{
		presenter:  presenter,
		logger:     logger,
		browser:    browser,
		userConfig: config,
		sleeper:    sleeper,
		github:     github,
	}
}

func (l LoginFeature) Execute(input LoginInput) error {
	handleError := func(err error) error {
		l.presenter.PresentToView(LoginResponse{
			Error: exceptions.ErrLoginError{
				Reason: err.Error(),
			},
		})

		return err
	}

	gitHubOAuthCbHandlerResp := struct {
		Error       error
		AccessToken string
	}{}

	gitHubOauthCbHandlerDoneChan := make(chan struct{})

	gitHubOAuthCbHandler := func(w http.ResponseWriter, r *http.Request) {
		defer close(gitHubOauthCbHandlerDoneChan)

		queryComponents, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			gitHubOAuthCbHandlerResp.Error = err
			return
		}

		errorInQuery, hasErrorInQuery := queryComponents["error"]

		if hasErrorInQuery {
			msg := "<h1>Error!</h1>"
			msg = msg + "<p>" + html.EscapeString(errorInQuery[0]) + "</p>"

			w.WriteHeader(500)
			w.Write([]byte(msg))

			gitHubOAuthCbHandlerResp.Error = errors.New(errorInQuery[0])
			return
		}

		accessTokenInQuery, hasAccessTokenInQuery := queryComponents["access_token"]

		if !hasAccessTokenInQuery {
			msg := "<h1>Error!</h1>"
			msg = msg + "<p>An unknown error occured during GitHub connection. Please retry.</p>"

			w.WriteHeader(500)
			w.Write([]byte(msg))

			gitHubOAuthCbHandlerResp.Error = errors.New("no access token returned after authorization")
			return
		}

		msg := "<h1>Success!</h1>"
		msg = msg + "<p>Your GitHub account is now connected. You can close this tab and go back to the Recode CLI.</p>"

		w.WriteHeader(200)
		w.Write([]byte(msg))

		gitHubOAuthCbHandlerResp.AccessToken = accessTokenInQuery[0]
	} // <- End of gitHubOAuthCbHandler

	http.HandleFunc(
		config.GitHubOAuthAPIToCLIURLPath,
		gitHubOAuthCbHandler,
	)

	// Assign a random port to our http server
	httpListener, err := net.Listen("tcp", ":0")

	if err != nil {
		return handleError(err)
	}

	httpServerServeErrorChan := make(chan error, 1)
	go func() {
		httpServerServeErrorChan <- http.Serve(httpListener, nil)
	}()

	httpListenPort := httpListener.Addr().(*net.TCPAddr).Port

	gitHubOAuthClient := &oauth2.Config{
		ClientID: config.GitHubOAuthClientID,
		Scopes:   config.GitHubOAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://github.com/login/oauth/authorize",
		},
		RedirectURL: config.GitHubOAuthCLIToAPIURL,
	}

	// Listen port is passed through OAuth
	// state because GitHub doesn't support
	// dynamic redirect URIs
	gitHubOAuthAuthorizeURL := gitHubOAuthClient.AuthCodeURL(
		fmt.Sprintf("%d", httpListenPort),
	)

	bold := constants.Bold
	l.logger.Log(bold("\nYou will be taken to your browser to connect your GitHub account...\n"))

	l.logger.Info("If your browser doesn't open automatically, go to the following link:\n")
	l.logger.Log("%s", gitHubOAuthAuthorizeURL)

	l.sleeper.Sleep(4 * time.Second)

	if err := l.browser.OpenURL(gitHubOAuthAuthorizeURL); err != nil {
		l.logger.Error(
			"\nCannot open browser! Please visit above URL ↑",
		)
	}

	l.logger.Warning("\nWaiting for GitHub authorization... (Press Ctrl-C to quit)\n")

	select {
	case httpServerServeError := <-httpServerServeErrorChan:
		return handleError(httpServerServeError)
	case <-gitHubOauthCbHandlerDoneChan:
		// We swallow the httpListener.Close() error here
		// given that the CLI will exit and force all
		// resources to be released
		_ = httpListener.Close()
	}

	if gitHubOAuthCbHandlerResp.Error != nil {
		return handleError(gitHubOAuthCbHandlerResp.Error)
	}

	githubUser, err := l.github.GetAuthenticatedUser(
		gitHubOAuthCbHandlerResp.AccessToken,
	)

	if err != nil {
		return handleError(err)
	}

	l.userConfig.Set(
		config.UserConfigKeyUserIsLoggedIn,
		true,
	)

	l.userConfig.Set(
		config.UserConfigKeyGitHubAccessToken,
		gitHubOAuthCbHandlerResp.AccessToken,
	)

	l.userConfig.PopulateFromGitHubUser(
		githubUser,
	)

	if err := l.userConfig.WriteConfig(); err != nil {
		return handleError(err)
	}

	l.presenter.PresentToView(LoginResponse{})
	return nil
}
