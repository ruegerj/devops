FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy

# toggle for coverage instrumentation of e2e tests
ARG WITH_COVERAGE=false
RUN if [ "$WITH_COVERAGE" = "true" ]; then \
    go build -cover -coverpkg=./... -o bin/api .; \
    else \
    go build -o bin/api; \
    fi


FROM alpine:3.22

WORKDIR /app
EXPOSE 3000
CMD [ "./api" ]
ENV HOST=0.0.0.0
ENV PORT=3000
ARG JWT_KEY
ENV JWT_KEY=${JWT_KEY}
COPY --from=builder /app/bin/api .