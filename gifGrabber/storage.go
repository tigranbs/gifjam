package gifGrabber

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"gifjam/config"
)

var (
	mongoURL = os.Getenv("GIFJAM_MONGO")
	session  *mgo.Session
	db       *mgo.Database
	storage  *mgo.GridFS
)

func initDB() {
	if len(config.GlobalConfig.Mongo) > 0 {
		mongoURL = config.GlobalConfig.Mongo
	}

	// making database on package init
	var err error
	session, err = mgo.Dial(mongoURL)
	if err != nil {
		log.Println("Unable to connect to MongoDB database as a storage backend at url[", mongoURL, "] -> ", err.Error())
		os.Exit(1)
	}

	session.SetMode(mgo.Monotonic, true)

	// Pinging every 1 second to keep connection alive
	go func() {
		for {
			time.Sleep(time.Second * 1)
			err = session.Ping()
			if err != nil {
				log.Println("Error while trying to ping to MongoDB, restting connection -> ", err.Error())
				os.Exit(1)
			}
		}
	}()

	db = session.DB("gifjam")
	storage = db.GridFS("fs")
}

func SaveItem(item *FeedItem) (bool, error) {
	found, err := storage.Find(bson.M{"filename": item.Link}).Count()
	if err != nil {
		return false, err
	}

	// if we already have this image, just returning
	if found > 0 {
		return false, nil
	}

	retry_count := 0

	// if we don't have it downloading it from url
	for {
		if retry_count > 20 {
			log.Println("Retry count existed for url -> ", item.Link)
			return false, nil
		}

		res, err := http.Get(strings.Replace(strings.Replace(item.Link, " ", "", -1), "\t", "", -1))
		if err != nil {
			log.Println("Unable to get Gif Image from url ", item.Link, " to save it to DB -> ", err.Error())
			retry_count++
			time.Sleep(time.Second * 1)
			continue
		}

		if !strings.Contains(res.Header.Get("Content-Type"), "image/gif") {
			return false, nil
		}

		file, err := storage.Create(item.Link)
		if err != nil {
			return false, err
		}

		file.SetMeta(item)
		_, err = io.Copy(file, res.Body)
		if err != nil {
			return false, err
		}

		file.Close()
		res.Body.Close()

		break
	}

	return true, nil
}

// Getting FileID's for downloading them from client
func GetItems(offset, limit int) (ids []string, err error) {
	files := []bson.M{}
	err = storage.Find(bson.M{"metadata.visible": true}).Skip(offset).Limit(limit).All(&files)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		ids = append(ids, file["_id"].(bson.ObjectId).Hex())
	}

	return ids, nil
}

func GetFileIO(id string) (int64, io.ReadCloser, error) {
	file, err := storage.OpenId(bson.ObjectIdHex(id))
	if err != nil {
		return 0, nil, err
	}

	return file.Size(), file, nil
}

func SetVisibility(id string, is_visible bool) error {
	return db.C("fs.files").Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"metadata.visible": is_visible}})
}
