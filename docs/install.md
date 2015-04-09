Installing and running ims
==========================

This section of the documentation provides details on how to install, configure
and run an ims instance.

**Contents**:

- [Installation](#installation)
    - [Local install](#local-install)
    - [Using Docker](#using-docker)
- [Configuration](#configuration)



# Installation


## Local install

ims is distributed in a single os-specific binary for each target platform. To
install ims just [download the executable](https://github.com/paddycarey/ims/releases/)
for your target platform, save somewhere on your `$PATH` (e.g. `/usr/local/bin/ims`),
and make it executable `chmod a+x /usr/local/bin/ims`.

ims has no required dependencies, you can run it straight away:

```bash
$ cd my_directory_full_of_images
$ ims <optional-configuration-options>
```

Once running, all images within the folder should be available at
http://localhost:5995 (e.g. http://localhost:5995/apple.png).


## Using Docker

If you prefer, a prebuilt docker image is available from the docker hub.
Provided you have Docker installed and running, you can run ims like so:

```bash
docker run -t -i -p 5995:5995 -v /path/to/your/images:/mnt/images paddycarey/ims:latest <optional-configuration-options>
```

Once running, all images within the folder should be available at
http://localhost:5995 (e.g. http://localhost:5995/apple.png).



# Configuration


## Command line arguments

ims is configured using command line arguments. The full range of available
options are shown below.

```bash
$ ims --help
ims.

Usage:
  ims [--storage=<src>] [--storage-credentials=<creds>] [--cache=<cch>]
      [--address=<address>] [--log-level=<level>] [--no-optimization]
  ims -h | --help
  ims --version

Options:
  -h --help                      Show this screen.
  --version                      Show version.
  --storage=<src>                Storage backend                    [default: ./].
  --storage-credentials=<creds>  Storage credentials file           [default: storage-credentials.json].
  --cache=<cch>                  Cache backend                      [default: ./.cache].
  --address=<address>            Address that ims should bind to    [default: :5995].
  --log-level=<level>            Log level (debug/info/warn/error)  [default: info].
  --no-optimization              Disables image optimization.
```


## Configuring storage

ims supports multiple storage backends from which to load source images. By
default ims will serve images from the current working directory, but this is
configurable using the `--storage` flag.

#### Local storage

Local filesystem storage is the simplest option. Simply point ims at the
directory containing your images like so:

```bash
$ ims --source=./
$ ims --source=/some/directory
```

#### Google Cloud Storage

ims supports reading images from a Google Cloud Storage bucket. You'll need to
create a service account using the Google cloud console, and ensure that the
Google Cloud Storage API is enabled for the project.

When creating the service account you should download the credentials in json
format and copy the values into your environment so that ims can access them.

The keys have the same names as the keys in the downloaded JSON file in
uppercase, prepended with `GCS_`. [Forego](https://github.com/ddollar/forego)
is a nice tool for automating this process during development. You can run ims
with cloud storage enabled like so:

```bash
# set the required environment variables
$ export GCS_PRIVATE_KEY_ID="0084nwejcgweuif7wleukcgw93"
$ export GCS_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\ngesrvberswefserf.....\n-----END PRIVATE KEY-----\n"
$ export GCS_CLIENT_EMAIL="123123123-xlaksxmlaksxmalskxm@developer.gserviceaccount.com"
$ export GCS_CLIENT_ID="123123123-xlaksxmlaksxmalskxm.apps.googleusercontent.com"

# run ims using the root of the GCS bucket
$ ims --storage=gcs://my-bucket

# or serve from a specific folder
$ ims --storage=gcs://my-bucket/my-folder
```
