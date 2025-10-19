package frontend

import "net/url"

func AddQueryParams(value string, params map[string]string) (string, error) {
	parsedUrl, err := url.Parse(value)
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	for k, v := range params {
		queryParams.Add(k, v)
	}

	parsedUrl.RawQuery = queryParams.Encode()

	return parsedUrl.String(), nil
}
