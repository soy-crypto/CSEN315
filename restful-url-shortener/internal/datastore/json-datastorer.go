package datastore

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"sync"
	"time"
)

// JsonFileStore uses a file as the datastore for a set of links and possibly users
type JsonFileStore struct {
	filename  string
	cache     map[string]Link
	writeLock sync.Mutex
}

// This innocuous check verifies that a pointer to a JsonFileStore implements all the
// methods specified in LinkStorer
var _ LinkStorer = (*JsonFileStore)(nil)

// NewLinkStorer is a constructor to create a FileStorer that is backed by
// a single JSON file
func NewJsonFileStore(filename string) (*JsonFileStore, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err)
	}

	contents, err := io.ReadAll((f))
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %s ", err)
	}

	cache := make(map[string]Link)
	if err := json.Unmarshal(contents, &cache); err != nil {
		return nil, fmt.Errorf("error while reading json file: %s", err)
	}

	j := &JsonFileStore{
		filename: filename,
		cache:    cache,
	}

	return j, nil
}

func (store *JsonFileStore) updateFile() {
	store.writeLock.Lock()
	defer store.writeLock.Unlock()
	bytes, err := json.Marshal(store.cache)

	if err != nil {
		fmt.Println("unable to store cache data because it could not be marshalled")
	}

	err = os.WriteFile(store.filename, bytes, fs.FileMode(os.O_WRONLY))
	if err != nil {
		fmt.Println("error writing file, will try again later")
	}
}

func (store *JsonFileStore) GetLink(id string) (*Link, error) {
	link, exists := store.cache[id]

	// might want to return a NotAuthorized error if we ever introduce
	// the notion of a private link
	if !exists {
		return nil, &NotFoundError{}
	}

	return &link, nil
}

func (store *JsonFileStore) CreateLink(url string, owner string) (*Link, error) {
	hasher := sha1.New()
	hasher.Write([]byte(url + owner + fmt.Sprintf("%d", rand.Int())))
	id := hex.EncodeToString(hasher.Sum(nil)[:3])
	print(id)

	link := Link{
		Id:        id,
		Url:       url,
		Owner:     owner,
		Views:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	store.cache[id] = link
	go store.updateFile()
	return &link, nil
}

// GetUserLinks retrieves all links from the datastore which are owned by the given user
func (store *JsonFileStore) GetUserLinks(user string) []Link {
	return []Link{}
}

// DeleteLink deletes a link from the store's cache if it's present and then updates the file
// in which all inkls are stored
func (store *JsonFileStore) DeleteLink(id string, user string) error {
	return nil
}
