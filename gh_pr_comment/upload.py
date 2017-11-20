import sys
import os
import requests


def post_main():
    argv = sys.argv
    if len(argv) < 2:
        sys.stderr.write('usage: gh-pr-upload filename\n')
        sys.stderr.write('env:\n')
        sys.stderr.write('- IMAGE_UPLOADER : imgur(default)\n')
        sys.stderr.write('- IMGUR_CLIENT_ID: custom-client-id\n')
        sys.stderr.write('return: image url\n')
        sys.exit(1)

    post(argv[1])


def post(filename):
    if 'IMAGE_UPLOADER' not in os.environ:
        os.environ['IMAGE_UPLOADER'] = 'imgur'
    if 'IMGUR_CLIENT_ID' not in os.environ:
        os.environ['IMGUR_CLIENT_ID'] = 'dd2b80c72f01f10'

    data = open(filename, 'rb').read()

    url = None
    if os.environ['IMAGE_UPLOADER'] == 'imgur':
        url = post_imgur(data)
    else:
        sys.stderr.write('Unknown IMAGE_UPLOADER.\n')
        sys.exit(1)

    if url is None:
        sys.stderr.write('Upload failed.\n')
        sys.exit(1)

    print(url)


def post_imgur(data):
    url = 'https://api.imgur.com/3/image?type=file'
    headers = {
        'Authorization': 'Client-ID ' + os.environ['IMGUR_CLIENT_ID']
    }

    r = requests.post(url, data=data, headers=headers)
    # sys.stderr.write(r.text)
    response = r.json()
    if not response['success']:
        return None
    return response['data']['link']
