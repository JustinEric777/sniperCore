package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filename := "/Users/iknow/Works/Go/Project/sniperCore/resources/11.jpeg"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	filePath := "/Users/iknow/Works/Go/Project/sniperCore/resources/1111.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(string(data))
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}
