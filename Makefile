build:
	go build -o main

clean:
	-rm main

test: clean build
	./main