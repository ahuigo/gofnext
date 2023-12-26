msg?=

######################### test ################
test: 
	go test -race -coverprofile cover.out -coverpkg "./..." -failfast ./...
cover: test
	go tool cover -html=cover.out
race: 
	go test -race -failfast ./...
fmt:
	gofmt -s -w .


###################### pkg ##########################
.ONESHELL:
gitcheck:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi

.ONESHELL:
pkg: gitcheck test fmt
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	git commit -am "$(msg)"
	#jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v" && git push origin HEAD
pkg0: test
	v=`cat version` && git tag "$$v" && git push origin "$$v" && git push origin HEAD
report:
	goreportcard-cli -v
