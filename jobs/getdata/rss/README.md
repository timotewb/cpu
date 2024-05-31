
sudo apt-get install gcc-arm*
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -o build/rss

GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabi-gcc go build -o build/rss

If building on machine running just do `go build -o`

### To test
1. Create `rss.json` file in same dir as `main.go`
2. Add below json to file and add your data 
`
{
    "urls":[
        {"url":"<url>","format":<1,2>}       
    ]
}
`
3. Create `all.json` file in the ame dir as `main.go`
4. Add below json to file and add your data (check dirs exist)
`
{
    "staging_dir":"tmp/staging/",
    "loading_dir":"tmp/loading/",
    "sqlite_max_size_mb":10
}
`
6. Run below command to run
`go run . -c .`

Dont for get to cleaup before pushing