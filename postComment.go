package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func getJwtToken(keyFileContent []byte, applicationId int64) string {

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyFileContent)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Add(-60 * time.Second).Unix(),
		"iss": applicationId,
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func main() {

	var keyPath string
	flag.StringVar(&keyPath, "key-from-file", "", "Path to the private key file")
	var keyEnvVar string
	flag.StringVar(&keyEnvVar, "key-from-env-var", "", "Name of the environment variable containing base64-encoded private key")
	var prNumber int
	flag.IntVar(&prNumber, "pr-comment", 0, "PR number to post comment to")
	var prefix string
	flag.StringVar(&prefix, "prefix", "", "Prefix to detect existing comment. If not set, the first existing comment will be edited. If set the first comment starting with the prefix will be edited.")
	var applicationId int64
	flag.Int64Var(&applicationId, "application-id", 0, "ID of the GitHub application")
	var repository string
	flag.StringVar(&repository, "repository", "", "Name of the repository")
	var org string
	flag.StringVar(&org, "org", "", "Name of the organization")
	flag.Parse()

	if applicationId == 0 {
		panic("no application id provided, please use -application-id flag")
	}
	if repository == "" {
		panic("no repository provided, please use -repository flag")
	}
	if org == "" {
		panic("no organization provided, please use -org flag")
	}

	// if --pr-number was set read text rom stdin
	var text string
	fmt.Println(prNumber)
	if prNumber != 0 {
		fmt.Println("read text from stdin")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			text += line + "\n"
		}
	}

	ctx := context.Background()
	var keyFileContent []byte
	if keyPath != "" {
		var err error
		keyFileContent, err = os.ReadFile(keyPath)
		if err != nil {
			panic(err)
		}
	} else if keyEnvVar != "" {
		var err error
		keyFileContent, err = base64.StdEncoding.DecodeString(os.Getenv(keyEnvVar))
		if err != nil {
			panic(err)
		}
	} else {
		panic("no private key provided, please use -key-from-file or -key-from-env-var flags")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getJwtToken(keyFileContent, applicationId)},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// find installation id
	// https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-an-installation

	installations, _, err := client.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil {
		panic(err)
	}

	var installId int64
	var appSlug string
	for _, installation := range installations {
		if *installation.Account.Login == org {
			installId = *installation.ID
			appSlug = *installation.AppSlug
		}
	}

	token, _, err := client.Apps.CreateInstallationToken(ctx, installId, &github.InstallationTokenOptions{})
	if err != nil {
		panic(err)
	}

	// create a client for the application installation
	installationClient := github.NewClient(
		oauth2.NewClient(ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: *token.Token},
			),
		),
	)

	if prNumber != 0 {
		comments, _, err := installationClient.Issues.ListComments(ctx, org, repository, prNumber, &github.IssueListCommentsOptions{})
		if err != nil {
			panic(err)
		}

		var existingComment int64 = 0

		for _, comment := range comments {

			if *comment.User.Login == fmt.Sprintf("%s[bot]", appSlug) {
				if strings.HasPrefix(comment.GetBody(), prefix) {
					existingComment = *comment.ID
					fmt.Printf("detected existing comment id:%d\n", existingComment)
					break
				}
			}
		}

		var comment *github.IssueComment
		if existingComment != 0 {
			comment, _, err = installationClient.Issues.EditComment(ctx, org, repository, existingComment,
				&github.IssueComment{
					Body: &text,
				})
		} else {
			comment, _, err = installationClient.Issues.CreateComment(ctx, org, repository, prNumber,
				&github.IssueComment{
					Body: &text,
				})
		}

		if err != nil {
			panic(err)
		}
		b, _ := json.MarshalIndent(comment, "", "  ")
		fmt.Println(string(b))
	}
}
