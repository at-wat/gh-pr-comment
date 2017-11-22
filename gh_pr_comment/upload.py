import sys
import os
import requests
import uuid
import mimetypes


def post_main():
    argv = sys.argv
    if len(argv) < 2:
        sys.stderr.write('usage: gh-pr-upload filename\n')
        sys.stderr.write('env:\n')
        sys.stderr.write('- IMAGE_UPLOADER: '
                         + 'imgur(default), '
                         + 's3(optional, need boto3)\n')
        sys.stderr.write('- ALLOW_PUBLIC_UPLOADER: '
                         + 'set it to enable public uploader\n')
        sys.stderr.write('env for imgur:\n')
        sys.stderr.write('- IMGUR_CLIENT_ID: custom-client-id\n')
        sys.stderr.write('env for s3:\n')
        sys.stderr.write('- AWS_DEFAULT_REGION\n')
        sys.stderr.write('- AWS_ACCESS_KEY_ID\n')
        sys.stderr.write('- AWS_SECRET_ACCESS_KEY\n')
        sys.stderr.write('- AWS_S3_BUCKET\n')
        sys.stderr.write('return: image url\n')
        sys.exit(1)

    post(argv[1])


def post(filename):
    if 'IMAGE_UPLOADER' not in os.environ:
        os.environ['IMAGE_UPLOADER'] = 'imgur'
    if 'IMGUR_CLIENT_ID' not in os.environ:
        os.environ['IMGUR_CLIENT_ID'] = 'dd2b80c72f01f10'

    basename, ext = os.path.splitext(filename)
    data = open(filename, 'rb').read()
    mime = mimetypes.guess_type(filename)

    url = None
    if os.environ['IMAGE_UPLOADER'] == 'imgur':
        url = post_imgur(data)
    elif os.environ['IMAGE_UPLOADER'] == 's3':
        url = post_s3(data, ext, mime[0])
    else:
        sys.stderr.write('Unknown IMAGE_UPLOADER.\n')
        sys.exit(1)

    if url is None:
        sys.stderr.write('Upload failed.\n')
        sys.exit(1)

    print(url)


def post_imgur(data):
    if 'ALLOW_PUBLIC_UPLOADER' not in os.environ:
        sys.stderr.write('Public uploader is not enabled.\n')
        sys.stderr.write('Set ALLOW_PUBLIC_UPLOADER to enable.\n')
        sys.exit(1)

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


def post_s3(data, ext, mime):
    try:
        import boto3
    except ImportError:
        sys.stderr.write('boto3 is not available.\n')
        sys.exit(1)

    if 'AWS_S3_BUCKET' not in os.environ:
        sys.stderr.write('AWS_S3_BUCKET is not specified.\n')
        sys.exit(1)

    if 'AWS_DEFAULT_REGION' not in os.environ:
        sys.stderr.write('AWS_DEFAULT_REGION is not specified.\n')
        sys.exit(1)

    bucket_name = os.environ['AWS_S3_BUCKET']
    path = str(uuid.uuid4()) + ext

    s3 = boto3.resource('s3')
    s3.Bucket(bucket_name).put_object(
        Key=path,
        Body=data,
        ACL='public-read',
        ContentType=mime)

    return 'https://s3-%s.amazonaws.com/%s/%s' \
        % (os.environ['AWS_DEFAULT_REGION'], bucket_name, path)
