curl -X POST http://localhost:8080/prince/generate/test.pdf \
  -F "html=@$PWD/test.html" \
  -F "css=@$PWD/test.css" \
  -H "Content-Type: multipart/form-data"

