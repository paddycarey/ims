#!/usr/bin/env python
import base64
import os


def decode_env_var():
    v = os.environ.get("GCS_DATA", None)
    return base64.urlsafe_b64decode(v)


def write_file(path, data):
    with open(path, "w") as f:
        f.write(data)


if __name__ == "__main__":

    print "Writing .env file"
    data = decode_env_var()
    path = os.path.join(os.path.dirname(__file__), os.pardir, ".env")
    write_file(path, data)
