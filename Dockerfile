FROM alpine:latest
RUN apk --no-cache add curl

# copy over web application
ADD server /app/
ADD public /app/public

WORKDIR /app
ENTRYPOINT ["/app/server"]
