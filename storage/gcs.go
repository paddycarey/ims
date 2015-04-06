package storage

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcs "google.golang.org/api/storage/v1"
)

type GCSFile struct {
	*bytes.Reader
	name    string
	modtime time.Time
}

func (gf *GCSFile) Close() error {
	return nil
}

func (gf *GCSFile) MimeType() string {
	return getMimeTypeFromFilename(gf.name)
}

func (gf *GCSFile) ModTime() time.Time {
	return gf.modtime
}

type GCSFileSystem struct {
	bucket  string
	service *gcs.Service
	client  *http.Client
}

func NewGCSFileSystem(uri, credentials string) (*GCSFileSystem, error) {

	data, err := ioutil.ReadFile(credentials)
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

	bucket, err := parseBucketFromURI(uri)
	if err != nil {
		return nil, err
	}

	gcss := &GCSFileSystem{bucket, service, oauth2Client}
	return gcss, nil
}

func (g *GCSFileSystem) Open(name string) (File, error) {

	logrus.WithFields(logrus.Fields{
		"bucket": g.bucket,
		"file":   name,
	}).Info("Fetching metadata from GCS")

	name = strings.TrimLeft(name, "/")
	res, err := g.service.Objects.Get(g.bucket, name).Do()
	if err != nil {
		return nil, err
	}

	logrus.WithField("file", name).Info("Fetching file from GCS")
	resp, err := g.client.Get(res.MediaLink)
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

	gcsf := &GCSFile{bytes.NewReader(b), name, ts}
	return gcsf, nil
}

func parseBucketFromURI(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}
