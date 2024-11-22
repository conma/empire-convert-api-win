## convert file .wmf, .emf sang .png  
- chỉ hỗ trợ convert 1 file 1 lúc

- Chạy api: windows:  required:  ImageMagick 7.1.1-39 Q16 x64, cài mathtype (wiris) cho windows  
`go run convert-api.go`  
- Request:  
`curl --location 'http://localhost:8080/convert' \  
--form 'file=@"/D:/work/empire/image234.wmf"' \  
--form 'path="2024/11/12/a"'`
- Response:
  xem ở file `convert-client.go`
