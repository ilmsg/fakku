package main

import (
	"archive/zip"
	"io/fs"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {

	var target string = "./data"

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		files := getFiles(target)

		return ctx.Render("index", fiber.Map{
			"Title":    "Hello, World!",
			"ZipFiles": files,
		})
	})

	app.Get("/files/:zipIndex", func(ctx *fiber.Ctx) error {
		files := getFiles(target)

		zipIndex, _ := strconv.ParseInt(ctx.Params("zipIndex"), 0, 32)
		filename := getFile(files, int(zipIndex))

		zf := readZipFile(target + "/" + filename)
		defer zf.Close()

		return ctx.Render("files", fiber.Map{
			"Title":        "Hello, World!",
			"ZipIndex":     zipIndex,
			"NextZipIndex": zipIndex + 1,
			"ZipFiles":     files,
			"ImageFiles":   zf.File,
		})
	})

	app.Get("/files/:zipIndex/:imageIndex", func(ctx *fiber.Ctx) error {
		files := getFiles(target)

		index, _ := strconv.ParseInt(ctx.Params("zipIndex"), 0, 32)
		filename := getFile(files, int(index))

		zf := readZipFile(target + "/" + filename)
		defer zf.Close()

		imageIndex, _ := strconv.ParseInt(ctx.Params("imageIndex"), 0, 32)
		content := readContent(zf.File[imageIndex])

		ctx.Write(content)
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}

func getFiles(target string) []fs.FileInfo {
	files, err := ioutil.ReadDir(target)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return files
}

func getFile(files []fs.FileInfo, index int) string {
	return files[index].Name()
}

func readZipFile(target string) *zip.ReadCloser {
	zf, err := zip.OpenReader(target)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return zf
}

func readContent(file *zip.File) []byte {
	fc, err := file.Open()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return content
}
