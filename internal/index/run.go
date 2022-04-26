package index

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/internal/storage"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Run 索引写入入口
func Run(meta *engine.Meta, conf *conf.Config) {

	e := engine.NewEngine(meta, conf)
	index, err := NewIndexEngine(e, conf)
	if err != nil {
		panic(err)
	}
	defer index.Close()

	addDoc(index)
}

func addDoc(in *Index) {
	docList := readFiles(in.Conf.Source.Files)
	for _, item := range docList {
		doc, err := doc2Struct(item)
		if err != nil {
			log.Errorf("doc2Struct err: %v", err)
			doc = new(storage.Document)
		}
		log.Debugf("doc_id:%v,title:%s", doc.DocID, doc.Title)
		err = in.AddDocument(doc)
		if err != nil {
			log.Errorf("AddDocument err: %v", err)
			break
		}
		// // 达到阈值
		if len(in.PostingsHashBuf) > 0 && (in.BufCount > in.BufSize) {
			in.Flush()
		}
	}
	// 读取结束 写入磁盘
	in.Flush()
}

func doc2Struct(docStr string) (*storage.Document, error) {

	d := strings.Split(docStr, "\t")

	if len(d) < 3 {
		return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	}
	doc := new(storage.Document)

	docID, err := strconv.Atoi(d[0])
	if err != nil {
		return nil, err
	}
	doc.DocID = uint64(docID)
	doc.Title = d[1]
	doc.Body = d[2]
	return doc, nil
}