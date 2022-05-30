package system

import (
	"bufio"
	"io"
	"strings"

	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/interfaces"
)

func AskForConfirmation(
	logger interfaces.Logger,
	stdin io.Reader,
	question string,
) (bool, error) {

	stdinReader := bufio.NewReader(stdin)

	logger.Log(constants.Bold(constants.Yellow("Warning!") + " " + question))

	logger.Log("\nOnly \"yes\" will be accepted to confirm. (You could use \"--force\" next time).\n")
	logger.LogNoNewline(constants.Bold("Confirm? "))

	response, err := stdinReader.ReadString('\n')

	if err != nil {
		return false, err
	}

	sanitizedResponse := strings.TrimSpace(response)

	if sanitizedResponse == "yes" {
		return true, nil
	}

	logger.Log("")

	return false, nil
}
