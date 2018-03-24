FROM alpine:latest
RUN apk --no-cache add curl nodejs

# copy over web application
ADD . /app/

# run webpack
RUN npm i quasar-cli yarn -g
WORKDIR /app/quasar/
RUN yarn
RUN quasar build -c

# debug dafuq is going on?
RUN ls -lrt dist/spa-mat/js/
RUN cat dist/spa-mat/js/app.*

# copy html files over to the public directory for serving
RUN mv dist/spa-mat/* ../public/

WORKDIR /app
ENTRYPOINT ["/app/server"]
