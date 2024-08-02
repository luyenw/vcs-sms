package service

import (
	"fmt"
	"os"
	"time"
	"vcs-sms/config/logger"
	config "vcs-sms/config/storage"
	"vcs-sms/model/entity"

	"cloud.google.com/go/storage"

	"github.com/xuri/excelize/v2"
)

type IXLSXService interface {
	ExportXLSX(servers []entity.Server) (string, error)
	ImportXLSX(filePath string) ([][]string, error)
}

type XLSXService struct {
	gcpService *storage.Client
}

func NewXLSXService() *XLSXService {
	return &XLSXService{
		// gcpService: nil,
		gcpService: config.GetGCPStorage(),
	}
}

func (service *XLSXService) ExportXLSX(servers []entity.Server) (string, error) {
	log := logger.NewLogger()
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(fmt.Sprintf("Error closing file: %s", err))
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		log.Error(fmt.Sprintf("Error reading sheet: %s", err))
		return "", err
	}
	f.SetCellValue("Sheet1", "A2", "Name")
	f.SetCellValue("Sheet1", "B2", "IPv4")
	f.SetCellValue("Sheet1", "C2", "Status")
	f.SetCellValue("Sheet1", "D2", "Created Time")
	f.SetCellValue("Sheet1", "E2", "Last Updated")
	for idx, server := range servers {
		f.SetCellValue("Sheet1", "A"+fmt.Sprint(idx+3), server.Name)
		f.SetCellValue("Sheet1", "B"+fmt.Sprint(idx+3), server.IPv4)
		f.SetCellValue("Sheet1", "C"+fmt.Sprint(idx+3), server.Status)
		f.SetCellValue("Sheet1", "D"+fmt.Sprint(idx+3), server.CreatedTime.Format("2006-01-02 15:04:05"))
		f.SetCellValue("Sheet1", "E"+fmt.Sprint(idx+3), server.LastUpdated.Format("2006-01-02 15:04:05"))
	}
	currentTimeF := time.Now().Format("06-01-02_15-04-05")
	os.MkdirAll("./tmp/export", os.ModePerm)
	exportPath := "./tmp/export/" + currentTimeF + ".xlsx"
	f.SetActiveSheet(index)
	if err := f.SaveAs(exportPath); err != nil {
		log.Error(fmt.Sprintf("Error saving file: %s", err))
		return "", err
	}
	exportFN := currentTimeF + ".xlsx"
	exportFile, err := os.Open(exportPath)
	gcpsClient := &ClientUploader{
		cli:        service.gcpService,
		bucketName: "vcs-sms-bucket",
		uploadPath: "export/",
	}
	if err = gcpsClient.UploadFileAndSetMetaData(exportFile, currentTimeF+".xlsx"); err != nil {
		log.Error(fmt.Sprintf("Error uploading file: %s", err))
		return "", err
	}
	exportURL, err := gcpsClient.GetFileURL(exportFN)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting file URL: %s", err))
		return "", err
	}
	return exportURL, nil
}

func (service *XLSXService) ImportXLSX(filePath string) ([][]string, error) {
	log := logger.NewLogger()
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("Error opening file: %s", err))
		return nil, err
	}
	defer func() {
		f.Close()
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Error(fmt.Sprintf("Error getting rows: %s", err))
		return [][]string{}, nil
	}
	return rows, nil
}
