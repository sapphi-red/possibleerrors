# possibleerrors

Golang linters for finding code which is likely a logic error.

- [fordirection](./fordirection): `for i := 0; i < 10; i--`
- [avoidaccesslen](./avoidaccesslen): `slice[len(slice)]`
- [mutexscope](./mutexscope): forgotten `.Unlock()`.

## Install
```
go get -u github.com/sapphi-red/possibleerrors/fordirection/cmd/fordirection
go get -u github.com/sapphi-red/possibleerrors/avoidaccesslen/cmd/avoidaccesslen
go get -u github.com/sapphi-red/possibleerrors/mutexscope/cmd/mutexscope
```

Or download from [releases](https://github.com/sapphi-red/possibleerrors/releases) and put it in PATH.

## Usage
```
go vet -vettool=`which fordirection` pkgname
go vet -vettool=`which avoidaccesslen` pkgname
go vet -vettool=`which mutexscope` pkgname
```

## Usage with GitHub Actions
### With go vet + reviewdog
```yaml
name: reviewdog
on: [pull_request]
jobs:
  govet:
    name: runner / govet
    runs-on: ubuntu-latest
    steps:
      - name: Check out code into the Go module directory
        uses: actions/checkout@v1
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.15'
      - name: build custom analyzer
        run: |
          go get -u github.com/gostaticanalysis/vetgen
          vetgen init _vettool
          cd _vettool
          vetgen add github.com/sapphi-red/possibleerrors/fordirection
          vetgen add github.com/sapphi-red/possibleerrors/avoidaccesslen
          vetgen add github.com/sapphi-red/possibleerrors/mutexscope
          go install .
      - uses: reviewdog/action-setup@v1
      - name: Run reviewdog
        env:
          REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          go vet -vettool=$(which _vettool) ./... | reviewdog -f=govet -reporter=github-pr-review
```

### With golangci-lint + reviewdog
```yaml
name: reviewdog
on: [pull_request]
jobs:
  golangci-lint:
    name: runner / golangci-lint
    runs-on: ubuntu-latest
    steps:
      - name: Check out code into the Go module directory
        uses: actions/checkout@v1
      - uses: actions/setup-go@v2
        with:
          go-version: '1.15'
      - name: Fetch custom linter
        run: |
          curl -o linters.tar.gz -L https://github.com/sapphi-red/possibleerrors/releases/download/v0.1.16/possibleerrors_golang-ci-lint-plugin_0.1.16_linux_64bit.tar.gz
          tar -xzf linters.tar.gz
      - name: Install golangci-lint with CGO_ENABLED=1
        run: |
          CGO_ENABLED=1 go install github.com/golangci/golangci-lint/cmd/golangci-lint
      - uses: reviewdog/action-setup@v1
      - name: Run reviewdog
        env:
          REVIEWDOG_GITHUB_API_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          golangci-lint run --out-format line-number | reviewdog -f=golangci-lint -reporter=github-pr-review
```
