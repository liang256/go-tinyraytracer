# go-tinyraytracer

A lightweight ray tracing renderer written in Go, designed for simulating realistic lighting effects such as reflection, refraction, and shadows. This project demonstrates the fundamentals of ray tracing, material modeling, and rendering geometric shapes, using simple but extensible code.

## Output Example
The program generates a 1024x768 PPM image with spheres and discs rendered under multiple light sources. Here’s a preview of the result:
![disc](https://user-images.githubusercontent.com/23650308/194513765-0c17446a-7ee6-42cf-a9af-02861a76e0ee.png)

![candy](https://user-images.githubusercontent.com/23650308/194514517-1535d436-812e-49a2-82ca-59118fb6891d.png)

## Core Concepts
Ray Tracing Algorithm
1. Rays are cast from the camera through each pixel on the canvas.
2. For each ray:
    - Intersection: Check which object the ray hits first.
    - Lighting: Calculate diffuse and specular contributions from all light sources.
    - Reflection and Refraction: Recursively trace reflected and refracted rays.
Shading Model
- Diffuse Shading: Light intensity proportional to the angle between the light ray and surface normal.
- Specular Shading: Simulates shiny surfaces using the Phong reflection model.
- Reflection: Simulates light bouncing off reflective surfaces.
- Refraction: Simulates light bending through transparent materials.

## Installation
1. Clone the Repository:
```
git clone https://github.com/liang256/go-raytracer.git
cd go-raytracer
```
2. Install Dependencies:
```
go mod tidy
```
3. Run the Program:
```
go run main.go
```
4. Output: The program generates an image file out.ppm in the project directory. Open the file using any PPM-compatible viewer or convert it to a more common format like PNG using tools such as ImageMagick:
```
convert out.ppm out.png
```
