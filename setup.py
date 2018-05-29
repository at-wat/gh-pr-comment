from setuptools import setup, find_packages

setup(
    name='gh_pr_comment',
    version='0.4.0',
    description='GitHub PR comment post',
    url='https://github.com/at-wat/gh-pr-comment',
    author='Atsushi Watanabe',
    author_email='atsushi.w@ieee.org',
    packages=find_packages(exclude=['tests']),
    install_requires=['requests', 'uuid'],
    entry_points={
        'console_scripts': [
            'gh-pr-comment = gh_pr_comment.comment:post_main',
            'gh-pr-upload= gh_pr_comment.upload:post_main'
        ]
    },
    license="BSD"
)
