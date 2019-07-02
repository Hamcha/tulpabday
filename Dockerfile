FROM golang:alpine AS builder

ARG MODULE_PATH

ENV GOPROXY https://modules.fromouter.space
ENV GO111MODULE=on

WORKDIR /app

# Get updated modules
COPY ./ ./
RUN cd src; go mod download

# Compile code
RUN CGO_ENABLED=0 go build -o /svc ./src

FROM scratch AS final

# Import the compiled executable from the first stage.
COPY --from=builder /svc /svc
COPY --from=builder /app /app

WORKDIR /app

CMD [ "/svc" ]