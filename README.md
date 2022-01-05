# DISYS-Exam
Template project for the Distributed Systems Exam at the IT University of Copenhagen

## Guide to run this application
If you have docker you can then use the command `docker compose up` at the root of the application.
Everything should then start up and show in the terminal.

## Guide to generate compiled proto files
It is necessary to the following in order to compile the proto file into something that you can use in your GoLang project.

1. Run the following two go install commands
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

2. After that run this in your command-line or add it to your environment variables file
```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

3. Now you can add the stuff you want in the `./proto/exam.proto` file
4. Then you can run the shell script by typing the following in the command-line `sh gen.sh`