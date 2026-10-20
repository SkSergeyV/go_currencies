package importer

import (
	"currencies/config"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"crypto/subtle"
)

type ImporterHandler struct {
	router     fiber.Router
	config     *config.ImportConfig
	repository *ImportRepository
}

func NewHandler(router fiber.Router, cfg *config.ImportConfig, repo *ImportRepository) {
	h := &ImporterHandler{
		router:     router,
		config:     cfg,
		repository: repo,
	}

	h.router.Get("/import", h.importer)
}

func (h *ImporterHandler) importer(c fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	b1 := []byte(key)
	b2 := []byte(h.config.ImportKey)

	if len(b1) != len(b2) || subtle.ConstantTimeCompare(b1, b2) != 1 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	importUrl := h.config.ImportUrl

	resp, err := http.Get(importUrl)
	if err != nil {
		log.Printf("Ошибка при запросе к URL: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d, expected %d", resp.StatusCode, http.StatusOK)
		return err
	}

	var jsonimport JsonImport
	err = json.NewDecoder(resp.Body).Decode(&jsonimport)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		return err
	}

	ctx := c.Context()
	err = h.repository.InsertExchangeRates(ctx, jsonimport)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.repository.UpdateExchangeUpd(ctx, jsonimport.Date)
	if err != nil {
		log.Printf("failed to update exchange_upd: %v", err)
		return err
	}

	return c.JSON(fiber.Map{"status": "ok", "currencies": len(jsonimport.Valute)})

}
