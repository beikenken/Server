/*
 * simple blog
 *
 * A Simple Blog
 *
 * API version: 1.0.0
 * Contact: apiteam@swagger.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package t

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//

	sw "github.com/SYSU-SimpleBlog/Server/go"

	"github.com/boltdb/bolt"
)

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func CreateTable() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Article"))
		if b == nil {
			//create table "xx" if not exits
			b, err = tx.CreateBucket([]byte("Article"))
			if err != nil {
				log.Fatal(err)
			}
		}
		if b != nil {
			var article sw.Article
			var tags []sw.Tag
			tags = append(tags, sw.Tag{"CS"})
			tags = append(tags, sw.Tag{"SC"})

			filePath := "./data"
			files, err := ioutil.ReadDir(filePath)
			if err != nil {
				log.Fatal(err)
			}
			for i := 1; i <= len(files); i++ {
				path := filePath + "/" + strconv.Itoa(i)
				fileInfoList, err := ioutil.ReadDir(path)
				var articleName string
				for i := 0; i < len(fileInfoList); i++ {
					if fileInfoList[i].IsDir() == false {
						articleName = fileInfoList[i].Name()
					}
				}
				if err != nil {
					log.Fatal(err)
				}
				content, err := ioutil.ReadFile(path + "/" + articleName)
				if err != nil {
					fmt.Println("获取失败", err)
					return err
				}

				//fmt.Println("文本内容为:", string(content))

				title := articleName[:len(articleName)-3]
				article = sw.Article{int32(i), title, tags, "2019", string(content)}
				v, err := json.Marshal(article)
				//insert rows
				err = b.Put(itob(i), v)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			return errors.New("Table Article doesn't exist")
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func GetArticleById(id int) {
	//connect to database
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//query the article by ID
	var article sw.Article
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Article"))
		if b != nil {
			v := b.Get(itob(id))
			if v == nil {
				fmt.Println(id, " Article Not Exists")
				return errors.New("Article Not Exists")
			} else {
				_ = json.Unmarshal(v, &article)
				return nil
			}
		} else {
			fmt.Println("Article Not Exists")
			return errors.New("Article Not Exists")
		}
	})
	//fmt.Println(article.Content)
}

func GetArticles(p int) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//display 10 articles per page
	IdIndex := (p-1)*10 + 1
	var articles sw.ArticlesResponse
	var article sw.ArticleResponse
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Article"))
		if b != nil {
			c := b.Cursor()
			k, v := c.Seek(itob(IdIndex))
			if k == nil {
				fmt.Println("Page is out of index")
				return errors.New("Page is out of index")
			}
			key := binary.BigEndian.Uint64(k)
			fmt.Print(key)
			if int(key) != IdIndex {
				fmt.Println("Page is out of index")
				return errors.New("Page is out of index")
			}
			count := 0
			var ori_artc sw.Article
			for ; k != nil && count < 10; k, v = c.Next() {
				err = json.Unmarshal(v, &ori_artc)
				if err != nil {
					return err
				}
				article.Id = ori_artc.Id
				article.Name = ori_artc.Name
				articles.Articles = append(articles.Articles, article)
				count = count + 1
			}
			return nil
		} else {
			return errors.New("Article Not Exists")
		}
	})
	for i := 0; i < len(articles.Articles); i++ {
		fmt.Println(articles.Articles[i])
	}
}

func DeleteArticleById(id int) {
	//connect to database
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//delete the article by ID
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Article"))
		if b != nil {
			c := b.Cursor()
			c.Seek(itob(id))
			err := c.Delete()
			if err != nil {
				//fmt.Println("Delete article failed")
				//log.Fatal(err)
				return errors.New("Delete article failed")
			}
		} else {
			//fmt.Println("Article Not Exists")
			return errors.New("Article Not Exists")
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully Delete article ", id)
}

func DBTestArticle() {
	fmt.Println()
	fmt.Println("DBTestArticle")
	CreateTable()
	GetArticleById(1)
	GetArticleById(5)
	GetArticles(1) /*
		DeleteArticleById(5)
		GetArticleById(5)
	*/
	CreateUser()
}
