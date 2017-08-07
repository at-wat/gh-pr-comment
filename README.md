# gh-pr-comment

## Usage

```
$ gh-pr-comment title comment
```

The script escapes special charactors for json encoding.

## Required environment variables

- ***TRAVIS\_REPO\_SLUG*** = user/repos
- ***TRAVIS\_PULL\_REQUEST*** = pull request number
- ***TRAVIS\_BOT\_GITHUB\_TOKEN*** = token with comment write permission

## Usage in your package

```
$ wget -q https://raw.githubusercontent.com/at-wat/gh-pr-comment/master/gh-pr-comment.sh
$ source gh-pr-comment.sh
```

