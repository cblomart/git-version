# git-version

REM: **WORK IN PROGRESS**

Generates a version.go file in function of the informations available in the git repository

To generate with "go generate" place the following line in your main package:

```golang
package main

//go:generate git-version
```

Generated "version.go" will contain constants with infomation from git

```golang
package main

const(
    GIT_COMMIT = "2428f6d35d68a8cad92c5c81c810a929adb89bcd"
    GIT_SHORT_COMMIT = "2428f6d"
    GIT_TAG = "v0.7"
)
```