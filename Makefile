deploy:
	appcfg.py --no_cookies update app/app.yaml subservice/subservice.yaml

run:
	dev_appserver.py --enable_watching_go_path=false app/app.yaml subservice/subservice.yaml

test:
	ginkgo -r src test

prod-logs:
    open https://console.cloud.google.com/logs/viewer?project=datastore-sandbox-1114

.PHONY: deploy run test prod-logs