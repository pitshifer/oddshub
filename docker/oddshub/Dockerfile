FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY oddshub .

EXPOSE 80

CMD ["./oddshub"]