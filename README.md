to run this project run this command from root

go run ./cmd/server/main.go

if you want change name of the root directory

first ensure you have defined the import path other modules correctly in main.go
and then run 

go mod init {directory name}

go mod tidy 

and then run the project using above command