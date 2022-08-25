# go-template-api

This is a template for a simple HTTP REST api server written in Go. It comes without any external dependencies.
Static files are embedded in the executable. They are organized in two directories:

### assets

Here are content like images, icons, stylesheets or Javascript files. They are delivered as they are. the mimetype is 
set based on the extension.

### templates

Here are the Go HTML templates. Go templates can contain variables and will be rendered with parameters.

## build

Checkout and build with `go get . && go build .` in the project folder.
