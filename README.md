# go-gcrmirrors

auto generate https://github.com/kbcx/mirrors.kb.cx api json from https://github.com/x-mirrors/gcr.io

## Environment Variables

- SOURCE_DIR: https://github.com/x-mirrors/gcr.io path
- PUBLIC_DIR: json public dir

## How to Use

```
    - name: Generate gcr.io mirrors api json
      uses: x-actions/go-gcrmirrors@main
      env:
        SOURCE_DIR: "<path-of>/gcr.io"
        PUBLIC_DIR: "<publish-dir>"
```
