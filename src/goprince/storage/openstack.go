package storage

import (
	"fmt"
	"time"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/objectstorage/v1"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
)

const CONFIG_FILE = "config.json"

type OpenstackConfig struct {
	host string
	project string
	username string
	password string
}

type OpenstackClient struct {
	config *OpenstackConfig
	endpoint string
	session *openstack.Session
}

func NewOpenstackClient() PrinceStorage {

	c := new(OpenstackClient)
	c.GetConfig(CONFIG_FILE)
	c.Authenticate(c.config.username, c.config.password)

	hdr, err := objectstorage.GetAccountMeta(c.session, c.endpoint)
	if err != nil {
		panicString := fmt.Sprint("There was an error getting account metadata:", err)
		panic(panicString)
	}
	_ = hdr
	return c
}

// Load configuration from file
func (c *OpenstackClient) GetConfig(filename string) {

	config, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("ReadFile json failed")
	}
	if err = json.Unmarshal(config, &c.config); err != nil {
		panic("Unmarshal json failed")
	}

	log.Printf("config: %+v\n", c.config)
}

// Authenticate with a project name, username, password.
func (c *OpenstackClient) Authenticate(username string, password string) {

	c.config.username = username
	c.config.password = password

	creds := openstack.AuthOpts{
		AuthUrl:     c.config.host,
		ProjectName: c.config.project,
		Username:    c.config.username,
		Password:    c.config.password,
	}
	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		panicString := fmt.Sprint("There was an error authenticating:", err)
		panic(panicString)
	}
	if !auth.GetExpiration().After(time.Now()) {
		panic("There was an error. The auth token has an invalid expiration.")
	}

	// Find the endpoint for object storage.
	c.endpoint, err = auth.GetEndpoint("object-store", "")
	if c.endpoint == "" || err != nil {
		panic("object-store url not found during authentication")
	}

	c.session, err = openstack.NewSession(nil, auth, nil)
	if err != nil {
		panicString := fmt.Sprint("Error crating new Session:", err)
		panic(panicString)
	}
}

// Get metadata of an object
// object: full path of the object
func (c *OpenstackClient) Metadata(object string) (metadata http.Header, err error) {
	return objectstorage.GetObjectMeta(c.session, c.endpoint+"/"+object)
}

// Retrieve an object
// object: full path of the object to retrieve
// Return array of bytes of the file
func (c *OpenstackClient) Read(object string) []byte {

	_, body, err := objectstorage.GetObject(c.session, c.endpoint+"/"+object)
	if err != nil {
		panicString := fmt.Sprint("GetObject Error:", err)
		panic(panicString)
	}
	return body
}

// Write an object on object container
// file: filename to use for saved file
// container: object container where to write file
// fContent: bytes stream to write
func (c *OpenstackClient) Write(filename string, container string, fContent []byte) {
	var err error

	headers := http.Header{}
	//headers.Add("X-Container-Meta-fubar", "false")
	object := container + "/" + filename
	if err := objectstorage.PutObject(c.session, &fContent, c.endpoint+"/"+object, headers); err != nil {
		panic(err)
	}
	objectsJson, err := objectstorage.ListObjects(c.session, 0, "", "", "", "",
		c.endpoint+"/"+container)

	type objectType struct {
		Name, Hash, Content_type, Last_modified string
		Bytes                                   int
	}
	objectsList := []objectType{}

	if err = json.Unmarshal(objectsJson, &objectsList); err != nil {
		panic(err)
	}
	found := false
	for i := 0; i < len(objectsList); i++ {
		if objectsList[i].Name == filename {
			found = true
		}
	}
	if !found {
		panic("created object is missing from the objectsList")
	}
}

// Delete an object on the container
//object: full path of the object to delete
func (c *OpenstackClient) Delete(object string) {
	if err := objectstorage.DeleteObject(c.session, c.endpoint+"/"+object); err != nil {
		panicString := fmt.Sprint("DeleteObject Error:", err)
		panic(panicString)
	}
}
/*
// Get a list of all the containers at the selected endoint.
func (c *OpenstackClient) Containers() []containerType{} {
	containersJson, err := objectstorage.ListContainers(sess, 0, "", url)
	if err != nil {
		panic(err)
	}

	type containerType struct {
		Name         string
		Bytes, Count int
	}
	containersList := []containerType{}

	if err = json.Unmarshal(containersJson, &containersList); err != nil {
		panic(err)
	}

	found := false
	for i := 0; i < len(containersList); i++ {
		if containersList[i].Name == config.Container {
			found = true
		}
	}
	if !found {
		panic("Created container is missing from downloaded containersList")
	}
}
*/
// Set debug mode
func (c *OpenstackClient) Debug(activate bool) {
	openstack.Debug = &activate
}
