FROM golang:1.17-alpine AS build

# RUN apk add --no-cache make
WORKDIR /go/src/bridge/
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/bridge .


FROM scratch

COPY --from=build /go/bin/bridge /bridge
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/bridge"]
