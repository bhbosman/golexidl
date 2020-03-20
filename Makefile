all: clean build run_files

clean:
	rm -rf golexidl
	rm -f ./output/*.orb.tokens
build:
	go build

run_files:
	./golexidl -v -remove_token WhiteSpace,SingleComment  -o ./output/JacORB.orb.tokens ../../JacORB/JacORB/idl/omg/orb.idl
	./golexidl -v -remove_token WhiteSpace,SingleComment  -o ./output/mine.orb.tokens ../CodeGenerators/idlgenerator/corbaFiles/orb.idl
