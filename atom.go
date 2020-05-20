package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

//This code is intended to work with the user input being the name of a plain text file with headers indicated by the
//keyword "topic" and paragraphs indicated by the keyword "body". Using the keywords outside their specific purpose will
//lead to unintended html code

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func textSplitter(fileName string) [][]string {
	f, err := os.Open(fileName)
	check(err)
	defer f.Close()
	fileInfo, err := os.Stat(fileName)
	check(err)
	length := fileInfo.Size()
	file := make([]byte, length)
	_, err = f.Read(file)
	text := string(file)
	topicSplit := strings.Split(text, "<topic>")
	topicSplit = topicSplit[1:]
	splitText := make([][]string, len(topicSplit))
	for i, topic := range topicSplit {
		splitText[i] = strings.Split(topic, "<body>")
	}
	return splitText
}

func newHTML(fileName string, splitText [][]string) {
	title := strings.ReplaceAll(fileName, "_", "-")
	header := fmt.Sprintf(`<! DOCTYPE HTML>
	<html>
		<head>
			<meta charset = “UTF-8”>
			<title>%s</title>
		</head>
		<body>
		`, title)
	htmlSlice := make([]string, len(splitText)+3)
	htmlSlice[0] = header
	for i, section := range splitText {
		topic := fmt.Sprintf("<h4>%s</h4> \n <p> \n %s \n </p>\n ", section[0], section[1])
		htmlSlice[i+1] = topic
	}
	footer := fmt.Sprintf("</body> \n </html>")
	htmlSlice[len(splitText)+1] = footer
	htmlCode := strings.Join(htmlSlice, "")
	htmlFileName := fmt.Sprintf("%s.html", fileName)
	htmlFile, err := os.Create(htmlFileName)
	check(err)
	defer htmlFile.Close()
	n, err := htmlFile.WriteString(htmlCode)
	fmt.Printf("%v bytes written. \n", n)
	check(err)
}

func updateFeed(fileName string, splitText [][]string) {
	f, err := os.Open(fileName)
	check(err)
	defer f.Close()
	var file []byte
	_, err = f.Read(file)
	check(err)
	link := fmt.Sprintf("https://wussell.github.io/%v", fileName)
	title := strings.TrimSuffix(fileName, ".html")
	title = strings.ReplaceAll(title, "_", "-")
	t := time.Now()
	formatted := t.Format(time.RFC3339)
	entryStart := fmt.Sprintf(
		` 
		<entry>
		<title>%s</title>
		<link href="%s"/>
		<id>%s</id>
		<updated>%s</updated>
		  <content>
		  <![CDATA[
		  `, title, link, link, formatted)
	entrySlice := make([]string, len(splitText)+2)
	entrySlice[0] = entryStart
	for i, section := range splitText {
		topic := fmt.Sprintf("<p>%s</p> \n <p>%s</p> \n", section[0], section[1])
		entrySlice[i+1] = topic
	}
	entryEnd := fmt.Sprintf(` ]]>
	</content>
	 </entry>
	  
	</feed>`)
	entrySlice[len(splitText)+1] = entryEnd
	entry := strings.Join(entrySlice, "")
	feed, err := os.OpenFile("/Users/moose1/Documents/Programming/Wussell.github.io/atomFeed.xml", os.O_RDWR, os.ModeDevice)
	check(err)
	defer feed.Close()
	ret, err := feed.Seek(-9, 2)
	fmt.Printf("written at offset %v \n", ret)
	check(err)
	_, err = feed.WriteString(entry)
	check(err)
}

func main() {
	//err := os.Chdir("/Users/moose1/Documents/Programming/Wussell.github.io")
	fmt.Println("Enter the name of the file to update the feed with:")
	var fileName string
	fmt.Scanln(&fileName)
	fmt.Println("inputted fileName: ", fileName)
	splitText := textSplitter(fileName)
	newHTML(fileName, splitText)
	updateFeed(fileName, splitText)
}
