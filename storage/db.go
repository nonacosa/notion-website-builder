package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

type Website struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Secret      string    `json:"secret"`
	PageID      string    `json:"pageId"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	Theme       string    `json:"theme"`
	Domain      string    `json:"domain"`
	CreateTime  time.Time `json:"CreateTime"`
}

type WebsiteItemMeta struct {
	LastUpdate time.Time `json:"lastUpdate"`
}

type WebsiteItem struct {
	WebsiteItemMeta
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	Categories  string      `json:"categories"`
	Tags        string      `json:"tags"`
	Status      string      `json:"status"`
	Date        string      `json:"date"`
	Description string      `json:"description"`
	FrontMatter FrontMatter `json:frontMatter`
}

type FrontMatter struct {
	//Image         interface{}   `yaml:",flow"`
	Title         interface{}   `yaml:",flow"`
	Status        interface{}   `yaml:",flow"`
	Position      interface{}   `yaml:",flow"`
	Categories    []interface{} `yaml:",flow"`
	Tags          []interface{} `yaml:",flow"`
	Keywords      []interface{} `yaml:",flow"`
	CreateAt      interface{}   `yaml:",flow"`
	Author        interface{}   `yaml:",flow"`
	Date          interface{}   `yaml:",flow"`
	Lastmod       interface{}   `yaml:",flow"`
	Description   interface{}   `yaml:",flow"`
	Draft         interface{}   `yaml:",flow"`
	ExpiryDate    interface{}   `yaml:",flow"`
	PublishDate   interface{}   `yaml:",flow"`
	Show_comments interface{}   `yaml:",flow"`
	//Calculate Chinese word count accurately. Default is true
	IsCJKLanguage interface{} `yaml:",flow"`
	Slug          interface{} `yaml:",flow"`
}

func Save(key string, value string) {
	root, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("could not get the default root directory to use for user-specific configuration data: %w", err)
	}

	if !filepath.IsAbs(root) {
		fmt.Errorf("storage root must be an absolute path, got %s", root)
	}

	// Fyne does not allow to customize the root for a storage
	// so we'll use the same
	storageRoot := filepath.Join(root, AppFolderName, settingFolderName)

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(storageRoot))
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	defer db.Close()
}

func Scan(key string, fn func(val []string, db *badger.DB)) {

	root, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("could not get the default root directory to use for user-specific configuration data: %w", err)
	}

	if !filepath.IsAbs(root) {
		fmt.Errorf("storage root must be an absolute path, got %s", root)
	}

	// Fyne does not allow to customize the root for a storage
	// so we'll use the same
	storageRoot := filepath.Join(root, AppFolderName, settingFolderName)

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(storageRoot))

	db.View(func(txn *badger.Txn) error {
		var rst []string
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(key)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()

			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				rst = append(rst, string(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		fn(rst, db)
		return nil
	})
	defer db.Close()
}
