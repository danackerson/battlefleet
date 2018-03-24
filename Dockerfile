FROM alpine:latest
RUN apk --no-cache add curl nodejs

# copy over web application
ADD . /app/

# run webpack
RUN npm i quasar-cli yarn -g
WORKDIR /app/quasar/
RUN yarn
RUN quasar build -c

# copy html files over to the public directory for serving
RUN mv /app/quasar/dist/spa-mat/* /app/public/

WORKDIR /app
ENTRYPOINT ["/app/server"]
