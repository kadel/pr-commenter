# PR Commenter

PR Commenter is a tool to automate the process of commenting on GitHub pull requests. It allows you to specify a PR and post or edit comments programmatically. This tool is particularly useful for CI/CD pipelines, automated feedback systems, or any application where automated comments on pull requests are required.

## Installation

```bash
go get github.com/kadel/pr-commenter
```

## Usage

After installation, you can run the tool via the command line.

```bash
pr-commenter [flags]
```

Certainly! Improving the descriptions can provide more clarity to users. Here's a more detailed version:

## Flags

- **-key-from-file**: Specifies the path where the private key file is located. This key is used for authentication and making authorized requests.

- **-key-from-env-var**:Instead of reading the private key from a file, you can store it as a **base64-encoded** string in an environment variable. This flag denotes the name of that environment variable.

- **-pr-comment**: Identifies the specific Pull Request (by its number) you want to post or edit a comment on.

- **-prefix**: When you want to edit a comment, this prefix helps identify which one. If there are multiple comments, the tool will look for a comment starting with this prefix to edit. If no prefix is set, it will default to editing the first comment.

- **-application-id**: Each GitHub application has a unique ID. This flag is used to specify the ID of the GitHub application in context.

- **-repository**: The name of the GitHub repository where the Pull Request resides.

- **-org**:If the repository is part of a GitHub organization (as opposed to a user), specify the organization's name with this flag.

## Examples

1. **Commenting on a PR using an environment variable**:

```bash
echo 'Test failed, see <a href="https://example.com/logs">logs</a> for more information.' | pr-commenter -key-from-file=/path/to/key/file -application-id=123 -pr-comment=123 -repository=my-repo -org=my-org
```

2. **Editing an existing comment with a specific prefix**:

```bash
echo 'status: failed, see <a href="https://example.com/logs">logs</a> for more information' | pr-commenter -key-from-env-var=MY_ENV_VARIABLE -application-id=123 -pr-comment=123 -repository=my-repo -org=my-org -prefix="MY TEST JOB:"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
