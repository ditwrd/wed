package page

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/ditwrd/wed/internal/model"
	"github.com/ditwrd/wed/internal/repository"
	"github.com/ditwrd/wed/internal/server/httputil"
	"github.com/ditwrd/wed/internal/web"
	"github.com/ditwrd/wed/internal/web/page"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Handler struct {
	rsvpRepo repository.RSVP
}

func NewHandler(rsvpRepo repository.RSVP) *Handler {
	return &Handler{
		rsvpRepo: rsvpRepo,
	}
}

func (h *Handler) Static() http.FileSystem {
	fsys, err := fs.Sub(web.WebStaticFS, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func (h *Handler) loadHomeBase(ctx context.Context) (page.HomePageProps, []page.MessageProps) {
	var homePageProps page.HomePageProps
	_ = viper.UnmarshalKey("data", &homePageProps)
	messages, err := h.rsvpRepo.GetLatestMessages(ctx, 20)
	if err != nil {
		messages = []model.RSVP{}
	}
	var messageProps []page.MessageProps
	for _, rsvp := range messages {
		if rsvp.Message != "" {
			messageProps = append(
				messageProps,
				page.MessageProps{Name: rsvp.Name, Group: rsvp.GroupName, Content: rsvp.Message},
			)
		}
	}
	return homePageProps, messageProps
}

func (h *Handler) Home(c echo.Context) error {
	homePageProps, messageProps := h.loadHomeBase(c.Request().Context())
	homeProps := page.HomeProps{
		PageProps: homePageProps,
		GuestName: "",
		Group:     "",
		Messages:  messageProps,
	}
	return Render(c, http.StatusOK, page.Home(homeProps))
}

// PersonalizedHome handles /p/<people_name> route
func (h *Handler) PersonalizedHome(c echo.Context) error {
	homePageProps, messageProps := h.loadHomeBase(c.Request().Context())
	guestName := strings.ReplaceAll(c.Param("name"), "_", " ")
	homeProps := page.HomeProps{
		PageProps: homePageProps,
		Group:     "",
		GuestName: guestName,
		Messages:  messageProps,
	}
	return Render(c, http.StatusOK, page.Home(homeProps))
}

// GroupHome handles /g/<group_name> route
func (h *Handler) GroupHome(c echo.Context) error {
	homePageProps, messageProps := h.loadHomeBase(c.Request().Context())
	group := strings.ReplaceAll(c.Param("group"), "_", " ")
	homeProps := page.HomeProps{
		PageProps: homePageProps,
		GuestName: "",
		Group:     group,
		Messages:  messageProps,
	}
	return Render(c, http.StatusOK, page.Home(homeProps))
}

// HandleRSVP processes RSVP form submissions
func (h *Handler) HandleRSVP(c echo.Context) error {
	// Parse form data
	name := c.FormValue("name")
	attending := c.FormValue("attending")
	message := c.FormValue("message")
	groupName := c.FormValue("group_name")

	// Validate required fields
	if name == "" {
		errorToast := ToastError("Name is required")
		return Render(c, http.StatusBadRequest, errorToast)
	}

	if attending != "yes" && attending != "no" {
		errorToast := ToastError("Please select whether you're attending or not")
		return Render(c, http.StatusBadRequest, errorToast)
	}

	// Create RSVP entry
	rsvp := &model.RSVP{
		Name:      name,
		Attending: attending == "yes",
		Message:   message,
		GroupName: groupName,
	}

	// Save to database
	if err := h.rsvpRepo.Create(c.Request().Context(), rsvp); err != nil {
		errorToast := ToastError("Failed to save RSVP. Please try again.")
		return Render(c, http.StatusInternalServerError, errorToast)
	}

	// Return success toast
	successToast := ToastSuccess("RSVP submitted successfully!")

	// If there's a message, also update the carousel with OOB swap
	if message != "" {
		// Get latest messages including the new one
		messages, err := h.rsvpRepo.GetLatestMessages(c.Request().Context(), 20)
		if err != nil {
			// Log error but continue with just the toast
			fmt.Printf("Error fetching messages: %v\n", err)
			return Render(c, http.StatusOK, successToast)
		}

		// Convert database RSVPs to MessageProps
		var messageProps []page.MessageProps
		for _, rsvp := range messages {
			if rsvp.Message != "" {
				messageProps = append(messageProps, page.MessageProps{
					Name:    rsvp.Name,
					Group:   rsvp.GroupName,
					Content: rsvp.Message,
				})
			}
		}

		// Create a component that includes both the toast and OOB carousel update
		combinedResponse := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			// Render the success toast
			if err := successToast.Render(ctx, w); err != nil {
				return err
			}

			// Render the OOB carousel update
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
			if _, err := w.Write([]byte(`<div id="message-carousel" hx-swap-oob="outerHTML">`)); err != nil {
				return err
			}
			if err := page.MessageCarousel(messageProps).Render(ctx, w); err != nil {
				return err
			}
			if _, err := w.Write([]byte(`</div>`)); err != nil {
				return err
			}
			return nil
		})

		return Render(c, http.StatusOK, combinedResponse)
	}

	return Render(c, http.StatusOK, successToast)
}

// AdminDashboard renders the admin dashboard
func (h *Handler) AdminDashboard(c echo.Context) error {
	limitInt, offsetInt, err := ParsePagination(c)
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid pagination")
	}

	// Get RSVP statistics
	stats, err := h.rsvpRepo.GetStats(c.Request().Context())
	if err != nil {
		return httputil.RespondError(c, http.StatusInternalServerError, "Failed to retrieve statistics")
	}

	// Get RSVP list with pagination
	rsvps, err := h.rsvpRepo.GetPaginated(c.Request().Context(), limitInt, offsetInt)
	if err != nil {
		return httputil.RespondError(c, http.StatusInternalServerError, "Failed to retrieve RSVPs")
	}

	// Get total count
	count, err := h.rsvpRepo.GetCount(c.Request().Context())
	if err != nil {
		return httputil.RespondError(c, http.StatusInternalServerError, "Failed to retrieve RSVP count")
	}

	// Prepare props for the template
	props := page.AdminDashboardProps{
		Stats:  stats,
		RSVPs:  rsvps,
		Count:  count,
		Limit:  limitInt,
		Offset: offsetInt,
	}

	return Render(c, http.StatusOK, page.AdminDashboard(props))
}

// AdminRSVPList returns paginated RSVP data
func (h *Handler) AdminRSVPList(c echo.Context) error {
	limitInt, offsetInt, err := ParsePagination(c)
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid pagination")
	}

	// Get RSVPs with pagination
	rsvps, err := h.rsvpRepo.GetPaginated(c.Request().Context(), limitInt, offsetInt)
	if err != nil {
		return httputil.RespondError(c, http.StatusInternalServerError, "Failed to retrieve RSVPs")
	}

	// Get total count
	count, err := h.rsvpRepo.GetCount(c.Request().Context())
	if err != nil {
		return httputil.RespondError(c, http.StatusInternalServerError, "Failed to retrieve RSVP count")
	}

	if err := Render(c, http.StatusOK, page.Pagination(page.AdminDashboardProps{
		Count:  count,
		Offset: offsetInt,
		Limit:  limitInt,
		RSVPs:  rsvps,
	})); err != nil {
		return err
	}
	return Render(c, http.StatusOK, page.AdminTable(rsvps))
}

// AdminRSVPStats returns RSVP statistics
func (h *Handler) AdminRSVPStats(c echo.Context) error {
	stats, err := h.rsvpRepo.GetStats(c.Request().Context())
	if err != nil {
		return httputil.RespondError(c, http.StatusInternalServerError, "Failed to retrieve statistics")
	}

	return httputil.RespondOK(c, stats)
}
