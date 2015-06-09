#!/usr/bin/env python
"""Integration/smoke tests for ims, designed to be run against the official
`testimages` dataset.
"""
# stdlib imports
import collections
import random
import urllib

# third-party imports
import grequests


IMAGE_PATHS = [
    "/testimages/apple.jpg",
    "/testimages/chocroll.jpg",
    "/testimages/didimissit.gif",
    "/testimages/doilookfunny.jpg",
    "/testimages/fireworks.gif",
    "/testimages/ihateyouredits.JPG",
    "/testimages/lemur.jpg",
]


class FilterGenerator(object):

    @classmethod
    def _all_filters(cls):
        yield "brightness", str(random.choice(range(-100, 101)))
        yield "contrast", str(random.choice(range(-100, 101)))
        yield "fliphorizontal", ""
        yield "flipvertical", ""
        yield "hue", str(random.choice(range(-180, 181)))
        yield "resize", ','.join([str(random.choice(range(0, 1001))) for _ in range(2)])
        yield "rotate", random.choice(["90", "180", "270"])
        yield "saturation", str(random.choice(range(-100, 501)))
        yield "transpose", ""
        yield "transverse", ""

    @classmethod
    def _random_filter(cls):
        return random.choice(list(cls._all_filters()))

    @classmethod
    def generate(cls, length):
        filters = collections.OrderedDict()
        for _ in xrange(length):
            k, v = cls._random_filter()
            filters[k] = v
        return urllib.urlencode(filters)


def _generate_urls(base_url, num_urls):
    """Generate a list of URLs we can use for testing
    """
    for _ in xrange(num_urls):
        _weights = [1] * 80 + [2] * 10 + [3] * 5 + [4] * 3 + [5] * 2
        _image_path = random.choice(IMAGE_PATHS)
        _filters = FilterGenerator.generate(random.choice(_weights))
        yield base_url + _image_path + "?" + _filters


def _check_response(response, expect=200):
    """Check that a HTTP response is as expected.
    """
    if not response.status_code == expect:
        raise AssertionError("Response status not as expected: Expected {0}, Got {1}".format(
            str(expect), str(response.status_code)
        ))
    if response.status_code == 200 and len(response.content) == 0:
        raise AssertionError("Response has 200 status but did not return any data")


def run_tests():
    """Run integration tests (make a whole bunch of HTTP requests to a local server)
    """
    _url_generator = _generate_urls("http://localhost:5995", 100)
    _req_generator = (grequests.get(url) for url in _url_generator)
    for response in grequests.imap(_req_generator, size=5):
        _check_response(response)


if __name__ == '__main__':

    run_tests()
