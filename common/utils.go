package common

import (
	"os/exec"
	"strings"
)

// Get the name of the author, return stupid placeholder if it cannot be obtained.
func GetAuthorName() string {
	// Get the git username for attribution purposes
	command := exec.Command("git", "config", "user.name")
	authNameTmp, err := command.Output()
	// shit happens when you party naked ~Socrates
	authName := "Thurman Merman"
	if err == nil {
		authName = string(authNameTmp)
	}
	authName = strings.ReplaceAll(authName, "\n", "")
	return authName
}
