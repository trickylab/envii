# Build stage
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/envii ./cmd/envii

# Runtime stage
FROM alpine:3.20
RUN adduser -D -h /home/envii envii
COPY --from=build /out/envii /usr/local/bin/envii
USER envii
WORKDIR /home/envii
ENTRYPOINT ["envii"]
