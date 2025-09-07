package security

import (
    "encoding/xml"
    "io"
    "net/http"

    "github.com/labstack/echo/v4"
)

type Request struct {
    Data string `xml:"data"`
}

func XMLHandler(c echo.Context) error {
    var req Request
    decoder := xml.NewDecoder(io.LimitReader(c.Request().Body, 1 << 20))

    decoder.Strict = true
    decoder.Entity = map[string]string{}

    if err := decoder.Decode(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid XML")
    }

    return nil
}
