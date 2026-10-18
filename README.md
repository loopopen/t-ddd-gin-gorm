# 

## This project is generated from template `t-ddd-gin-huma-gorm`.

## 1. Manage Go versions with goenv (recommended).

This project includes a `.go-version` file. We recommend using [goenv](https://github.com/go-nv/goenv) to manage multiple Go versions and automatically select the version required by this project.

After installing and configuring `goenv`, install the project-specific Go version:

```sh
goenv install "$(cat .go-version)"
go version
```

When you enter the project directory, `goenv` reads `.go-version` and switches to the matching Go version automatically.

## 2. Use make.

* Use make to build OpenAPI documents from the Huma route definitions.
```sh
make doc
```

* Use make to run your project.
```sh
make run

make run env=local
```

* Use make to wire dependency injection
```sh
make wire
```

* Use make to generate all.
```sh
make gen
```

## 3. Run and access the Huma API documentation at http://127.0.0.1:8080/docs

The OpenAPI documents are also available at `http://127.0.0.1:8080/openapi.json` and `http://127.0.0.1:8080/openapi.yaml`.

## 4. Delete file PLEASE_DELETE_ME.go and fix your code.

## 5. Create a new git repository of the command line:
```sh
git init
git add -A
git commit -m "init commit"
git branch -M main
git remote add origin git@github.com:<org-name>/<app-name>.git
git push -u origin main
```

---

* ## How to use gonew template?
### 1. Install gonew.
```sh
go install golang.org/x/tools/cmd/gonew@latest
```
### 2. Use this template to create your project.
```sh
mkdir <app-name> && \
gonew github.com/loopopen/t-ddd-gin-huma-gorm@<version> github.com/<org-name>/<app-name> ./<app-name>
make tidy
```

