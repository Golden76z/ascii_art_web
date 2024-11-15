# syntax=docker/dockerfile:1

# specify the base image to  be used for the application, alpine or ubuntu

FROM golang:1.22-alpine

LABEL author="Team"
LABEL description="Description"

# create a working directory inside the image

WORKDIR /app

# copy directory files i.e all files ending with .go

# copy all files/folder into /app
COPY . ./

# download Go modules and dependencies

# RUN go mod download

# compile application

RUN go build -o bin .
RUN apk update
RUN apk add --no-cache bash 

ENTRYPOINT [ "/app/bin" ]

# tells Docker that the container listens on specified network ports at runtime

EXPOSE 8080

# command to be used to execute when the image is used to start a container

CMD [ "./ascii_web" ]

#docker container ps (liste tout les containers de mon image)
#Mon image stocke tout mes containers

#docker container stop <container-id> (stopper mon container en cours d'execution)