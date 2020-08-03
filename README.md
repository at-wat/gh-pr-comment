# gh-pr-comment

Command line tool for GitHub pull-request (issue) comment.

## Usage

```
$ gh-pr-comment title comment
```

## Install

### Download compiled binary
```shell
# Install latest version under ~/.local/bin
$ curl -sL https://raw.githubusercontent.com/at-wat/gh-pr-comment/master/install.sh | sh -s

# Install latest version under /path/to/bin
$ curl -sL https://raw.githubusercontent.com/at-wat/gh-pr-comment/master/install.sh \
  | sh -s latest /path/to/bin

# Install specific version under /usr/local/bin
$ curl -sL https://raw.githubusercontent.com/at-wat/gh-pr-comment/master/install.sh \
  | sh -s v0.5.0 /path/to/bin
```

## Required environment variables

### Travis-CI

- ***TRAVIS***: true
- ***TRAVIS\_PULL\_REQUEST\_SLUG***: owner/repos
- ***TRAVIS\_PULL\_REQUEST***: pull request number
- ***TRAVIS\_BOT\_GITHUB\_TOKEN*** or ***GITHUB\_TOKEN***: token with comment write permission
