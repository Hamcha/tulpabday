all: dep tulpaweb

tulpaweb:
	go build -o tulpaweb ./src

dep:
	go get github.com/cloudflare/gokabinet/kc
	go get github.com/gorilla/mux
	go get github.com/hoisie/mustache
	touch dep

clean:
	rm -f dep
	rm -f tulpaweb
