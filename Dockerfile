FROM golang:1.17.3 as build
WORKDIR /build
ADD . .
RUN CGO_ENABLED=1 go build -o kudos ./backend
FROM scratch
COPY --from=build /build/kudos .
ENTRYPOINT ["./kudos"]
