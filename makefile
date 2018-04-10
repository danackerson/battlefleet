# Makefile for MacOS v10.13.14
staging:
	cd quasar; \
	gsed -i 's/PORT: 443/PORT: 8083/' quasar.conf.js && quasar build; \
	rm -Rf ../public/* && mv dist/spa-mat/* ../public/; \
	cd ../public; \
	sed -i -e 's/href=app/href=css\/app/' index.html; \
	sed -i -e 's/href=\/app/href=\/css\/app/' index.html; \
	mkdir css && mv app*.css css/; \
	cd ../ && fresh -c fresh.conf

.PHONY: staging
