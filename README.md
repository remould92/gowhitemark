# gowhitemark
## 一、项目简介
该项目利用golang，实现了一个简单的为图片批量加水印的方法。
## 二、使用方法
```
git clone https://github.com/remould92/gowhitemark.git
cd gowhitemark
```
将需要加水印的图片放在`srcimg`文件夹下，在whitemark.txt中填入水印文字。
```
go run whitemark.go
```
加完水印的文件出现在同一项目目录下.
效果如下所示
![example1](srcimg/beijing.jpeg)

![example1](beijing_new.jpg)