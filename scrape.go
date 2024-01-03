// Package newsletter is ----------------
package newsletter

import (
	"bytes"
	"net/http"
)

// GetReferences returns the content of the references
func GetReferences(references []string) ([]string, error) {
	refContent := make([]string, len(references))

	for _, ref := range references {
		resp, err := http.Get(ref)
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		if resp.StatusCode == 200 {
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(resp.Body)
			if err != nil {
				return nil, err
			}
			bodyString := buf.String()
			refContent = append(refContent, bodyString[:100])
		}
	}
	return refContent, nil
}
