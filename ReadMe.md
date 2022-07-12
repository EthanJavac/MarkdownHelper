### ReadMe

解决了本地markdown文件依赖文件夹内部图片的问题  

markdown文件需要完全只依赖文件夹内的图片，以`![abc](./defg.jpg)` 或者 `![](cde.png)`
或者`![./abc/def/defg.jpg]`等形式引用图片，那么可以根据相对路径把图片用base64替换  

使用要求 `path`输入markdown所在文件夹与go程序运行所在文件夹的相对路径，因此如果是Windows环境，markdown文件夹的根目录
和go程序运行目录需要是在同一个磁盘。  

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
