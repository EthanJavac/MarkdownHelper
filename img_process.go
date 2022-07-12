package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//填充imageMap key=image小名称 value=image的路径+文件名
//abc.jpg -> ../c/a/abc.jpg
//递归读取markdown根文件夹下所有文件(排除dir和markdown类型),key记录shortName,value记录相对路径+shortName
func cacheImagePath(pathname string) error {
	pathname = strings.TrimSuffix(pathname, `/`)
	ds, err := os.ReadDir(pathname)
	if err != nil {
		log.Printf("os_read_dir path=%v, err=%v\n", pathname, err)
		return err
	}
	for _, v := range ds {
		if v.IsDir() {
			name := v.Name()
			err = cacheImagePath(pathname + "/" + name)
			if err != nil {
				return err
			}
		} else {
			log.Printf("path=%v,v= %v\n", pathname, v.Name())
			//abc.png
			if strings.Contains(v.Name(), ".") && !strings.HasSuffix(v.Name(), ".md") {

				imageMap[v.Name()] = strings.TrimSuffix(pathname, `/`) + "/" + v.Name()
			}
		}
	}
	return nil
}

//图片路径转base64
//输入就是 ../abc/def/a.jpg  abc/edf/g.jpg  ./a/b/c.png
//输出就是 data:image/?;base64,wer32r432... base64编码的图片
func readImageToBase64(imageFullPath string) (string, error) {
	//jpg png jpeg
	sls := strings.Split(imageFullPath, `.`)
	class := sls[len(sls)-1]
	log.Printf("图片类型是%v\n", class)

	ff, err := os.Open(imageFullPath)
	if err != nil {
		log.Printf("打开图片%v失败%v\n", imageFullPath, err)
		return "", err
	}
	defer ff.Close()

	data, err := ioutil.ReadAll(ff)
	if err != nil {
		log.Printf("io读取失败%v\n", err)
		return "", err
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	//data:image/png;base64,
	//base64图片的前缀

	prefix0 := `data:image/?;base64,`
	prefix1 := strings.ReplaceAll(prefix0, `?`, class)

	base64Data = prefix1 + base64Data
	return base64Data, nil
}
