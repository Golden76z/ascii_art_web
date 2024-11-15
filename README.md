ASCII Art Web
=====================

This is a simple web application that generates ASCII art from user input. The application uses the Go programming language and the `net/http` package to create a web server that listens on port 8080.

Features
--------

* Generates ASCII art from user input
* Allows users to select a banner style
* Displays the generated ASCII art on a result page
* Allows users to download the generated ASCII art as a text file

Routes
------

* `/`: Home page with text input and style selection
* `/error404`: Error 404 page
* `/ascii-art`: Generates ASCII art from user input
* `/result`: Displays the generated ASCII art
* `/download`: Downloads the generated ASCII art as a text file

Templates
---------

* `template/error404.html`: Error 404 page template
* `template/index.html`: Home page template
* `template/generator.html`: ASCII art generator template
* `template/result.html`: Result page template

Files
-----

* `static/`: Directory for static files (e.g. CSS, images)
* `template/`: Directory for HTML templates
* `text.txt`: Temporary file for storing generated ASCII art

How to Run
-----------

1. Run the application using `go run main.go`
2. Open a web browser and navigate to `http://localhost:8080`

Note: This application uses temporary files to store generated ASCII art. Make sure to clean up these files periodically to avoid 
disk space issues.