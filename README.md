# git-version

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
        gitCommit = "e1ed2c7b8b1d1af5c4213c87f7fe3a8062a5eeb8"
        gitShortCommit = "e1ed2c7"
        gitTag = "NA"
        gitBranch = "NA"
        gitStatus = "dirty"
)
```
