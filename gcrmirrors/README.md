# gcrmirrors

## Build

```
git clone https://github.com/x-actions/go-gcrmirrors.git
cd go-gcrmirrors/gcrmirrors

make build/linux/mac/windows
or
GOOS=linux GOARCH=amd64 go build -tags netgo
```

## Usage

```
./gcrmirrors
Usage of ./gcrmirrors:
  -publicDir string
        json public dir (default "./public")
  -sourceDir string
        https://github.com/kbcx/gcr.io dir
```

demo:

```
./gsync \
  -sourceDir <some-dir>
```
