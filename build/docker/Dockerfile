##########################
# Build stage
##########################

ARG DIR=/go/src/github.com/influenzanet/counter-service

FROM golang:1.17-alpine as builder
ARG DIR # Populate global ARG for this build
ARG VERSION=Unknown
ARG COMMIT=Unknown
ARG GOPKG=github.com/influenzanet/counter-service
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN mkdir -p ${DIR}
ENV GO111MODULE=on
ADD . ${DIR}
WORKDIR ${DIR}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR $DIR/cmd/counter-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -X ${GOPKG}/version.Version="${VERSION}" -X ${GOPKG}/version.Revision="${COMMIT}"' -o app .

##########################
# Image 
##########################
FROM scratch
ARG DIR
COPY --from=builder $DIR/cmd/counter-service/app /app/
WORKDIR /app
ENV GIN_MODE=release
ENV PORT=5021
EXPOSE 5021:5021
CMD ["./app"]
