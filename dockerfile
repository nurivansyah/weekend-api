FROM golang:1.18 as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.5
COPY . .

RUN swag init
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian11
COPY --from=builder /go/bin/app /
EXPOSE 8080
CMD [ "/app" ]