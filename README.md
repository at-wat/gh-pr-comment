# gh-pr-comment

Command line tool for GitHub pull-request (issue) comment.

## Usage

```
$ gh-pr-comment title comment
```

## Install

```
$ pip install gh-pr-comment
```

## Required environment variables

- ***TRAVIS\_REPO\_SLUG***: user/repos
- ***TRAVIS\_PULL\_REQUEST***: pull request number
- ***TRAVIS\_BOT\_GITHUB\_TOKEN***: token with comment write permission

Two of first variables are automatically set by Travis-CI on pull-request builds.
