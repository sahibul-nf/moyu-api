package handler

import (
	"fmt"
	"moyu/campaign"
	"moyu/helper"
	"moyu/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

// api/v1/campaigns (Get Campaigns by userID or all)
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.campaignService.FindCampaigns(userID)
	if err != nil {
		campaignsFormatter := campaign.FormatCampaigns(campaigns)

		response := helper.APIResponse("Failed to get campaigns", "error", http.StatusBadRequest, campaignsFormatter)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	campaignsFormatter := campaign.FormatCampaigns(campaigns)

	response := helper.APIResponse("Successfuly get list of campaigns", "success", http.StatusOK, campaignsFormatter)
	c.JSON(http.StatusOK, response)
}

// api/v1/campaigns/{1} (Get Campaign by ID or Campaign Detail)
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	// handler : mapping id yg di url ke struct input => kirim service, call formatter
	// service : tangkap input dari struct input, panggil repo
	// repository : get campaign by id

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign detail", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaigns, err := h.campaignService.FindCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign detail", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatDetailCampaign(campaigns)

	response := helper.APIResponse("Successfuly to get campaign detail", "success", http.StatusOK, formatter)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateNewCampaign(c *gin.Context) {
	// get input user dan mapping ke struct input
	// cek or get current user via jwt/handler
	// service passing struct input dan panggil repo + generete struct
	// repo save to db
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorFormatter := helper.ErrorValidationFormat(err)

		errorMessage := gin.H{"errors": errorFormatter}

		response := helper.APIResponse("Create new campaign is failed", "error", http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		// errorFormatter := helper.ErrorValidationFormat(err)
		// errorMessage := gin.H{"errors": errorFormatter}

		response := helper.APIResponse("Create new campaign is failed", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaign(newCampaign)

	response := helper.APIResponse("Successfuly to create new campaign detail", "success", http.StatusOK, formatter)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	// user masukkan input
	// handler
	// mapping dari input ke input struct
	// input dari user dan juga input dari Uri (passing ke service)
	// service (find campaign by ID, tangkap parameter)
	// repository update data campaign

	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errorFormatter := helper.ErrorValidationFormat(err)

		errorMessage := gin.H{"errors": errorFormatter}

		response := helper.APIResponse("Failed to update campaign", "error", http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", "error", http.StatusUnprocessableEntity, nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Successfuly to updated campaign", "success", http.StatusOK, updatedCampaign)
	c.JSON(http.StatusOK, response)
}

// handler
// 1. tangkap input dan passing ke struct input
// 2. save campaign image ke folder on server
// service
// 1. call 2 point on repo
// repository
// 1. create or save data image to on table campaign-images
// 2. ubah is_primary true ke false

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.SaveCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errorFormatter := helper.ErrorValidationFormat(err)

		errorMessage := gin.H{"errors": errorFormatter}

		response := helper.APIResponse("Failed to upload campaign image", "error", http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse("Failed to upload campaign image", "error", http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// simpan avatar di folder "images/"
	path := fmt.Sprintf("campaign-images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse("Failed to upload campaign image", "error", http.StatusBadRequest, data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.campaignService.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}

		response := helper.APIResponse("Failed to upload campaign image", "error", http.StatusBadGateway, data)

		c.JSON(http.StatusBadGateway, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image successfuly uploaded", "success", http.StatusOK, data)

	c.JSON(http.StatusOK, response)
}
