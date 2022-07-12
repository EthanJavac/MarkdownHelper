### ReadMe

解决了本地markdown文件依赖文件夹内部图片的问题  

markdown文件需要完全只依赖文件夹内的图片，以`![abc](./defg.jpg)` 或者 `![](cde.png)`
或者`![./abc/def/defg.jpg]`等形式引用图片，那么可以根据相对路径把图片用base64替换  

使用要求 `path`输入markdown所在文件夹与go程序运行所在文件夹的相对路径，因此如果是Windows环境，markdown文件夹的根目录
和go程序运行目录需要是在同一个磁盘。  

>如果markdown文件夹是`C://User/123/DeskTop/MarkDownFiles`    
而go程序运行在`C://User/456`  
那么输入`path='../123/DeskTop/MarkDownFiles'`


