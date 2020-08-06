package sonic

import (
	"bufio"
	"github.com/expectedsh/go-sonic/sonic"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	. "searchIndex/common"
	"searchIndex/entity"
	"searchIndex/utils"
	"strings"
)

func CreateIndex(posts []entity.Post) {

	var list []sonic.IngestBulkRecord

	for _, item := range posts {
		record := initRecord(item.Title)
		if record.Object != "" && record.Text != "" {
			list = append(list, record)
		}
	}

	if len(list) == 0 {
		log.Warn("组装索引数据没有数据")
		return
	}

	for _, item := range list {
		log.Info(item.Object, "-", item.Text)
	}

	log.Info("开始连接Sonic服务")
	//开始插入索引数据
	ingester, err := sonic.NewIngester(Conf.SonicHost, 1491, Conf.SonicPassword)
	if err != nil {
		log.Error(err)
	}

	defer ingester.Quit()

	for _, item := range list {
		ingester.FlushObject(Conf.SonicCollection, Conf.SonicBucket, item.Object)
	}

	log.Info("开始写入索引数据")
	bulks := ingester.BulkPush(Conf.SonicCollection, Conf.SonicBucket, 3, list)
	if len(bulks) > 0 {
		log.Error("写入索引有异常")
		for _, item := range bulks {
			log.Errorf("Object: %s - Error: %s", item.Object, item.Error)
		}
	} else {
		log.Info("成功写入索引：", len(list))
	}
}

//组装索引数据
//Object,Text
func initRecord(strPath string) sonic.IngestBulkRecord {

	var res sonic.IngestBulkRecord

	//提取相对路径，为了网页可以使用
	res.Object = strings.Replace(strPath, Conf.RootPath, "", 1) + "|"

	//读取单个文件，准备提取内部关键字
	file, err := os.Open(strPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	/*
		ScanLines (默认)
		ScanWords
		ScanRunes (遍历UTF-8字符非常有用)
		ScanBytes
	*/

	var postTitle, postDesc string

	//是否有下一行
	for scanner.Scan() {

		//提取标题内容
		title, titleStr := utils.MatchString(`^title: "(.+)"$`, scanner.Text())
		if title {
			postTitle = titleStr
		}

		//提取描述内容
		desc, descStr := utils.MatchString(`^description: "(.+)"$`, scanner.Text())
		if desc {
			postDesc = descStr
		}

		//截止提取位置
		matched, _ := regexp.MatchString(`^author: ""$`, scanner.Text())
		if matched {
			break
		}

	}

	if postTitle == "" || postDesc == "" {
		res.Object = ""
		res.Text = ""
		return res
	}

	res.Object += postTitle
	res.Text = postTitle + " " + postDesc

	//Object加密，内容包含了地址和标题
	encryptString, _ := utils.DESEncryptString(res.Object)
	res.Object = encryptString

	return res

}
