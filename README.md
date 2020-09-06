# s3-gpg-proxy-server

Simple solution for adding asymmetric encryption layer to your backup scripts (or something else)

## Why I need this?
Some applications (backup solutions for example) unable to work with GPG,
but can save results to S3-compatible server.
This applicaion provides additional local security layer that encrypts _locally_ all of received files
before sending to target

## How to use

1. `cp example_settings.env settings.env`
2. Set correct credentials to target S3 server in `settings.env`
3. Run containers with `docker-compose up`
4. Connect to proxy server on http://127.0.0.1:9001 with `ACCESS_KEY=minio` and `SECRET_KEY=miniostorage`.


## FAQ

### Why do I need asymmetric encryption?
With asymmetric encryption you can use one single key for any purposes without sacrificing security.

It also means that even if even if an attacker gains access to this server,
he will not be able to decrypt old data objects..

### Is it secure?
It is containerized application,
only encrypted data may be uploaded to target server.
And no one can access this containers from the Internet.

### I never used Docker. Where I can get it?
You need Docker [Engine](https://docs.docker.com/engine/install/) and Docker [Compose](https://docs.docker.com/compose/install/).
It is possible to install this tools on almost any existing server,
even on RaspberriPi.

### Is it production ready?
Quick answer is **no**.
Use it at your own risk.

But realization is pretty simple inside.
Should work stably.
Even if encryption part fails it is possible to recover initial,
unencrypted data from local Minio server instance.

