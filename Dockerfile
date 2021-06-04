FROM golang:1.16 as build
ADD . /healthchecker
WORKDIR /healthchecker
RUN go build -o healthchecker healthchecker.go

# Copy the scheduler binary to smaller container for deployment.
FROM gcr.io/distroless/base
WORKDIR /healthchecker
COPY --from=build /healthchecker/healthchecker /healthchecker/
ENTRYPOINT ["/healthchecker/healthchecker"]
