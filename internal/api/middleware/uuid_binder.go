package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UUIDBinder extends Echo's default binder to properly handle UUID types
type UUIDBinder struct {
	DefaultBinder echo.DefaultBinder
}

// Bind implements the Binder interface
func (b *UUIDBinder) Bind(i interface{}, c echo.Context) error {
	// Use the default binder first
	if err := b.DefaultBinder.Bind(i, c); err != nil {
		return err
	}

	// Process path parameters
	if len(c.ParamNames()) > 0 {
		for _, name := range c.ParamNames() {
			// Check if this parameter is a UUID and needs to be parsed
			param := c.Param(name)
			if len(param) == 36 || len(param) == 32 {
				// Try to parse it as UUID
				_, err := uuid.Parse(param)
				if err == nil {
					// It's a valid UUID, so let the handler process it
					return nil
				}
			}
		}
	}

	return nil
}
