FROM golang:1.17-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY main.go kubernetes.go handler.go /app/
RUN go build -o=main

# =====================================================

FROM alpine AS production
RUN apk add --no-cache musl-dev
# Copy app and front
WORKDIR /app
COPY ./templates templates

COPY --from=build /app/main /main

EXPOSE 3000
ENTRYPOINT ["/main"]