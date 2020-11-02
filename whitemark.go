package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/freetype"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "simhei.ttf", "filename of the ttf font")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
)

func main() {
	pwd, _ := os.Getwd()
	//所有需要加水印的图片请放在srcimg文件夹下
	dir_list, e := ioutil.ReadDir(pwd + "/srcimg")
	if e != nil {
		fmt.Println("read dir error")
		return
	}
	for i, v := range dir_list {
		fmt.Println(i, "=", v.Name())
		tempfile := v.Name()
		addWhitemark(pwd+"/srcimg/"+v.Name(), tempfile[0:len(tempfile)-len(filepath.Ext(tempfile))])
	}

}

//加水印函数需要两个参数，文件路径以及文件名称。
func addWhitemark(imgpath string, imgname string) {
	imgorgin, _ := os.Open(imgpath)
	img, _ := jpeg.Decode(imgorgin)
	defer imgorgin.Close()

	flag.Parse()
	//字体为黑体，字体需要提前下载
	fontBytes, err := ioutil.ReadFile(*fontfile)
	checkError(err)
	f, err := freetype.ParseFont(fontBytes)
	checkError(err)

	fg, bg := image.White, image.Transparent
	//获取原始图片尺寸，新建水印图片
	rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	maxxlat := rgba.Bounds().Dx()
	maxylat := rgba.Bounds().Dy()
	//设置水印字体大小，此处设置为图片高度的三十二分之一
	size := float64(maxylat / 32)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// 渲染水印文字
	//whitemark.txt保存了需要添加到图片上的水印文字
	fi, err := os.Open("whitemark.txt")
	checkError(err)
	defer fi.Close()
	//设置水印位置，位置为距离图片右下角，按比例缩放
	pt := freetype.Pt(maxxlat-int(maxxlat/8), maxylat-int(maxylat/8)+int(c.PointToFixed(size)>>6))
	//按行读取数据，逐行渲染
	br := bufio.NewReader(fi)
	for {
		a, _, ch := br.ReadLine()
		_, err = c.DrawString(string(a), pt)
		checkError(err)
		pt.Y += c.PointToFixed(size * *spacing)
		if ch == io.EOF {
			break
		}
	}

	// 保存水印图片
	outFile, err := os.Create("out.png")
	checkError(err)
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	checkError(err)

	err = b.Flush()
	checkError(err)
	fmt.Println("Wrote out.png OK.")
	//读取水印图片
	wmb, _ := os.Open("out.png")
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()
	//把水印图片盖在原始图片上
	offset := image.Pt(0, 0)
	bou := img.Bounds()
	m := image.NewNRGBA(bou)

	draw.Draw(m, bou, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	//保存新的图片
	imgw, _ := os.Create(imgname + "_new.jpg")
	jpeg.Encode(imgw, m, &jpeg.Options{100})

	defer imgw.Close()

	fmt.Printf("水印添加结束,请查看%s_new.jpg图片...\n", imgname)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
