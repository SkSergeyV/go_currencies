package v1

import (
	"currencies/internal/models"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"

	"currencies/pkg/tadapter"
	"currencies/views"
)

type V1Handler struct {
	router     fiber.Router
	repository *V1Repository
}

func NewHandler(router fiber.Router, repo *V1Repository) {
	h := &V1Handler{
		router:     router,
		repository: repo,
	}

	router.Get("/", h.apiInfo)

	router.Get("/exchange_rates/current/json", h.exchangeRatesJson)
	router.Get("/exchange_rates/current/xml", h.exchangeRatesXml)

	router.Get("/exchange_rates/date/json", h.exchangeRatesDateJson)
	router.Get("/exchange_rates/date/xml", h.exchangeRatesDateXml)
}

func (h *V1Handler) apiInfo(c fiber.Ctx) error {
	component := views.Api_info()
	return tadapter.Render(c, component, http.StatusOK)
}

func (h *V1Handler) validateKey(ctx fiber.Ctx) error {
	keyValue := strings.TrimSpace(ctx.Query("key"))
	if keyValue == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Key parameter is required")
	}

	_, err := h.repository.GetActiveByValue(ctx.Context(), keyValue)
	if errors.Is(err, models.ErrKeyNotFound) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired key")
	}
	if err != nil {
		log.Printf("validateKey unexpected error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return nil
}

func (h *V1Handler) exchangeRatesJson(c fiber.Ctx) error {

	isDemo, ok := c.Locals("isDemo").(bool)
	if !ok {
		// ошибка: middleware не отработал
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	if !isDemo {
		//Проверка ключа
		if err := h.validateKey(c); err != nil {
			return err
		}
	}

	currency := c.Query("currencies", "all")
	if currency == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "currencies parameter is required",
		})
	}

	exportData, err := h.repository.GetExportData(currency)
	if err != nil {
		log.Printf("GetExportData error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	exportDates, err := h.repository.GetUpdateExportDates()
	if err != nil {
		log.Printf("GetUpdateExportDates error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	response := models.ExportResponse{
		Dates:      exportDates,
		Currencies: exportData,
	}

	return c.JSON(response)
}

func (h *V1Handler) exchangeRatesXml(c fiber.Ctx) error {

	isDemo, ok := c.Locals("isDemo").(bool)
	if !ok {
		// ошибка: middleware не отработал
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	if !isDemo {
		//Проверка ключа
		if err := h.validateKey(c); err != nil {
			return err
		}
	}

	currency := c.Query("currencies", "all")
	if currency == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	exportData, err := h.repository.GetExportData(currency)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	exportDates, err := h.repository.GetUpdateExportDates()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	currencies := make([]models.ExportCurrencies, len(exportData))
	for i, e := range exportData {
		currencies[i] = models.ExportCurrencies{
			CurrencyId:        e.CurrencyId,
			CurrencyName:      e.CurrencyName,
			CharCode:          e.CharCode,
			NumCode:           e.NumCode,
			Nominal:           e.Nominal,
			ExchangeValue:     e.ExchangeValue,
			PrevExchangeValue: e.PrevExchangeValue,
			Difference:        e.Difference,
		}
	}

	response := models.ExportResponse{
		Dates:      exportDates,
		Currencies: currencies,
	}

	return c.XML(response)
}

func (h *V1Handler) exchangeRatesDateJson(c fiber.Ctx) error {

	isDemo, ok := c.Locals("isDemo").(bool)
	if !ok {
		// ошибка: middleware не отработал
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	if !isDemo {
		//Проверка ключа
		if err := h.validateKey(c); err != nil {
			return err
		}
	}

	date := c.Query("date")
	if date == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	exportData, err := h.repository.GetExportDataDates(date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect date format", //err.Error(),
		})
	}
	exportDates, err := h.repository.GetUpdateExportDatesPlug(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	response := models.ExportResponse{
		Dates:      exportDates,
		Currencies: exportData,
	}

	return c.JSON(response)
}

func (h *V1Handler) exchangeRatesDateXml(c fiber.Ctx) error {

	isDemo, ok := c.Locals("isDemo").(bool)
	if !ok {
		// ошибка: middleware не отработал
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	if !isDemo {
		//Проверка ключа
		if err := h.validateKey(c); err != nil {
			return err
		}
	}

	date := c.Query("date")
	if date == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	exportData, err := h.repository.GetExportDataDates(date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect date format",
		})
	}
	exportDates, err := h.repository.GetUpdateExportDatesPlug(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	currencies := make([]models.ExportCurrencies, len(exportData))
	for i, e := range exportData {
		currencies[i] = models.ExportCurrencies{
			CurrencyId:        e.CurrencyId,
			CurrencyName:      e.CurrencyName,
			CharCode:          e.CharCode,
			NumCode:           e.NumCode,
			Nominal:           e.Nominal,
			ExchangeValue:     e.ExchangeValue,
			PrevExchangeValue: e.PrevExchangeValue,
			Difference:        e.Difference,
		}
	}

	response := models.ExportResponse{
		Dates:      exportDates,
		Currencies: currencies,
	}

	return c.XML(response)

}
