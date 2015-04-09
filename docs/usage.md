Using the ims API
=================

This section of the documentation provides details on how to use the ims API,
what filters are available, and how to use them.

**Contents**:

- [Available filters](#available-filters)
- [Chaining filters](#chaining-filters)



# Available Filters

Assuming you have an image called `apple.png` being served by ims, applying
filters is simple.


## Resizing

Resize an image using arbitrary dimensions, maintaining (or not) the aspect
ratio.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/resize_200_0.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?resize=200,0`
    - `http://localhost:5995/apple.png?resize=200,0,nearestneighbor`

- **URL Parameters** (in order):
    - **width** *(int)*: required.
        - Desired width of output image `0 - 10000`.
        - Can be `0` to preserve aspect ratio only if height specified.
    - **height** *(int)*: required.
        - Desired height of output image `0 - 10000`.
        - Can be `0` to preserve aspect ratio only if width specified.
    - **resampling algorithm** *(string)*: optional.
        - Allowed values:
            - box
            - cubic
            - lanczos *(default)*
            - linear
            - nearestneighbour


## Cropping

Crop an image to the specified dimensions/coordinates.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/crop_0_0_100_200.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?crop=0,0,100,200`
    - `http://localhost:5995/apple.png?crop=50,50,200,200`

- **URL Parameters** (in order):
    - **x1** *(int)*: required.
        - `x` coordinate of top-left corner of box to use for cropping.
    - **y1** *(int)*: required.
        - `y` coordinate of top-left corner of box to use for cropping.
    - **x2** *(int)*: required.
        - `x` coordinate of bottom-right corner of box to use for cropping.
    - **y2** *(int)*: required.
        - `y` coordinate of bottom-right corner of box to use for cropping.


## Rotation

Rotate an image by 90, 180 or 270 degrees anti-clockwise.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/rotate_270.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?rotate=270`

- **URL Parameters** (in order):
    - **resampling algorithm** *(int)*: required.
        - Allowed values:
            - 90
            - 180
            - 270


## Flipping

Flip an image, horizontally or vertically.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/flipvertical.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?fliphorizontal`
    - `http://localhost:5995/apple.png?flipvertical`


## Transpose/Transverse

Transpose an image (flip horizontally and rotate 90 degrees counter-clockwise)
and get its transverse (flipped vertically and rotated 90 degrees
counter-clockwise).

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/transpose.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?transpose`
    - `http://localhost:5995/apple.png?transverse`


## Adjust contrast

Modify the contrast of an image.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/contrast_-50.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?contrast=50`
    - `http://localhost:5995/apple.png?contrast=-50`

- **URL Parameters** (in order):
    - **percentage** *(int)*: required.
        - Allowed range `-100 - 100`
        - `-100` gives a solid grey image.
        - `0` gives the original image.
        - `100` gives an overcontrasted image.


## Adjust brightness

Modify the brightness of an image.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/brightness_-50.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?brightness=50`
    - `http://localhost:5995/apple.png?brightness=-50`

- **URL Parameters** (in order):
    - **percentage** *(int)*: required.
        - Allowed range `-100 - 100`
        - `-100` gives a solid black image.
        - `0` gives the original image.
        - `100` gives a solid white image.


## Adjust saturation

Modify the saturation of an image.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/saturation_50.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?saturation=50`
    - `http://localhost:5995/apple.png?saturation=-50`

- **URL Parameters** (in order):
    - **percentage** *(int)*: required.
        - Allowed range `-100 - 500`
        - `0` gives the original image.


## Adjust hue

Modify the hue of an image.

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/hue_100.png)

- **Example URLs**:
    - `http://localhost:5995/apple.png?hue=50`
    - `http://localhost:5995/apple.png?hue=-50`

- **URL Parameters** (in order):
    - **shift** *(int)*: required.
        - Allowed range `-180 - 180`
        - `0` gives the original image.



# Chaining filters

Filters can be chained arbitrarily in the URL, they will be applied in the
order they are specified, e.g.

```
http://localhost:5995/apple.png?resize=500,0&crop=0,0,100,200&resize=100,100
```

Original image | Filtered image
--- | ---
![original](https://storage.googleapis.com/ims-examples/apple.png) | ![filtered](https://storage.googleapis.com/ims-examples/chained.png)
