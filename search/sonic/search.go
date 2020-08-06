package sonic

import (
	"github.com/expectedsh/go-sonic/sonic"
	. "search/common"
	"search/entity"
	"search/utils"
	"strings"
)

func Search(key string) []entity.Post {

	search, err := sonic.NewSearch(Conf.SonicHost, 1491, Conf.SonicPassword)
	if err != nil {
		panic(err)
	}

	defer search.Quit()

	results, _ := search.Query(Conf.SonicCollection, Conf.SonicBucket, key, Conf.SonicQueryLimit, Conf.SonicQueryOffset)

	if len(results) != 0 && results[0] != "" {

		var list []entity.Post

		for _, item := range results {
			post := entity.Post{}
			ds, _ := utils.DESDecryptString(item)
			split := strings.Split(ds, "|")
			post.Path = strings.Replace(split[0], ".md", "", 1)
			post.Name = split[1]
			list = append(list, post)
		}

		return list
	} else {
		return nil
	}

}
