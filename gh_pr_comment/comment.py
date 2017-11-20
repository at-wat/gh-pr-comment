import sys
import os
import json
import requests


def post_main():
    argv = sys.argv
    if len(argv) < 3:
        sys.stderr.write(
            'usage: gh-pr-comment title comment\n')
        sys.stderr.write(
            'env:\n')
        sys.stderr.write(
            '- TRAVIS_REPO_SLUG       :'
            + ' user/repos\n')
        sys.stderr.write(
            '- TRAVIS_PULL_REQUEST    :'
            + ' pull request number\n')
        sys.stderr.write(
            '- TRAVIS_BOT_GITHUB_TOKEN:'
            + ' token with comment write permission\n')
        sys.exit(1)

    post(argv[1], argv[2])


def post(title, contents):
    if 'TRAVIS_BOT_GITHUB_TOKEN' not in os.environ:
        sys.stderr.write(
            'gh-pr-comment: TRAVIS_BOT_GITHUB_TOKEN is not set.\n')
        sys.exit(1)
    if 'TRAVIS_REPO_SLUG' not in os.environ:
        sys.stderr.write(
            'gh-pr-comment: TRAVIS_REPO_SLUG is not set.\n')
        sys.exit(1)
    if 'TRAVIS_PULL_REQUEST' not in os.environ:
        sys.stderr.write(
            'gh-pr-comment: TRAVIS_PULL_REQUEST is not set.\n')
        sys.exit(1)

    if os.environ['TRAVIS_PULL_REQUEST'] == 'false':
        sys.stderr.write(
            'gh-pr-comment: TRAVIS_PULL_REQUEST is false.\n')
        sys.exit(0)

    url = 'https://api.github.com/repos/' \
        + os.environ['TRAVIS_REPO_SLUG'] \
        + '/issues/' \
        + os.environ['TRAVIS_PULL_REQUEST'] \
        + '/comments?access_token=' \
        + os.environ['TRAVIS_BOT_GITHUB_TOKEN']
    headers = {"Content-Type": "application/json"}

    body = {
        "body": '## %s' % title + '\n\n' + contents
    }

    r = requests.post(url, data=json.dumps(body), headers=headers)
    print(r.text)
