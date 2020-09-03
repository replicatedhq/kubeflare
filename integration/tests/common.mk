SHELL := /bin/bash

.PHONY: commonrun
commonrun:
	./before.sh
	../../../bin/integration run \
		--spec ./spec.yaml \
		--email $(CF_API_EMAIL) \
		--key $(CF_API_KEY) \
		--zone-name $(CF_ZONE_NAME)
	rm -rf out
	mkdir -p out
	./after.sh > out/result.json
	./verify.sh
	