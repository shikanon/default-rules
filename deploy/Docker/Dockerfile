FROM golang:alpine3.15 as builder

COPY . /rules

WORKDIR /rules

ARG GOPROXY
ENV GOPROXY=$GOPROXY

RUN go build -o default-rule main.go 

FROM alpine:3.15.4

COPY --from=builder /rules/default-rule /rules/default-rule

RUN chmod +x /rules/default-rule

CMD [ "/rules/default-rule"]