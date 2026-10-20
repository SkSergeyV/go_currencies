package home

import (
	"currencies/pkg/tadapter"
	"net/http"

	"currencies/views"

	"github.com/gofiber/fiber/v3"
)

type HomeHandler struct {
	router     fiber.Router
	repository *CurrenciesRepository
}

func NewHandler(router fiber.Router, repository *CurrenciesRepository) {
	h := &HomeHandler{
		router:     router,
		repository: repository,
	}

	h.router.Get("/", h.home)
}

func (h *HomeHandler) home(c fiber.Ctx) error {

	currencies, err := h.repository.getCurrencies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	lastUpdateDate, lastUpdateDateFrom, err := h.repository.GetUpdateDates()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	component := views.Home(currencies, lastUpdateDate, lastUpdateDateFrom)
	return tadapter.Render(c, component, http.StatusOK)

}
