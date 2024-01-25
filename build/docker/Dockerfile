##########################
# Build stage
##########################

ARG DIR=/go/src/github.com/influenzanet/counter-service

FROM golang:1.19-alpine as builder
ARG DIR # Populate global ARG for this build
ARG VERSION=Unknown
ARG COMMIT=Unknown
ARG GOPKG=github.com/influenzanet/counter-service/pkg/version
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN mkdir -p ${DIR}
ENV GO111MODULE=on
ADD . ${DIR}
WORKDIR ${DIR}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR ${DIR}/cmd/counter-service
ARG LDFLAG="-extldflags '-static' -X ${GOPKG}.Version=${VERSION} -X ${GOPKG}.Revision=${COMMIT}"
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "${LDFLAG}" -o app .
RUN echo "${VERSION};${COMMIT}" > version
##########################
# Image 
##########################
FROM scratch
ARG DIR
COPY --from=builder ${DIR}/cmd/counter-service/app /app/
COPY --from=builder ${DIR}/cmd/counter-service/version /app/
WORKDIR /app
ENV GIN_MODE=release
ENV PORT=5021
EXPOSE 5021:5021
CMD ["./app"]
