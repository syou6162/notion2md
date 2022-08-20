package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjk/notionapi"
)

func getExportPageURL(client *notionapi.Client, pageId string) (string, error) {
	if strings.HasPrefix(pageId, "https://") {
		pageId = notionapi.ExtractNoDashIDFromNotionURL(pageId)
	}
	return client.RequestPageExportURL(pageId, "markdown", false)
}

func readMainMarkdownFileFromZip(exportURL string) (string, error) {
	resp, err := http.Get(exportURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", err
	}

	markdownContent := ""

	for _, zipFile := range zipReader.File {
		// 階層配下はひとまず出力しない
		if filepath.Dir(zipFile.Name) != "." {
			continue
		}

		// 画像などはコピーしない
		if filepath.Ext(zipFile.Name) != ".md" {
			continue
		}

		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return "", err
		}
		markdownContent += string(unzippedFileBytes)
	}

	return markdownContent, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func main() {
	authToken := os.Getenv("NOTION_TOKEN")
	if authToken == "" {
		log.Fatalln("You have to set the env var NOTION_TOKEN.")
	}

	var (
		pageId = flag.String("page_id", "", "notion page id to export")
	)
	flag.Parse()
	client := &notionapi.Client{
		AuthToken: authToken,
	}

	exportURL, err := getExportPageURL(client, *pageId)
	if err != nil {
		log.Fatal(err)
	}

	markdown, err := readMainMarkdownFileFromZip(exportURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
}
