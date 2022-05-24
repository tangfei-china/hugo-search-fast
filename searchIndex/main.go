package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	. "searchIndex/common"
	_ "searchIndex/common"
	"searchIndex/entity"
	"searchIndex/sonic"
	"sort"
	"strings"
	"time"
)

func filterPost(s []entity.Post, filter func(x entity.Post) bool) []entity.Post {
	newS := s[:0]
	for _, x := range s {
		if !filter(x) {
			newS = append(newS, x)
		}
	}
	return newS
}

func main() {

	log.Info("开始扫描文件")

	var posts []entity.Post

	//获取当前目录下的所有文件或目录信息
	filepath.Walk(Conf.PostPath, func(path string, info os.FileInfo, err error) error {

		suffix := strings.HasSuffix(info.Name(), "md")

		if suffix {
			p := new(entity.Post)
			p.Title = path
			p.DateTime = info.ModTime()

			posts = append(posts, *p)
		}

		return nil
	})

	//排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].DateTime.Before(posts[j].DateTime)
	})

	if !Conf.IndexAllInit {
		now := time.Now().Add(-time.Hour * 2)

		//过滤数据，通过修改时间
		newPost := filterPost(posts, func(x entity.Post) bool {
			return x.DateTime.Before(now)
		})

		posts = newPost
	}

	for _, item := range posts {
		log.Infof("名称：%s - 日期：%s \n", item.Title, item.DateTime)
	}

	log.Info("需要更新索引文件数量：", len(posts))

	if len(posts) == 0 {
		log.Warn("没有文件需要更新索引")
	} else {
		sonic.ProcessIndex(posts)
	}

	log.Info("搜索索引处理结束")

}
