# go-synology
Synology NAS  API Client in golang

# Example Usage
```
syno, err := synology.Login("https://10.0.0.2:5001, "user", "password")
if err != nil {
   log.Panic(err)
}

shares, err := syno.ListShares()
for _, share := range shares {
   ...
   ... 

```


# Warning
This package is not even at "alpha" quality. Beware, use, test, contribute. 


# Notes:
This is a hobbu project based on API documentation provided by Synology:
https://global.download.synology.com/download/Document/DeveloperGuide/Synology_File_Station_API_Guide.pdf

