package scenarios

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"

    "tests/scenarios/payloads"

    "github.com/form3tech-oss/f1/v2/pkg/f1/testing"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/walletera/payments-types/publicapi"
)

const (
    authUrl     = "http://localhost:3099/api/v1/auth/identity/sessions"
    paymentsUrl = "http://localhost:3099/api/v1/payments"
)

func DinopayOutboundSucceed(t *testing.T) testing.RunFn {
    basicAuthUsername := os.Getenv("BASIC_AUTH_USERNAME")
    if basicAuthUsername == "" {
        t.Errorf("missing BASIC_AUTH_USERNAME env var")
    }
    basicAuthPassword := os.Getenv("BASIC_AUTH_PASSWORD")
    if basicAuthPassword == "" {
        t.Errorf("missing BASIC_AUTH_PASSWORD env var")
    }

    sessionCookie, err := BasicAuthLogin(authUrl, basicAuthUsername, basicAuthPassword)
    if err != nil {
        t.Error(err)
    }

    return func(t *testing.T) {
        createPaymentAndWaitForConfirmation(t, sessionCookie)
    }
}

func createPaymentAndWaitForConfirmation(t *testing.T, sessionCookie *http.Cookie) {
    rawPayment := fmt.Sprintf(payloads.DinopayOutboundSucceedJSON, uuid.New())
    req, err := http.NewRequest(http.MethodPost, paymentsUrl, strings.NewReader(rawPayment))
    if err != nil {
        t.Error(err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.AddCookie(sessionCookie)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }

    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Error(err)
    }

    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("faile creating payment: status code %d - error %s", resp.StatusCode, respBody)
    }

    var payment = &publicapi.Payment{}
    err = payment.UnmarshalJSON(respBody)
    if err != nil {
        t.Error(err)
    }

    t.Logf("Payment created successfully with id %s", payment.ID)

    // Wait for payment to reach confirmed status and propagate to the read-model
    time.Sleep(50 * time.Millisecond)

    assert.Eventually(t, func() bool {
        p := readPayment(t, sessionCookie, payment.ID.String())
        if p != nil {
            if p.Status == "confirmed" {
                return true
            }
            t.Logf("payment id %s did not reached confirmed status yet (status %s), retrying...", p.ID, p.Status)
        }
        return false
    }, 200*time.Millisecond, 50*time.Millisecond, "payment did not reach confirmed status")
}

func readPayment(t *testing.T, sessionCookie *http.Cookie, paymentID string) *publicapi.Payment {
    readPaymentUrl := paymentsUrl + "/" + paymentID

    req, err := http.NewRequest(http.MethodGet, readPaymentUrl, nil)
    if err != nil {
        t.Error(err)
    }

    req.AddCookie(sessionCookie)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }

    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return nil
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Error(err)
    }

    var payment = &publicapi.Payment{}
    err = payment.UnmarshalJSON(respBody)
    if err != nil {
        t.Error(err)
    }

    return payment
}
