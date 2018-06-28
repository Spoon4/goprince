curl -X POST http://localhost:8080/generate/test.pdf \
  -F "html=@$PWD/test.html" \
  -F "css=@$PWD/test.css" \
  -H "Content-Type: multipart/form-data"

