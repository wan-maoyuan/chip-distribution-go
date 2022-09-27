package service

import (
	"chip-distribution-go/pkg/database"
	"chip-distribution-go/pkg/entity"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ChipService struct {
	engine *database.DataBaseEngine
}

func NewChipService() (*ChipService, error) {
	var service = &ChipService{}

	engine, err := database.NewDataBaseEngine()
	if err != nil {
		return service, err
	}

	service.engine = engine

	return service, nil
}

func (service *ChipService) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (service *ChipService) GetAllStockInfos(c *gin.Context) {
	company := c.Query("company")

	infos, err := service.engine.QueryStockInfoByCompanyName(company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.GetAllStockInfosResponse{
		StockInfos: infos,
	})
}

func (service *ChipService) UploadExcel(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	formFiles := form.File["file"]
	if len(formFiles) != 1 {
		c.JSON(http.StatusBadRequest, entity.ErrResponse{
			Message: "only allow one excel file",
		})
		return
	}

	fileHeader := formFiles[0]
	filePath := ExcelFileDir + string(os.PathSeparator) + fileHeader.Filename
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrResponse{
			Message: err.Error(),
		})
		return
	}
	defer os.Remove(filePath)

	// 读取excel中的数据
	stockInfoList, companyMap, err := readExcelData(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	go service.handleExcelData(stockInfoList, companyMap)

	c.JSON(http.StatusOK, "OK")
}

func (service *ChipService) GetAllCompanies(c *gin.Context) {
	companies, err := service.engine.QueryAllCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.GetAllCompaniesResponse{
		Companies: companies,
	})
}

func (service *ChipService) Close() {
	service.engine.Close()
}
