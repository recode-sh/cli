package presenters

import (
	"errors"
	"fmt"
	"strings"

	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/exceptions"
	"github.com/recode-sh/cli/internal/system"
	"github.com/recode-sh/recode/entities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ViewableError struct {
	Title   string
	Message string
}

type ViewableErrorBuilder interface {
	Build(error) *ViewableError
}

type RecodeViewableErrorBuilder struct{}

func NewRecodeViewableErrorBuilder() RecodeViewableErrorBuilder {
	return RecodeViewableErrorBuilder{}
}

func (RecodeViewableErrorBuilder) Build(err error) (viewableError *ViewableError) {
	viewableError = &ViewableError{}

	if typedError, ok := err.(entities.ErrClusterNotExists); ok {
		viewableError.Title = "Cluster not found"
		viewableError.Message = fmt.Sprintf(
			"The cluster \"%s\" was not found.",
			typedError.ClusterName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrClusterAlreadyExists); ok {
		viewableError.Title = "Cluster already exists"
		viewableError.Message = fmt.Sprintf(
			"The cluster \"%s\" already exists.",
			typedError.ClusterName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrDevEnvNotExists); ok {
		viewableError.Title = "Development environment not found"

		if typedError.ClusterName != entities.DefaultClusterName {
			viewableError.Message = fmt.Sprintf(
				"The development environment \"%s\" was not found in the cluster \"%s\".",
				typedError.DevEnvName,
				typedError.ClusterName,
			)
			return
		}

		viewableError.Message = fmt.Sprintf(
			"The development environment \"%s\" was not found.",
			typedError.DevEnvName,
		)
		return
	}

	if errors.Is(err, exceptions.ErrUserNotLoggedIn) {
		viewableError.Title = "GitHub account not connected"
		viewableError.Message = fmt.Sprintf(
			"You must first connect your GitHub account using the command \"recode login\".\n\n"+
				"Recode requires the following permissions:\n\n"+
				"  - \"Public SSH keys\" and \"Repositories\" to let you access your repositories from your development environments\n\n"+
				"  - \"GPG Keys\" and \"Personal user data\" to configure Git and sign your commits (verified badge)\n\n"+
				"All your data (including the OAuth access token) will only be stored locally (in \"%s\").",
			system.UserConfigFilePath(),
		)

		return
	}

	if typedError, ok := err.(entities.ErrInvalidDevEnvUserConfig); ok {
		viewableError.Title = "Invalid user configuration"

		bold := constants.Bold
		viewableError.Message = fmt.Sprintf(
			"The repository \"%s/%s\" contains an invalid configuration.\n\n"+bold("Reason:")+" %s",
			typedError.RepoOwner,
			entities.DevEnvUserConfigRepoName,
			typedError.Reason,
		)

		return
	}

	if typedError, ok := err.(entities.ErrDevEnvRepositoryNotFound); ok {
		viewableError.Title = "Repository not found"
		viewableError.Message = fmt.Sprintf(
			"The repository \"%s/%s\" was not found.\n\n"+
				"Please double check that this repository exists and that you can access it.",
			typedError.RepoOwner,
			typedError.RepoName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrStartRemovingDevEnv); ok {
		viewableError.Title = "Invalid development environment state"
		viewableError.Message = fmt.Sprintf(
			"The development environment \"%s\" cannot be started because it is currently removing.\n\n"+
				"You must wait for the removing process to terminate.",
			typedError.DevEnvName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrStartStoppingDevEnv); ok {
		viewableError.Title = "Invalid development environment state"
		viewableError.Message = fmt.Sprintf(
			"The development environment \"%s\" cannot be started because it is currently stopping.\n\n"+
				"You must wait for the stopping process to terminate.",
			typedError.DevEnvName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrStopRemovingDevEnv); ok {
		viewableError.Title = "Invalid development environment state"
		viewableError.Message = fmt.Sprintf(
			"The development environment \"%s\" cannot be stopped because it is currently removing.\n\n"+
				"You must wait for the removing process to terminate.",
			typedError.DevEnvName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrStopCreatingDevEnv); ok {
		viewableError.Title = "Invalid development environment state"
		viewableError.Message = fmt.Sprintf(
			"The development environment \"%s\" cannot be stopped because it is currently creating.\n\n"+
				"You must wait for the creation process to terminate.",
			typedError.DevEnvName,
		)

		return
	}

	if typedError, ok := err.(entities.ErrStopStartingDevEnv); ok {
		viewableError.Title = "Invalid development environment state"
		viewableError.Message = fmt.Sprintf(
			"The development environment \"%s\" cannot be stopped because it is currently starting.\n\n"+
				"You must wait for the starting process to terminate.",
			typedError.DevEnvName,
		)

		return
	}

	if typedError, ok := err.(exceptions.ErrLoginError); ok {
		viewableError.Title = "GitHub connection error"
		viewableError.Message = fmt.Sprintf(
			"An error has occured during the authorization of the Recode application.\n\n%s",
			typedError.Reason,
		)

		return
	}

	if typedError, ok := err.(exceptions.ErrMissingRequirements); ok {
		viewableError.Title = "Missing requirements"
		viewableError.Message = fmt.Sprintf(
			"The following requirements are missing:\n\n  - %s",
			strings.Join(typedError.MissingRequirements, "\n\n  - "),
		)

		return
	}

	bold := constants.Bold

	if status, ok := status.FromError(err); ok {
		viewableError.Title = "Recode agent error"

		errorMessage := status.Message()

		if len(errorMessage) >= 2 {
			errorMessage = strings.ToTitle(errorMessage[0:1]) + errorMessage[1:] + "."
		}

		viewableError.Message = errorMessage

		if status.Code() != codes.Unknown {
			viewableError.Message += "\n\n" +
				bold("Error code: ") +
				status.Code().String()
		}

		return
	}

	viewableError.Title = "Unknown error"
	viewableError.Message = fmt.Sprintf(
		"An unknown error occurred.\n\n"+
			"You could try to fix it (using the details below) or open a new issue at: https://github.com/recode-sh/cli/issues/new\n\n"+
			bold("%s"),
		err.Error(),
	)

	return
}
