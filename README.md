## convert file .wmf, .emf sang .png  
- chỉ hỗ trợ convert 1 file 1 lúc

- Chạy api: windows:  
`go run convert-api.go`  
- Request:  
`curl --location 'http://localhost:8080/convert' \  
--form 'file=@"/D:/work/empire/image234.wmf"' \  
--form 'path="2024/11/12/a"'`
- Response:
  xem ở file `convert-client.go`
