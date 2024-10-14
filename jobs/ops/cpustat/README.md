#Setup
`go get github.com/Azure/azure-sdk-for-go/sdk/storage/azblob`
`go get github.com/Azure/azure-sdk-for-go/sdk/azidentity`

#Deployment

1. Build `go build -o <output location>`
2. add files to same directory that the binary will be called `.env config.json`
3. schedule with crontab
