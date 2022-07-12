package img_process

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//markdown图片要在此文件夹内
//图片是同级关系？ 子级关系？
//返回一个集合 里面是可能的图片->path
var ImageMap map[string]struct{}

/*

testmd
|-- 111111.png
|-- 222.jpg
|-- dir00
|   `-- test_c.jpg
|-- dir01
|   `-- dir02
|       |-- 333.bmp
|       `-- photo.jpg
|-- test_base64.md
`-- test.md


*/
//处理md_root目录下的所有非md类型 非目录类型的文件，
//把其相对root的关系存起来
//就是存成  ../testmd/dir01/dir02/photo.jpg的集合
func CacheImagePath(pathname string) error {
	pathname = strings.TrimSuffix(pathname, `/`)
	log.Printf("]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]pathname=%v\n", pathname)
	ds, err := os.ReadDir(pathname)
	if err != nil {
		log.Printf("os_read_dir path=%v, err=%v\n", pathname, err)
		return err
	}
	for _, v := range ds {
		log.Printf("]]]]]]]]]]]]]ds-v=%v\n", v.Name())
		if v.IsDir() {
			name := v.Name()
			err = CacheImagePath(pathname + "/" + name)
			if err != nil {
				return err
			}
		} else {
			log.Printf("path=%v,v= %v\n", pathname, v.Name())
			//abc.png
			if strings.Contains(v.Name(), ".") && !strings.HasSuffix(v.Name(), ".md") {

				ImageMap[strings.TrimSuffix(pathname, `/`)+"/"+v.Name()] = struct{}{}
			}
		}
	}
	return nil
}

//处理结果 ImageMap : map[111111.png:{} 222.jpg:{} c/test_c.jpg:{} dir01/dir02/333.bmp:{} dir01/dir02/photo.jpg:{}]
func CacheImagePathAndProcess(pathname string) error {
	err := CacheImagePath(pathname)
	if err != nil {
		return err
	}
	log.Printf("map1==================map1\n")
	log.Printf("%v\n", ImageMap)
	log.Printf("map1==================map1\n")

	//遍历map 拿掉头部信息， 这样集合里面都是相对root文件夹的文件路径
	if len(ImageMap) == 0 {
		log.Printf("找不到目录下的图片信息\n")
		return fmt.Errorf("找不到目录下的图片信息")
	}

	m2 := make(map[string]struct{})
	for k := range ImageMap {
		k1 := strings.TrimPrefix(k, pathname)
		//下行代码估计是没用的 先写着吧
		k1 = strings.TrimPrefix(k1, `./`)
		m2[k1] = struct{}{}
	}

	//m2 替代原来map
	ImageMap = m2
	log.Printf("map2==================map2\n")
	log.Printf("%v\n", ImageMap)
	log.Printf("map2==================map2\n")
	return nil
}

//图片路径转base64
//输入就是 ../abc/def/a.jpg  abc/edf/g.jpg  ./a/b/c.png
//输出就是 data:image/?;base64,wer32r432... base64编码的图片
func ReadImageToBase64(imageFullPath string) (string, error) {
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
