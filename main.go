package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

//markdown的图片检测
//todo 检测markdown的图片引用，包括提取图片名，有待改进检测的正则表达式以及提取图片名策略
const md_image = `!\[(.*)\]\((.*)\)`

var mdImagePattern = regexp.MustCompile(md_image)

//markdown图片要在此文件夹内
//图片是同级关系？ 子级关系？
//返回一个集合 里面是可能的图片->path
var imageMap map[string]string

//只要path填对
func main() {
	path := `../testmd`
	where := false
	processMDRootPath(path, where)
}

//输入markdown的根目录文件夹相对go运行程序的相对路径
//比如 ../channel/ 或者 ../a/b/mdFiles
//会在里面递归查找类似abc.md的文件，然后处理后生成abc_base64.md的新文件
//输入参数where:
//当where=true 新md文件放在原来md的相同位置那里
//当where=false 新md文件放在此go程序的运行位置
func processMDRootPath(markDownRootPath string, where bool) {
	imageMap = make(map[string]string)

	path := strings.TrimSuffix(markDownRootPath, "/") + "/"

	//预处理图片 存imageMap
	err := cacheImagePath(path)
	if err != nil {
		log.Printf("预处理图片出错%v\n", err)
		return
	}

	if len(imageMap) == 0 {
		log.Printf("文件夹内没有非markdown类型的文件了\n")
		return
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Printf("err occur read dir %v\n", err)
		return
	}

	for _, dir := range dirs {
		dirName := dir.Name()
		//log.Printf("cur dir name is %v, \n", dirName)
		if dir.IsDir() {
			//log.Printf("这个dir %v 还是一个目录 \n", dirName)
			continue
		}
		if !strings.Contains(dirName, `.md`) {
			//log.Printf("这个dir %v 不是一个md文件 \n", dirName)
			continue
		}
		if strings.Contains(dirName, `base64`) {
			//log.Printf("这个dir %v 已经是转换后的base64_md \n", dirName)
			continue
		}

		//只剩待转换的md文件!
		log.Printf("======path:[%v]====file:[%v] \n", path, dirName)
		generateNewMDFile(path, dirName, where)
	}
}

//把md的图片改成base64的
//where=true就是放在md那里, where=false就是放在go程序运行的地方
func generateNewMDFile(path, name string, where bool) error {
	finalName := path + name
	finalName2 := strings.TrimSuffix(finalName, `.md`) + "_base64.md"

	if !where {
		slice := strings.Split(finalName2, `/`)
		finalName2 = `./` + slice[len(slice)-1]
	}
	f2, err := os.Create(finalName2)
	if err != nil {
		log.Printf("创建base64_md文件出错%v\n", err)
		return err
	}
	//_ = f2

	fileHandle, err := os.OpenFile(finalName, os.O_RDONLY, 0666)
	if err != nil {
		log.Printf("读取文件出错%v\n", err)
		return err
	}

	defer fileHandle.Close()

	reader := bufio.NewReader(fileHandle)
	var str string

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("读取出错 line=%v err=%v \n", string(line), err)
		}

		//line
		str = string(line)
		if mdImagePattern.MatchString(str) {
			if len(str) <= 6 { //![](x)
				continue
			}
			if !strings.Contains(str, "!") || !strings.Contains(str, "(") || !strings.Contains(str, ")") || !strings.Contains(str, "[") || !strings.Contains(str, "]") {
				continue
			}
			var builder strings.Builder
			log.Printf("=================================\n")
			log.Printf("str match %v", str)
			//![asd](dfs.png)
			// abcd![12](333.jpg)dfweefrw![](./dfsf.jpeg)
			//读取括号 ( )
			//揪出括号对内的字符串聚拢
			var innerWrite = false
			var innerBuilder strings.Builder
			var inners []string
			//40和41分别是()
			for i := 0; i < len(str); i++ {
				v := str[i]
				if v == 40 {
					//( 可以开始内部写了
					builder.WriteByte(v)
					innerWrite = true
					innerBuilder.Reset()
				} else if v == 41 {
					//) 结束内部写
					builder.WriteByte(v)
					innerWrite = false
					if len(innerBuilder.String()) != 0 {
						inners = append(inners, innerBuilder.String())
					}
				} else {
					if innerWrite {
						innerBuilder.WriteByte(v)
					} else {
						builder.WriteByte(v)
					}
				}
			}

			//处理之后 	abcd![12](333.jpg)dfweefrw![](./dfsf.jpeg)
			//就会变成  abcd![12]()dfweefrw![]() inners切片里面放了`333.jpg` `./dfsf.jpeg`

			//()()()

			bs := builder.String()

			for _, inner := range inners {
				for k1, v1 := range imageMap {
					//v1可能是 ../abc/edf/tom.png
					//k1 必定是 tom.png
					//inner 可能是 tom.png 可能是./tom.png
					if strings.Contains(inner, k1) {
						b64Str, err := readImageToBase64(v1)
						if err != nil {
							log.Printf("k1=%v,v=1%v,图片处理出错%v\n", k1, v1, err)
						}
						bs = strings.Replace(bs, `()`, `(`+b64Str+`)`, 1)
					}
				}
			}

			f2.WriteString(bs)
			f2.WriteString("\n")

			//for image, imagePath := range imageMap {
			//	if strings.Contains(str, image) {
			//		imageFinalPath := strings.TrimSuffix(imagePath, `/`) + `/`
			//		imageFullName := imageFinalPath + image
			//		//把imageFullName替换成base64
			//		b64Str, err := readImageToBase64(imageFullName)
			//		if err != nil {
			//			log.Printf("图片处理出错%v\n", err)
			//		}
			//		localMap[image] = b64Str
			//	}
			//}

			log.Printf("=================================\n")
		} else {
			f2.WriteString(str)
			f2.WriteString("\n")
		}

	}
	return nil
}
