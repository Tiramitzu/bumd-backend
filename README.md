# Kemendagri SIPD Service Auth
Kemendagri SIPD Service Auth.-

## Prerequisites
Prequisites package:
* [Docker](https://www.docker.com/get-started) - for developing, shipping, and running applications (Application Containerization).
* [Go](https://golang.org/) - Go Programming Language.
* [Make](https://golang.org/) - Automated Execution using Makefile.
* [swag](https://github.com/swaggo/swag) Converts Go annotations to Swagger Documentation 2.0. We've created a variety of plugins for popular Go web frameworks.
* [golang-migrate/migrate](https://github.com/golang-migrate/migrate#cli-usage) Database migrations written in Go. Use as CLI or import as library for apply migrations.

Optional package:
* [gocritic](https://github.com/go-critic/go-critic) Highly extensible Go source code linter providing checks currently missing from other linters.
* [gosec](https://github.com/securego/gosec) Golang Security Checker. Inspects source code for security problems by scanning the Go AST.
* [golangci-lint](https://github.com/golangci/golangci-lint) Go linters runner. It runs linters in parallel, uses caching, supports yaml config, has integrations with all major IDE and has dozens of linters included.

## ⚡️ Quick start
These instructions will get you a copy of the project up and running on docker container and on your local machine.
1. Install Prequisites and optional package to your system:
2. Rename `Makefile.example` to `Makefile` then fill it with your make setting.
3. Generate swagger api doc by this command
```shell
make swag
```
4. Instant run by this command
```shell
make instant_run
```
5. Bulid go binary file
```shell
make build
```
6. Build go binary file and run
```shell
make run
```
7. Run in docker container
```shell
make docker_run
```

## Akun Testing
```
Password 123456 hash: $2a$12$e8dx90LgZVq7/nF6DHHurOEZgYZlZwfyKvoUjNFnmaOmz/6hMF7Fu

Admin
Username: 000000000000000000
Password: 123456

BUD
Username: 197308241992031001
Password: 123456

PA
Username: 196411271990031002 
Password: 123456

KPA
Username: 196502231992031001
Password: 123456

PPTK
Username: 196502231992031001
Password: 123456 

```

## Test Notes

### go-critic failed test
- commentedOutCode: //some comment
- example failed code
```
#utils/validator.go 37
if regexp.MustCompile(`^[a-zA-Z0-9/*-]*$`).MatchString(fieldstr) {
    return true
}else{
    return false
}
```

### gosec failed test
```
h := md5.New()
log.Println(h)
```

### lint failed test
```
#controllers/menu.go 143
if _, found := data[name]; found {
    data[name] = append(data[name], role)
} else {
    data[name] = []string{role}
}
```

## AUTH References---
- [Bcrypt Hash Generator](https://bcrypt-generator.com) Online Bcrypt Hash Generator & Checker.
