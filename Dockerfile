FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN apt-get update && apt-get install -y openssl
RUN openssl genpkey -algorithm ED25519 -out refresh.ed
RUN openssl pkey -pubout -in refresh.ed -out refresh.ed.pub
RUN openssl genpkey -algorithm ED25519 -out auth.ed

RUN CGO_ENABLED=0 GOOS=linux go build -o /sportsvoting

FROM alpine:3.18 AS build-release-stage

WORKDIR /

COPY --from=build-stage /sportsvoting /sportsvoting
COPY --from=build-stage /app/*.ed* /

EXPOSE 8080

ENTRYPOINT ["/sportsvoting"]