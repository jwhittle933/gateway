# Gateway
This package is for building an API Gateway with ease, and can be used with both the `net/http` std lib package and with `fasthttp`.

## Reverse Proxies
A reverse proxy moves a request from one place to another and back again, and is done for many reasons (auth, micro-services, etc.). Data transforms in the middle should normally be kept to a minimum, though there are times when this is good and appropriate. The `Gateway` package in this repo is designed to make gateway building very easy. 

## Approaches
There are two ways to use this package to get a Gateway off the ground: custom implementation or the pre-built, configurable server. The package is like any other Go package: you can use what you want and embed it in any application of your choosing. The server, included under `cmd`, was built to give development teams or individuals a quick way to gateway their services. It is built using the technologies and methodologies that we find most compelling (namely `fasthttp` rather than `net/http`).