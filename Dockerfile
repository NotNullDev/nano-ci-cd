
# build go image
FROM golang:1.18 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY ./go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
RUN go build -v -o server .


# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds

FROM docker:20.10-git

RUN apk update
# RUN apk upgrade
RUN apk add bash

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

EXPOSE 3000

CMD ["/app/server"]

