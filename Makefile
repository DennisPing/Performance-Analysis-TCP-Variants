.PHONY: all clean

PWD := $(shell pwd)

all: exp01 exp02 exp03

exp01:
	@cd cmd/exp01 && go build -o $(PWD)/bin/exp01 && echo Successful build exp01

exp02:
	@cd cmd/exp02 && go build -o $(PWD)/bin/exp02 && echo Successful build exp02

exp03:
	@cd cmd/exp03 && go build -o $(PWD)/bin/exp03 && echo Successful build exp03

clean:
	@rm -rf bin/*