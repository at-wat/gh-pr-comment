from setuptools import setup, find_packages

setup(
    name='gh_pr_comment',
    version='0.1.0.dev0',
    description='GitHub PR comment post',
    url='https://github.com/at-wat/gh-pr-comment',
    author='Atsushi Watanabe',
    author_email='atsushi.w@ieee.org',
    packages=find_packages(exclude=['tests']),
    install_requires=['requests'],
    entry_points={
        'console_scripts': [
            'gh-pr-comment = gh_pr_comment.post:post_main',
            'gh-pr-upload= gh_pr_comment.upload:post_main'
        ]
    },
    license="BSD"
)
