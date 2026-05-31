FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/contractgen ./cmd/cgen-web

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/contractgen /contractgen

EXPOSE 8080
ENTRYPOINT ["/contractgen"]