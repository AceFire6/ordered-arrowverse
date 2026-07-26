# syntax=docker/dockerfile:1.7
#
# Multi-stage build for the Arrowverse Go service.
#
# Stage 1 (build) compiles a stripped static binary. Tools (templ + sqlc)
# run from a tools/go.mod module so the application go.mod stays free of
# codegen dependencies.
#
# Stage 2 (runtime) ships the binary on a distroless non-root image.

ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

# Install only what we need for codegen + go build.
RUN apk add --no-cache git make bash ca-certificates

# Pull dependencies first for layer caching.
COPY go.mod go.sum ./
COPY tools/go.mod tools/go.sum ./tools/
RUN go mod download && cd tools && go mod download

# Install templ + sqlc binaries.
RUN go install github.com/a-h/templ/cmd/templ@v0.3.943 \
 && go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

# Bring in the rest of the source.
COPY . .

# Generate templ + sqlc code, then build the stripped binary.
RUN templ generate -path ./components
RUN sqlc generate -f ./sqlc.yaml
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/arrowverse \
    ./cmd/ordered-arrowverse

FROM gcr.io/distroless/static-debian12:nonroot AS runtime
COPY --from=build /out/arrowverse /usr/local/bin/arrowverse
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot
EXPOSE 8080
ENV PORT=8080

# Container-level healthcheck is provided by compose.production.yaml
# (the distroless base has no shell, so the Dockerfile form is omitted).

ENTRYPOINT ["/usr/local/bin/arrowverse"]
