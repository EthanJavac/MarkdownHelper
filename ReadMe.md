### ReadMe

解决了`本地markdown文件依赖文件夹内部图片，生成预览版html只能在本地文件夹内浏览而不能拿到其他地方正确展示图片`的问题  

  

可以被处理的markdown文件需要完全只依赖文件夹内的图片，以`![abc](./defg.jpg)` 或者 `![](cde.png)`
或者`![./abc/def/defg.jpg]`等形式引用图片，那么可以根据相对路径把图片用base64替换  

使用要求 `path`输入markdown所在文件夹与go程序运行所在文件夹的相对路径，因此如果是Windows环境，markdown文件夹的根目录
和go程序运行目录需要是在同一个磁盘。  

__能够被正确处理的markdown文件只能放在待处理文件夹的根目录内，
不能放在子文件夹内。
md文件里面引用的图片也只能是放在待处理文件夹根目录或者根目录的子目录下(允许嵌套)__  

```
## 允许的目录结构, testmd作为root文件夹, 合法的md文件只能在root文件夹下
## md内部引用本地图片都在root目录下或者root的子目录下,且通过相对路径引用。
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

## 其中test.md引用了root目录及其子目录的各个图片, 
## test_base64.md就是程序生成的解引用后的文件(通过base64图片编码方式)
```


>如果markdown文件夹是`C://User/123/DeskTop/MarkDownFiles`    
而go程序运行在`C://User/456`  
那么输入`path='../123/DeskTop/MarkDownFiles'`  

```go

/*
https://www.cnblogs.com/gyyyl/p/13606214.html
当main包中有多个go文件时;

package main:
main.go
aa.go
bb.go

此时main包中包含了三个go文件：main.go,aa.go,bb.go，
其中main.go文件中有main函数（必须有main函数，但是main函数不一定必须在main.go文件中）

此时有两种编译或者运行方式

复制代码
#列出所有的文件名
go run main.go aa.go bb.go
go build main.go aa.go bb.go

#使用*.go代替所有文件
go run *.go  //发现不可行 可能windows下不可行
go build *.go //发现不可行 可能windows下不可行

*/
```

todo:  
1. 现在的实现: 一个md文件内所有图片的小名称不可重复 `![](abc.jpg)`和`![](./1/22/abc.jpg)`现在是不允许出现的。
后续需要解决这个问题。  
2. 以后要做base64编码的图片统一放在md结尾处，这样文件更整洁。现在的做法是在插入图片出直接替换为base64编码。
