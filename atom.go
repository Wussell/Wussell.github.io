package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func update() {
	fmt.Println("Enter the name of the file to update the feed with:")
	scanner := bufio.NewScanner(os.Stdin)

	input := make([]byte, 64)
	//n, err := os.Stdin.Read(input)
	fmt.Printf("%v bytes read \n", n)
	check(err)
	fileName := string(input)
	fmt.Println("raw fileName: ", fileName)
	fileName = strings.Trim(fileName, "\n ")
	fmt.Println("trimmed fileName: ", fileName)
	link := fmt.Sprintf("https://wussell.github.io/%v", fileName)
	title := strings.TrimSuffix(fileName, ".html")
	title = strings.ReplaceAll(title, "_", "-")
	t := time.Now()
	formatted := t.Format(time.RFC3339)
	entry := fmt.Sprintf(
		` 
		<entry>
		<title>%s</title>
		<link href="%s"/>
		<id>%s</id>
		<updated>%s</updated>
	  </entry>
	  
	</feed>`, title, link, link, formatted)
	feed, err := os.OpenFile("/Users/moose1/Documents/Programming/Wussell.github.io/atomFeed.xml", os.O_RDWR, os.ModeDevice)
	check(err)
	defer feed.Close()
	//feedInfo, err := feed.Stat()
	//length := feedInfo.Size()
	//fmt.Println(length)
	ret, err := feed.Seek(-9, 2)
	fmt.Printf("written at offset %v \n", ret)
	check(err)
	_, err = feed.WriteString(entry)
	check(err)
}

func main() {
	//err := os.Chdir("/Users/moose1/Documents/Programming/Wussell.github.io")
	update()
}
