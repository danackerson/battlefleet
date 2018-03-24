# Quasar App

> To test locally:

$ docker run -it --rm -p 8443:8443 -v /Users/ackersond/go/src/github.com/danackerson/battlefleet:/battlefleet alpine
/ # apk update && apk --no-cache add curl nodejs
/ # npm i quasar-cli yarn -g
/ # cd battlefleet/quasar/
/ # yarn
/ # quasar dev

!WARNING : fetching/building node_modules on this MacOS mounted volume
from Alpine Linux is stupid!

https://localhost:8443/

/ # quasar build
/ # mv mv dist/spa-mat/* ../public/
