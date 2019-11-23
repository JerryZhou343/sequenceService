.PHONY: all clean

output= sequencesvc


all: clean
	go build -o ./bin/${output} main.go version.go
clean:
	rm -f ./bin/${output}
