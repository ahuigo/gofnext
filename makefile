
msg?=
.ONESHELL:
gitcheck:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi

.ONESHELL:
pkg: gitcheck test
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	git commit -am "$(msg)"
	#jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v" && git push origin HEAD

test:
	go test -coverprofile coverage.out -failfast ./...

cover:
	go tool cover -html=coverage.out
