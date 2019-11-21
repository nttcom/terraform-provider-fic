package fic

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/nttcom/go-fic/fic/eri/v1/ports"
)

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper to allow for logging.
type LogRoundTripper struct {
	Rt      http.RoundTripper
	OsDebug bool
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	defer func() {
		if request.Body != nil {
			request.Body.Close()
		}
	}()

	// for future reference, this is how to access the Transport struct:
	//tlsconfig := lrt.Rt.(*http.Transport).TLSClientConfig

	var err error

	if lrt.OsDebug {
		log.Printf("[DEBUG] FIC Request URL: %s %s", request.Method, request.URL)
		log.Printf("[DEBUG] FIC Request Headers:\n%s", FormatHeaders(request.Header, "\n"))

		if request.Body != nil {
			request.Body, err = lrt.logRequest(request.Body, request.Header.Get("Content-Type"))
			if err != nil {
				return nil, err
			}
		}
	}

	response, err := lrt.Rt.RoundTrip(request)
	if response == nil {
		return nil, err
	}

	if lrt.OsDebug {
		log.Printf("[DEBUG] FIC Response Code: %d", response.StatusCode)
		log.Printf("[DEBUG] FIC Response Headers:\n%s", FormatHeaders(response.Header, "\n"))

		response.Body, err = lrt.logResponse(response.Body, response.Header.Get("Content-Type"))
	}

	return response, err
}

// logRequest will log the HTTP Request details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *LogRoundTripper) logRequest(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	// Handle request contentType
	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := lrt.formatJSON(bs.Bytes())
		log.Printf("[DEBUG] FIC Request Body: %s", debugInfo)
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

// logResponse will log the HTTP Response details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *LogRoundTripper) logResponse(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	if strings.HasPrefix(contentType, "application/json") {
		var bs bytes.Buffer
		defer original.Close()
		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}
		debugInfo := lrt.formatJSON(bs.Bytes())
		if debugInfo != "" {
			log.Printf("[DEBUG] FIC Response Body: %s", debugInfo)
		}
		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}

	log.Printf("[DEBUG] Not logging because FIC response body isn't JSON")
	return original, nil
}

// formatJSON will try to pretty-format a JSON body.
// It will also mask known fields which contain sensitive information.
func (lrt *LogRoundTripper) formatJSON(raw []byte) string {
	var data map[string]interface{}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse FIC JSON: %s", err)
		return string(raw)
	}

	// Mask known password fields
	if v, ok := data["auth"].(map[string]interface{}); ok {
		if v, ok := v["identity"].(map[string]interface{}); ok {
			if v, ok := v["password"].(map[string]interface{}); ok {
				if v, ok := v["user"].(map[string]interface{}); ok {
					v["password"] = "***"
				}
			}
		}
	}

	// Ignore the catalog
	if v, ok := data["token"].(map[string]interface{}); ok {
		if _, ok := v["catalog"]; ok {
			return ""
		}
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("[DEBUG] Unable to re-marshal FIC JSON: %s", err)
		return string(raw)
	}

	return string(pretty)
}

/*
For FIC specific resources definition
*/

// PortCreateOpts represents the attributes used when creating a new port.
type PortCreateOpts struct {
	ports.CreateOpts
}

// ToPortCreateMap casts a CreateOpts struct to a map.
// It overrides ports.ToPortCreateMap to add the ValueSpecs field.
func (opts PortCreateOpts) ToPortCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "port")
}
