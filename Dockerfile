FROM golang:1.20.0-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/moai cmd/moai/**.go
USER nobody

# ---

FROM scratch
COPY --from=build /app/bin /
COPY --from=build /etc/passwd /etc/passwd
USER nobody
EXPOSE 8080 8443
ENV MOAI_DOCKER_HOST 0.0.0.0
CMD ["/moai"]
