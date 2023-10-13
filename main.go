package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func main() {
    err := os.MkdirAll("./files/", 0755)
    if err != nil {
        panic(err)
    }

    gin.SetMode(gin.ReleaseMode)

    router := gin.Default()
    router.POST("/upload", UploadFile)
    router.POST("/download", DownloadFile)
    router.GET("/download/*", DownloadFile)
    router.POST("/ebook-convert", EbookConvertFile)
    router.GET("/ebook-convert/version", EbookConvertVersion)
    router.GET("/files", ListFiles)

    router.SetTrustedProxies(nil)

    router.Run()
}

func UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uuid := uuid.New()

    fileName := uuid.String() + filepath.Ext(file.Filename)
    filePath := filepath.Join("./files/", fileName)

    f, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer f.Close()

    out, err := os.Create(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer out.Close()

    _, err = io.Copy(out, f)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "filename": fileName,
        "filesize": file.Size,
    })
}

func ListFiles(c *gin.Context) {
    files, err := ioutil.ReadDir("files")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    jsonFiles := make([]gin.H,0)
    for _, file := range files {
        jsonFiles = append(jsonFiles, gin.H{
            "filename": file.Name(),
            "filesize": file.Size(),
        })
    }

    c.JSON(http.StatusOK, jsonFiles)
}

func DownloadFile(c *gin.Context) {

    var request struct {
        FileName string `json:"filename" binding:"required"`
    }

    if c.Request.Method == "POST" {
        err := c.ShouldBindJSON(&request)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    } else if c.Request.Method == "GET" {
        request.FileName = filepath.Base(c.Request.URL.Path)
    }

    filePath, _ := filepath.Abs("./files/" + filepath.Base(request.FileName))
    fileData, err := os.Open(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer fileData.Close()

    fileHeader := make([]byte, 512)
    _, err = fileData.Read(fileHeader)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    fileContentType := http.DetectContentType(fileHeader)
    fileInfo, err := fileData.Stat()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(request.FileName)))
    c.Header("Content-Type", fileContentType)
    c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    c.File(filePath)
}

func EbookConvertFile(c *gin.Context) {
    var request struct {
        FileName string `json:"filename" binding:"required"`
        Format   string `json:"format" binding:"required,alphanum,max=10"`
        Params   []string `json:"params" binding:""`
        Debug    bool `json:"debug" binding:""`
    }
    err := c.ShouldBindJSON(&request)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uuid := uuid.New()
    fileNameOut := uuid.String() + "." + request.Format
    fullFileNameOut, _ := filepath.Abs("./files/" + fileNameOut)
    fullRequestFileName, _ := filepath.Abs("./files/" + filepath.Base(request.FileName))

    cmdParams := append([]string{ fullRequestFileName, fullFileNameOut }, request.Params...)
    cmdOutErr, err1 := exec.Command("/usr/bin/ebook-convert", cmdParams...).CombinedOutput()
    if err1 != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error(), "debug": string(cmdOutErr)})
        return
    }
    if request.Debug == true {
        c.JSON(http.StatusOK, gin.H{"filename": fileNameOut, "debug": string(cmdOutErr)})
    } else {
        c.JSON(http.StatusOK, gin.H{"filename": fileNameOut})
    }
}

func EbookConvertVersion(c *gin.Context) {

    cmdParams := append([]string{ "--version" })
    cmdOutErr, err1 := exec.Command("/usr/bin/ebook-convert", cmdParams...).CombinedOutput()
    if err1 != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error(), "debug": string(cmdOutErr)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"version": string(cmdOutErr)})

}
