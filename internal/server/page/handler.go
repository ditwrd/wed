package page

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/ditwrd/wed/internal/model"
	"github.com/ditwrd/wed/internal/repository"
	"github.com/ditwrd/wed/internal/web"
	"github.com/ditwrd/wed/internal/web/component/toast"
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

func (h *Handler) Home(c echo.Context) error {
	var homePageProps page.HomePageProps
	viper.UnmarshalKey("data", &homePageProps)

	// Fetch latest messages for carousel
	messages, err := h.rsvpRepo.GetLatestMessages(20)
	if err != nil {
		// Log error but continue without messages
		fmt.Printf("Error fetching messages: %v\n", err)
		messages = []model.RSVP{}
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

	// Get donation props from config

	// Create home props with empty name and group for default route
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
	var homePageProps page.HomePageProps
	viper.UnmarshalKey("data", &homePageProps)

	// Get the name parameter (guest name)
	guestName := c.Param("name")

	// Fetch latest messages for carousel
	messages, err := h.rsvpRepo.GetLatestMessages(20)
	if err != nil {
		// Log error but continue without messages
		fmt.Printf("Error fetching messages: %v\n", err)
		messages = []model.RSVP{}
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

	// Create home props with route parameters
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
	var homePageProps page.HomePageProps
	viper.UnmarshalKey("data", &homePageProps)

	// Get the group parameter
	group := c.Param("group")

	// For group route, we don't modify the welcome section
	// The group info will be handled in the RSVP form

	// Fetch latest messages for carousel
	messages, err := h.rsvpRepo.GetLatestMessages(20)
	if err != nil {
		// Log error but continue without messages
		fmt.Printf("Error fetching messages: %v\n", err)
		messages = []model.RSVP{}
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

	// Create home props with route parameters
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
		// Return error toast
		errorToast := toast.Toast(toast.Props{
			Title:       "Error",
			Description: "Name is required",
			Variant:     toast.VariantError,
			Dismissible: true,
			Duration:    5000,
		})
		return Render(c, http.StatusBadRequest, errorToast)
	}

	if attending != "yes" && attending != "no" {
		// Return error toast
		errorToast := toast.Toast(toast.Props{
			Title:       "Error",
			Description: "Please select whether you're attending or not",
			Variant:     toast.VariantError,
			Dismissible: true,
			Duration:    5000,
		})
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
	if err := h.rsvpRepo.Create(rsvp); err != nil {
		// Return error toast
		errorToast := toast.Toast(toast.Props{
			Title:       "Error",
			Description: "Failed to save RSVP. Please try again.",
			Variant:     toast.VariantError,
			Dismissible: true,
			Duration:    5000,
		})
		return Render(c, http.StatusInternalServerError, errorToast)
	}

	// Return success toast
	successToast := toast.Toast(toast.Props{
		Title:       "Success",
		Description: "RSVP submitted successfully!",
		Variant:     toast.VariantSuccess,
		Dismissible: true,
		Duration:    5000,
	})

	// If there's a message, also update the carousel with OOB swap
	if message != "" {
		// Get latest messages including the new one
		messages, err := h.rsvpRepo.GetLatestMessages(20)
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
			w.Write([]byte("\n"))
			w.Write([]byte(`<div id="message-carousel" hx-swap-oob="outerHTML">`))
			if err := page.MessageCarousel(messageProps).Render(ctx, w); err != nil {
				return err
			}
			w.Write([]byte(`</div>`))
			return nil
		})

		return Render(c, http.StatusOK, combinedResponse)
	}

	return Render(c, http.StatusOK, successToast)
}

// AdminDashboard renders the admin dashboard
func (h *Handler) AdminDashboard(c echo.Context) error {
	// Get pagination parameters
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	// Set defaults
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	// Convert to integers
	limitInt := 10
	offsetInt := 0

	fmt.Sscanf(limit, "%d", &limitInt)
	fmt.Sscanf(offset, "%d", &offsetInt)

	// Get RSVP statistics
	stats, err := h.rsvpRepo.GetStats()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve statistics"},
		)
	}

	// Get RSVP list with pagination
	rsvps, err := h.rsvpRepo.GetPaginated(limitInt, offsetInt)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve RSVPs"},
		)
	}

	// Get total count
	count, err := h.rsvpRepo.GetCount()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve RSVP count"},
		)
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
	// Get pagination parameters
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	// Set defaults
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	// Convert to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return err
	}

	// Get RSVPs with pagination
	rsvps, err := h.rsvpRepo.GetPaginated(limitInt, offsetInt)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve RSVPs"},
		)
	}

	// Get total count
	count, err := h.rsvpRepo.GetCount()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve RSVP count"},
		)
	}

	Render(c, http.StatusOK, page.Pagination(page.AdminDashboardProps{
		Count:  count,
		Offset: offsetInt,
		Limit:  limitInt,
		RSVPs:  rsvps,
	}))
	return Render(c, http.StatusOK, page.AdminTable(rsvps))
	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"rsvps":  rsvps,
	// 	"count":  count,
	// 	"limit":  limitInt,
	// 	"offset": offsetInt,
	// })
}

// AdminRSVPStats returns RSVP statistics
func (h *Handler) AdminRSVPStats(c echo.Context) error {
	stats, err := h.rsvpRepo.GetStats()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve statistics"},
		)
	}

	return c.JSON(http.StatusOK, stats)
}
