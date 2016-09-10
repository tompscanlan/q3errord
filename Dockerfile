# calls to slack url require ca-certificates to be installed
# hence base debian over scratch
FROM debian
#FROM scratch
#FROM golang:1.6
MAINTAINER Tom Scanlan <tscanlan@vmware.com>

EXPOSE 9999

# Add the microservice
ADD q3errord /q3errord

CMD ["/q3errord", "--port", "9999"]
