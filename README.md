# TaglessContent
Simple Go app which retrieves Capi V2 content without xml tags

## Endpoints

### /striptags
#### PUT
Id: valid CAPI V2 content uuid
Token: Authorization token

Example:
`curl -X PUT --data "{\"Id\":\"ff682b3e-fea7-4f10-8526-d8ad9ab7848b\",\"Token\":\"MyToken\"}" localhost:8080/striptags`  