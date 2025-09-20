package scenarios

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
)

type SessionCookie struct {
    Name  string
    Value string
}

func BasicAuthLogin(url, username, password string) (*http.Cookie, error) {
    payload := &bytes.Buffer{}

    writer := multipart.NewWriter(payload)
    err := writer.WriteField("email", username)
    if err != nil {
        return nil, err
    }

    err = writer.WriteField("password", password)
    if err != nil {
        return nil, err
    }

    err = writer.Close()
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPost, url, payload)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := http.DefaultClient
    res, err := client.Do(req)
    if err != nil {
        return nil, err

    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        resBody, err := io.ReadAll(res.Body)
        if err != nil {
            return nil, err
        }
        return nil, fmt.Errorf("basic auth login failed: status code %d - error %s", res.StatusCode, resBody)
    }

    setCookieHeader := res.Header.Get("Set-Cookie")
    if setCookieHeader == "" {
        return nil, fmt.Errorf("no cookie header found")
    }

    cookie, err := http.ParseSetCookie(setCookieHeader)
    if err != nil {
        return nil, err
    }

    return cookie, nil
}
