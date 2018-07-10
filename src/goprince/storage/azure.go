package storage

import (
	"net/http"
)

type AzureClient struct {
	vnetClient string
}

func NewAzureClient() PrinceStorage {
	c := new(AzureClient)
	return c
}

func (c *AzureClient) Authenticate(username string, password string) {

}

func (c *AzureClient) GetConfig(filename string) {

}

func (c *AzureClient) Read(object string) []byte {
	return nil
}

func (c *AzureClient) Write(filename string, container string, fContent []byte) {

}

func (c *AzureClient) Delete(object string) {

}

func (c *AzureClient) Metadata(object string) (metadata http.Header, err error) {
	return nil, nil
}

func (c *AzureClient) Debug(activate bool) {

}
