# gh-pr-comment

Command line tool for GitHub pull-request (issue) comment.

## Usage

```
$ gh-pr-comment title comment
```

## Install

### GitHub Actions

```yaml
      - uses: at-wat/setup-gh-pr-comment@v0
      ...
      - run: gh-pr-comment "title" "message"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

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
`env.sh` is also bundled as `gh-pr-comment-env.sh` in release archives.

### General
- ***GITHUB\_TOKEN***: token with comment write permission

### Optional
- ***GITHUB\_API\_URL\_BASE***: specify GitHub Enterprise or any custom endpoint URL

### Travis-CI
- ***TRAVIS***: true
- ***TRAVIS\_PULL\_REQUEST\_SLUG***: owner/repos
- ***TRAVIS\_PULL\_REQUEST***: pull request number

## Security

If your environment handles sensitive information, it is recommended to download `install.sh` and `env.sh` from https://github.com/at-wat/gh-pr-comment/releases and verify the signature using GPG.
My public key is available at https://github.com/at-wat.gpg and SKS keyserver pool.

(fingerprint: `358B DF63 B4AE D76A 871A  E62E 1BF1 686B 468C 35B2`)
