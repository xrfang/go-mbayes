all: clean
	go build -ldflags '-s -w' .
	ln -s train untrain
clean:
	rm -f train untrain
