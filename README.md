ims
===

ims is a standalone image manipulation and serving service written in Go. ims
provides on-the-fly resizing, cropping, rotation etc. (See full docs for an
explanation of all available filters).

ims uses the excellent [gift](https://github.com/disintegration/gift) library
under the hood to perform all image manipulation.



## Work in progress

ims is very much a work in progress and isn't ready for production use. Please
treat it as such and file issues where appropriate. As always, pull requests or
other offers of help are greatly appreciated.



## Usage

Until a release build of ims is made, you'll need to build it if you want to use it:

```
$ go build
```

Once built, you just need to point the binary at an existing directory full of images.

```
$ ./ims --source ./myimages
```

And then all images within the folder should be available at
http://localhost:5995. Assuming you have an image called `apple.png`, applying
filters is simple.


### Resizing

To resize an image you can use a URL like:

```
http://localhost:5995/apple.png?resize=500,0
```

Where the first parameter is the desired width, and the second is the desired
height. Either value can be set to `0` (not both) to maintain aspect ratio
during resize.


### Cropping

To crop an image you can use a URL like:

```
http://localhost:5995/apple.png?crop=0,0,100,200
```

Where the first 2 parameters are the top left corner of the box to use for
cropping, the second 2 values are the bottom corner of the box.


### Rotation

You can rotate an image by 90, 180 or 270 degrees anti-clockwise using a URL in
the following format:

```
http://localhost:5995/apple.png?rotate=270
```

The only acceptable parameters for `rotate` are `90`, `180` and `270`.


### Flipping

You can flip an image (horizontally or vertically) like so:

```
http://localhost:5995/apple.png?fliphorizontal
http://localhost:5995/apple.png?flipvertical
```

Any parameters passed to either `fliphorizontal` or `flipvertical` are ignored.


### Transpose/Transverse

You can transpose an image (flip horizontally and rotate 90 degrees
counter-clockwise) and get its transverse (flipped vertically and rotated 90
degrees counter-clockwise) like so:

```
http://localhost:5995/apple.png?transpose
http://localhost:5995/apple.png?transverse
```

Any parameters passed to either `transpose` or `transverse` are ignored.


### Adjust contrast

The contrast of an image can be manipulated like so:

```
http://localhost:5995/apple.png?contrast=50
```

The percentage parameter must be in range (-100, 100). The percentage = 0 gives
the original image. The percentage = -100 gives solid grey image. The
percentage = 100 gives an overcontrasted image.


### Adjust brightness

The brightness of an image can be manipulated like so:

```
http://localhost:5995/apple.png?brightness=50
```

The percentage parameter must be in range (-100, 100). The percentage = 0 gives
the original image. The percentage = -100 gives solid black image. The
percentage = 100 gives solid white image.


### Adjust saturation

The saturation of an image can be manipulated like so:

```
http://localhost:5995/apple.png?saturation=50
```

The percentage parameter must be in range (-100, 500). The percentage = 0 gives
the original image.


### Adjust hue

The hue of an image can be manipulated like so:

```
http://localhost:5995/apple.png?hue=150
```

The shift parameter is the hue angle shift, typically in range (-180, 180). The
shift = 0 gives the original image.


### Chaining filters

Filters can be chained arbitrarily in the URL like so:

```
http://localhost:5995/apple.png?resize=500,0&crop=0,0,100,200&resize=100,100
```

Filters will be applied in the order they are specified.



## Caching

ims provides a simple on-disk cache, ensuring that unnecessary encoding work
isn't repeated. When the first request comes in for a given transformation, the
resulting image is written to cache and served from there if a subsequent
request comes in.

ims' caching implementation is very naive at present and has no concept of
expiration times or detection when the source image changes. It will be
expanded to include these features in future.


## TODO

- Docs
- Expose more manipulations/filters
- Unit/integration tests
- Release testing builds



## Copyright & License

- Copyright © 2015 Patrick Carey (https://github.com/paddycarey)
- Licensed under the **MIT** license.
