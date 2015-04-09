ims
===

ims is a standalone image manipulation, optimisation and serving service
written in Go. ims provides on-the-fly resizing, cropping, rotation etc. (See
full docs for an explanation of all available filters).

ims uses the excellent [gift](https://github.com/disintegration/gift) library
under the hood to perform all image manipulation.


## Work in progress

ims is very much a work in progress and isn't ready for production use. Please
treat it as such and file issues where appropriate. As always, pull requests or
other offers of help are greatly appreciated.


## Documentation

Full documentation is available at the following pages:

- [Installing and running ims](docs/install.md)
- [Using the ims API](docs/usage.md)


## Optimisation

When available and configured to do so (on by default), ims will use one of a
number of third-party tools to optimise the images being served. Optimisation
is applied after the image has been processed by ims. Installation of the
third-party tools is outside the scope of this README.

ims uses the following tools when available:

- GIF: [gifsicle](http://www.lcdf.org/gifsicle/)
- JPEG: [jpegtran](http://jpegclub.org/jpegtran/)
- PNG: [optipng](http://optipng.sourceforge.net/)


## Caching

ims provides a simple on-disk cache, ensuring that unnecessary encoding work
isn't repeated. When the first request comes in for a given transformation, the
resulting image is written to cache and served from there if a subsequent
request comes in.

ims' caching implementation is very naive at present and has no concept of
expiration times or detection when the source image changes. It will be
expanded to include these features in future.


## Copyright & License

- Copyright © 2015 Patrick Carey (https://github.com/paddycarey)
- Licensed under the **MIT** license.
