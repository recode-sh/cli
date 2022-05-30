package features

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/recode-sh/agent/constants"
	"github.com/recode-sh/agent/proto"
	"github.com/recode-sh/cli/internal/agent"
	"github.com/recode-sh/cli/internal/config"
	cliConstants "github.com/recode-sh/cli/internal/constants"
	cliEntities "github.com/recode-sh/cli/internal/entities"
	"github.com/recode-sh/cli/internal/interfaces"
	"github.com/recode-sh/recode/actions"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/features"
)

type StartResponse struct {
	Error   error
	Content StartResponseContent
}

type StartResponseContent struct {
	DevEnvName           string
	DevEnvAlreadyStarted bool
	DevEnvRebuilt        bool
}

type StartPresenter interface {
	PresentToView(StartResponse)
}

type StartOutputHandler struct {
	userConfig         interfaces.UserConfigManager
	presenter          StartPresenter
	agentClientBuilder agent.ClientBuilder
	github             interfaces.GitHubManager
	logger             interfaces.Logger
	sshConfig          interfaces.SSHConfigManager
	sshKeys            interfaces.SSHKeysManager
	sshKnownHosts      interfaces.SSHKnownHostsManager
	vscodeProcess      interfaces.VSCodeProcessManager
	vscodeExtensions   interfaces.VSCodeExtensionsManager
}

func NewStartOutputHandler(
	userConfig interfaces.UserConfigManager,
	presenter StartPresenter,
	agentClientBuilder agent.ClientBuilder,
	github interfaces.GitHubManager,
	logger interfaces.Logger,
	sshConfig interfaces.SSHConfigManager,
	sshKeys interfaces.SSHKeysManager,
	sshKnownHosts interfaces.SSHKnownHostsManager,
	vscodeProcess interfaces.VSCodeProcessManager,
	vscodeExtensions interfaces.VSCodeExtensionsManager,
) StartOutputHandler {

	return StartOutputHandler{
		userConfig:         userConfig,
		presenter:          presenter,
		agentClientBuilder: agentClientBuilder,
		github:             github,
		logger:             logger,
		sshConfig:          sshConfig,
		sshKnownHosts:      sshKnownHosts,
		sshKeys:            sshKeys,
		vscodeProcess:      vscodeProcess,
		vscodeExtensions:   vscodeExtensions,
	}
}

func (s StartOutputHandler) HandleOutput(output features.StartOutput) error {

	stepper := output.Stepper

	handleError := func(err error) error {
		stepper.StopCurrentStep()

		s.presenter.PresentToView(StartResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	devEnv := output.Content.DevEnv

	devEnvCreated := output.Content.DevEnvCreated
	devEnvStarted := output.Content.DevEnvStarted
	devEnvRebuildAsked := output.Content.DevEnvRebuildAsked

	devEnvAlreadyStarted := !devEnvCreated && !devEnvStarted && !devEnvRebuildAsked

	var devEnvAdditionalProperties *cliEntities.DevEnvAdditionalProperties

	if len(devEnv.AdditionalPropertiesJSON) > 0 {
		err := json.Unmarshal(
			[]byte(devEnv.AdditionalPropertiesJSON),
			&devEnvAdditionalProperties,
		)

		if err != nil {
			return handleError(err)
		}
	}

	if devEnvAdditionalProperties == nil {
		devEnvAdditionalProperties = &cliEntities.DevEnvAdditionalProperties{}
	}

	if devEnvCreated || devEnvRebuildAsked {
		if !devEnvRebuildAsked {
			stepper.StartTemporaryStep(
				"Building your development environment",
			)
		}

		agentClient := s.agentClientBuilder.Build(
			agent.NewDefaultClientConfig(
				[]byte(devEnv.SSHKeyPairPEMContent),
				devEnv.InstancePublicIPAddress,
			),
		)

		err := agentClient.InitInstance(&proto.InitInstanceRequest{
			DevEnvNameSlug:  devEnv.GetNameSlug(),
			GithubUserEmail: s.userConfig.GetString(config.UserConfigKeyGitHubEmail),
			UserFullName:    s.userConfig.GetString(config.UserConfigKeyGitHubFullName),
		}, func(stream agent.InitInstanceStream) error {

			for {
				initInstanceReply, err := stream.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					return err
				}

				if initInstanceReply.GithubSshPublicKeyContent != nil &&
					devEnvAdditionalProperties.GitHubCreatedSSHKeyId == nil {

					sshKeyInGitHub, err := s.github.CreateSSHKey(
						s.userConfig.GetString(config.UserConfigKeyGitHubAccessToken),
						fmt.Sprintf("recode-%s", devEnv.GetNameSlug()),
						initInstanceReply.GetGithubSshPublicKeyContent(),
					)

					if err != nil {
						return err
					}

					devEnvAdditionalProperties.GitHubCreatedSSHKeyId = sshKeyInGitHub.ID
					err = devEnv.SetAdditionalPropertiesJSON(devEnvAdditionalProperties)

					if err != nil {
						return err
					}

					err = actions.UpdateDevEnvInConfig(
						stepper,
						output.Content.CloudService,
						output.Content.RecodeConfig,
						output.Content.Cluster,
						devEnv,
					)

					if err != nil {
						return err
					}
				}

				if initInstanceReply.GithubGpgPublicKeyContent != nil &&
					devEnvAdditionalProperties.GitHubCreatedGPGKeyId == nil {

					gpgKeyInGitHub, err := s.github.CreateGPGKey(
						s.userConfig.GetString(config.UserConfigKeyGitHubAccessToken),
						initInstanceReply.GetGithubGpgPublicKeyContent(),
					)

					if err != nil {
						return err
					}

					devEnvAdditionalProperties.GitHubCreatedGPGKeyId = gpgKeyInGitHub.ID
					err = devEnv.SetAdditionalPropertiesJSON(devEnvAdditionalProperties)

					if err != nil {
						return err
					}

					err = actions.UpdateDevEnvInConfig(
						stepper,
						output.Content.CloudService,
						output.Content.RecodeConfig,
						output.Content.Cluster,
						devEnv,
					)

					if err != nil {
						return err
					}
				}
			}

			return nil
		})

		if err != nil {
			return handleError(err)
		}

		resolvedDevEnvUserConfig := devEnv.ResolvedUserConfig
		resolvedRepository := devEnv.ResolvedRepository

		err = agentClient.BuildAndStartDevEnv(
			&proto.BuildAndStartDevEnvRequest{
				DevEnvRepoOwner:     resolvedRepository.Owner,
				DevEnvRepoName:      resolvedRepository.Name,
				UserConfigRepoOwner: resolvedDevEnvUserConfig.RepoOwner,
				UserConfigRepoName:  resolvedDevEnvUserConfig.RepoName,
			},
			func(stream agent.BuildAndStartDevEnvStream) error {

				var hasUncompletedLogLine = false
				var uncompletedLogLineBuf bytes.Buffer

				newLineIndentSpaces := "                      "
				formatLogLine := func(logLine string) string {
					return strings.TrimPrefix(strings.ReplaceAll(
						strings.ReplaceAll(
							logLine, "\n", "\n"+newLineIndentSpaces,
						),
						"\r",
						"\r"+newLineIndentSpaces,
					), "\n"+newLineIndentSpaces)
				}

				for {
					startDevEnvReply, err := stream.Recv()

					// Make sure to display all logs
					// especially in case of error
					if uncompletedLogLineBuf.Len() > 0 {
						s.logger.LogNoNewline(strings.TrimSuffix(uncompletedLogLineBuf.String(), "\n"))
						uncompletedLogLineBuf = bytes.Buffer{}
					}

					if err == io.EOF {
						break
					}

					if err != nil {
						return err
					}

					stepper.StopCurrentStep()

					if len(startDevEnvReply.LogLineHeader) > 0 {
						bold := cliConstants.Bold
						blue := cliConstants.Blue
						s.logger.Log(bold(blue("> "+startDevEnvReply.LogLineHeader)) + "\n")
					}

					if len(startDevEnvReply.LogLine) > 0 {
						goLoggerNoNewLine := log.New(&uncompletedLogLineBuf, "  ", log.Flags())
						goLoggerNewLine := log.New(s.logger, "  ", log.Flags())

						// No prefix for empty new lines
						if startDevEnvReply.LogLine == "\n" {
							hasUncompletedLogLine = false

							s.logger.Log("\n")

							continue
						}

						if strings.HasSuffix(startDevEnvReply.LogLine, "\n") &&
							!strings.Contains(startDevEnvReply.LogLine, "\r") {

							if hasUncompletedLogLine {
								hasUncompletedLogLine = false

								// Ends the log line
								// and add a new line at end
								s.logger.Log(strings.TrimSuffix(
									formatLogLine(
										startDevEnvReply.LogLine,
									),
									"\n"+newLineIndentSpaces,
								) + "\n")

								continue
							}

							hasUncompletedLogLine = false

							// Start a complete log line by prefix
							// and add a new line at end
							goLoggerNewLine.Print(strings.TrimSuffix(
								formatLogLine(
									startDevEnvReply.LogLine,
								),
								"\n"+newLineIndentSpaces,
							) + "\n")

							continue
						}

						if !hasUncompletedLogLine {
							// Start an uncompleted log line by prefix
							// but without a new line at end
							goLoggerNoNewLine.Print(formatLogLine(
								startDevEnvReply.LogLine,
							))

							hasUncompletedLogLine = true
							continue
						}

						// Continue logging uncompleted log line
						// without adding prefix or new line
						s.logger.LogNoNewline(formatLogLine(
							startDevEnvReply.LogLine,
						))
					}
				}

				return nil
			},
		)

		if err != nil {
			return handleError(err)
		}
	}

	if !devEnvAlreadyStarted {
		err := output.Content.SetDevEnvAsStarted()

		if err != nil {
			return handleError(err)
		}
	}

	stepper.StartTemporaryStepWithoutNewLine(
		"Updating your local SSH configuration",
	)

	sshPEMPath, err := s.sshKeys.CreateOrReplacePEM(
		devEnv.GetSSHKeyPairName(),
		devEnv.SSHKeyPairPEMContent,
	)

	if err != nil {
		return handleError(err)
	}

	sshServerListenPort, err := strconv.ParseInt(
		constants.SSHServerListenPort,
		10,
		64,
	)

	if err != nil {
		return handleError(err)
	}

	sshConfigHostKey := devEnv.Name

	err = s.sshConfig.AddOrReplaceHost(
		sshConfigHostKey,
		devEnv.InstancePublicIPAddress,
		sshPEMPath,
		entities.DevEnvRootUser,
		sshServerListenPort,
	)

	if err != nil {
		return handleError(err)
	}

	for _, sshHostKey := range devEnv.SSHHostKeys {
		err := s.sshKnownHosts.AddOrReplace(
			devEnv.InstancePublicIPAddress,
			sshHostKey.Algorithm,
			sshHostKey.Fingerprint,
		)

		if err != nil {
			return handleError(err)
		}
	}

	stepper.StartTemporaryStepWithoutNewLine(
		"Installing Visual Studio Code Remote - SSH extension",
	)

	_, err = s.vscodeExtensions.Install("ms-vscode-remote.remote-ssh")

	if err != nil {
		return handleError(err)
	}

	stepper.StartTemporaryStepWithoutNewLine(
		"Opening Visual Studio Code",
	)

	_, err = s.vscodeProcess.OpenOnRemote(
		sshConfigHostKey,
		constants.DevEnvVSCodeWorkspaceConfigFilePath,
	)

	if err != nil {
		return handleError(err)
	}

	stepper.StopCurrentStep()

	s.presenter.PresentToView(StartResponse{
		Content: StartResponseContent{
			DevEnvName:           devEnv.Name,
			DevEnvAlreadyStarted: devEnvAlreadyStarted,
			DevEnvRebuilt:        devEnvRebuildAsked,
		},
	})

	return nil
}
