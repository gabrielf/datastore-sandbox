deploy:
	appcfg.py --no_cookies update app/

run:
	dev_appserver.py app/app.yaml

test:
	ginkgo -r src test

.PHONY: deploy run test