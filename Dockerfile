FROM golang:1.25-trixie AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -ldflags "-s -w" -o ./graphql-playground ./cmd

FROM gcr.io/distroless/static-debian13:latest
COPY --from=build /src/graphql-playground /graphql-playground
EXPOSE 8080
ENTRYPOINT ["/graphql-playground"]