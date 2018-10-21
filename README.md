# go-latestver
get golang repos latest version format [v0.0.0-timestamp-commit], used in go.mod 

```
go get -u github.com/hysios/go-latestver
```

```
go-latestver https://github.com/hysios/go-latestver

%> v0.0.0-20181021184428-c2d8df786ab8
```

copy version and paste to go.mod require section, like this
```

require (
    github.com/hysios/go-latestver v0.0.0-20181021184428-c2d8df786ab8
)
```
