package page

import (
	"github.com/a-h/templ"
	"github.com/ditwrd/wed/internal/web/component/toast"
)

func ToastError(msg string) templ.Component {
	return toast.Toast(toast.Props{
		Title:         "Error",
		Description:   msg,
		Variant:       toast.VariantError,
		Dismissible:   true,
		Duration:      5000,
		Position:      toast.PositionTopCenter,
		ShowIndicator: true,
		Icon:          true,
	})
}

func ToastSuccess(msg string) templ.Component {
	return toast.Toast(toast.Props{
		Title:         "Success",
		Description:   msg,
		Variant:       toast.VariantSuccess,
		Dismissible:   true,
		Duration:      5000,
		Position:      toast.PositionTopCenter,
		ShowIndicator: true,
		Icon:          true,
	})
}
