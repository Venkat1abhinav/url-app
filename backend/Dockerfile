FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /url-app ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /url-app /url-app
EXPOSE 4000
USER nonroot:nonroot
ENTRYPOINT ["/url-app"]
