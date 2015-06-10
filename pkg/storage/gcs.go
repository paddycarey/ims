package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	gocache "github.com/pmylund/go-cache"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcs "google.golang.org/api/storage/v1"

	"github.com/paddycarey/ims/pkg/config"
)

type cachedGCSFile struct {
	b  []byte
	ts time.Time
}

type GCSFileSystem struct {
	bucket  string
	dir     string
	service *gcs.Service
	client  *http.Client
	cache   *gocache.Cache
}

func NewGCSFileSystem(uri string) (*GCSFileSystem, error) {

	jsonStruct := &struct {
		PrivateKeyId string `json:"private_key_id"`
		PrivateKey   string `json:"private_key"`
		ClientEmail  string `json:"client_email"`
		ClientId     string `json:"client_id"`
		Type         string `json:"type"`
	}{
		PrivateKeyId: config.GetConfigFromEnv("GCS_PRIVATE_KEY_ID", true),
		PrivateKey:   config.GetConfigFromEnv("GCS_PRIVATE_KEY", true),
		ClientEmail:  config.GetConfigFromEnv("GCS_CLIENT_EMAIL", true),
		ClientId:     config.GetConfigFromEnv("GCS_CLIENT_ID", true),
		Type:         "service_account",
	}

	data, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/devstorage.read_only")
	if err != nil {
		return nil, err
	}

	oauth2Client := conf.Client(oauth2.NoContext)
	service, err := gcs.New(oauth2Client)
	if err != nil {
		return nil, err
	}

	bucket, dir, err := parseURI(uri)
	if err != nil {
		return nil, err
	}

	cache := gocache.New(5*time.Minute, 30*time.Second)
	gcss := &GCSFileSystem{bucket, dir, service, oauth2Client, cache}
	return gcss, nil
}

func (g *GCSFileSystem) Open(name string) (File, error) {

	path := strings.TrimLeft(fmt.Sprintf("%s%s", g.dir, name), "/")
	u := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucket, path)

	file, found := g.cache.Get(u)
	if found {
		cgf := file.(cachedGCSFile)
		return &InMemoryFile{bytes.NewReader(cgf.b), getMimeTypeFromFilename(name), cgf.ts}, nil
	}

	logrus.WithFields(logrus.Fields{
		"bucket": g.bucket,
		"file":   name,
	}).Info("Fetching metadata from GCS")

	res, err := g.service.Objects.Get(g.bucket, path).Do()
	if err != nil {
		return nil, err
	}

	logrus.WithField("file", name).Info("Fetching file from GCS")

	resp, err := g.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ts, err := time.Parse(time.RFC3339, res.Updated)
	if err != nil {
		return nil, err
	}

	if len(b) > 0 {
		g.cache.Set(u, cachedGCSFile{b, ts}, gocache.DefaultExpiration)
	}
	gcsf := &InMemoryFile{bytes.NewReader(b), getMimeTypeFromFilename(name), ts}
	return gcsf, nil
}

func parseURI(uri string) (string, string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}

	return u.Host, strings.TrimRight(u.Path, "/"), nil
}
