package utils

import "net/url"

// RequestPath parse request path string.
//
//	Eg Get `/news/Avatar2023-jcer.jpeg`
//			from `https://902-local.s3.us-west-1.amazonaws.com/news/Avatar2023-jcer.jpeg`
//		Or `/storage/tmp/63e85ba1.jpeg`
//			from `https://www.dancefitvn.com/storage/tmp/63e85ba1.jpeg`
func RequestPath(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return "", err
	}

	return u.Path, nil
}

// RequestParam get a specific request parameter string.
//
//	Eg Get `G-Key=e2c32ed0807ce086083281f131&G-Time=20240803130551`
//		from `http://localhost:7789/api/v1/uploads?G-Key=e2c32ed0807ce086083281f131&G-Time=20240803130551`
func RequestParam(urlStr, param string) (string, error) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return "", err
	}

	return u.Query().Get(param), nil
}
