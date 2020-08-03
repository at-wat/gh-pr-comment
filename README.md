# gh-pr-comment

Command line tool for GitHub pull-request (issue) comment.

## Usage

```
$ gh-pr-comment title comment
```

## Install

### Download compiled binary
```
$ curl -sL https://raw.githubusercontent.com/at-wat/gh-pr-comment/master/install.sh | sh -s
```

## Required environment variables

### Travis-CI

- ***TRAVIS***: true
- ***TRAVIS\_PULL\_REQUEST\_SLUG***: owner/repos
- ***TRAVIS\_PULL\_REQUEST***: pull request number
- ***TRAVIS\_BOT\_GITHUB\_TOKEN*** or ***GITHUB\_TOKEN***: token with comment write permission
