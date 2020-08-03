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

To run gh-pr-comment inside docker container, required environemnt variables can be automatically specified by the following script.
```shell
$ docker run [your options...] \
    $(bash <(curl -s https://raw.githubusercontent.com/at-wat/gh-pr-comment/master/env.sh)) \
    image commands...
```

### General
- ***GITHUB\_TOKEN*** (or ***TRAVIS\_BOT\_GITHUB\_TOKEN***): token with comment write permission

### Optional
- ***GITHUB\_API\_URL\_BASE***: specify GitHub Enterprise or any custom endpoint URL

### Travis-CI
- ***TRAVIS***: true
- ***TRAVIS\_PULL\_REQUEST\_SLUG***: owner/repos
- ***TRAVIS\_PULL\_REQUEST***: pull request number
