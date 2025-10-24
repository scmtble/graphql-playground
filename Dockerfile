FROM golang:1.25-trixie AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -ldflags "-s -w" -o /bin/app ./cmd

FROM gcr.io/distroless/static-debian13:latest
COPY --from=build /bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]