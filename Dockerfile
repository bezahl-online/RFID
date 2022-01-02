FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64\
    RFID_PORT_NAME=/dev/serial/by-id/usb-1a86_USB2.0-Ser_-if00-port0 \
    PRODUCTIVE=YES

#RUN apk update && apk upgrade && \
#    apk add --no-cache bash git openssh
# Move to working directory /build
WORKDIR /rfid

# Copy and download dependency using go mod
#COPY go.mod .
#COPY go.sum .
#RUN go mod download

# Copy the build image into the container
# need to build like this:
# $ CGO_ENABLED=0 go build -o gm65server
ADD rfidserver .
ADD localhost.crt .
ADD localhost.key .

# Build the application
#RUN go build -o server .

# Export necessary port
EXPOSE 8040

# Command to run when starting the container
CMD ["/rfid/rfidserver"]