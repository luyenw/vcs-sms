package service

import (
	"fmt"
	"log"
	"os"
	"time"
	"vcs-sms/config/storage"
	"vcs-sms/model/entity"

	"github.com/xuri/excelize/v2"
)

type XLSXService struct {
}

func NewXLSXService() *XLSXService {
	return &XLSXService{}
}

func (service *XLSXService) ExportXLSX(servers []entity.Server) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return "", err
	}
	// f.SetCellValue("Sheet1", "A1", fmt.Sprintf("Sorted by %s %s, Page %d, PageSize %d", queryParam.SortBy, queryParam.Order, queryParam.Page, queryParam.PageSize))
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
	exportPath := "./tmp/export/" + currentTimeF + ".xlsx"
	f.SetActiveSheet(index)
	if err := f.SaveAs(exportPath); err != nil {
		fmt.Println(err)
	}
	exportFN := currentTimeF + ".xlsx"
	exportFile, err := os.Open(exportPath)
	gcpsClient := &ClientUploader{
		cli:        storage.GetGCPStorage(),
		bucketName: "vcs-sms-bucket",
		uploadPath: "export/",
	}
	if err = gcpsClient.UploadFileAndSetMetaData(exportFile, currentTimeF+".xlsx"); err != nil {
		return "", err
	}
	exportURL, err := gcpsClient.GetFileURL(exportFN)
	if err != nil {
		return "", err
	}
	return exportURL, nil
}

func (service *XLSXService) ImportXLSX(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
