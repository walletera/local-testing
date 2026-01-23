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
	accountsUrl              = "http://localhost:3099/api/v1/accounts"
	dinopayGatewayWebhookUrl = "http://localhost:8686/webhooks"
)

func DinopayInboundSucceed(t *testing.T) testing.RunFn {
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
		createAccount(t, sessionCookie)
		createDinopayPaymentAndWaitForConfirmation(t, sessionCookie)
	}
}

func createAccount(t *testing.T, sessionCookie *http.Cookie) {
	accountId := uuid.New()
	rawAccount := fmt.Sprintf(payloads.DinopayAccountSucceedJSON, accountId)
	req, err := http.NewRequest(http.MethodPost, accountsUrl, strings.NewReader(rawAccount))
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(sessionCookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		t.Log("account already exists, skipping account creation")
		return
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed creating account: status code %d", resp.StatusCode)
	}

	t.Logf("dinopay account created successfully with id %s", accountId)
}

func createDinopayPaymentAndWaitForConfirmation(t *testing.T, sessionCookie *http.Cookie) {
	externalPaymentId := uuid.New()
	rawPayment := fmt.Sprintf(payloads.DinopayInboundSucceedJSON, externalPaymentId)
	req, err := http.NewRequest(http.MethodPost, dinopayGatewayWebhookUrl, strings.NewReader(rawPayment))
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(sessionCookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed creating dinopay payment: status code %d", resp.StatusCode)
	}

	t.Logf("Dinopay Payment created successfully with external id %s", externalPaymentId)

	// Wait for payment to reach confirmed status and propagate to the read-model
	time.Sleep(50 * time.Millisecond)

	assert.Eventually(t, func() bool {
		p := readPaymentUsingExternalId(t, sessionCookie, externalPaymentId.String())
		if p != nil {
			if p.Status == "confirmed" {
				return true
			}
			t.Logf("payment id %s did not reached confirmed status yet (status %s), retrying...", p.ID, p.Status)
		}
		return false
	}, 200*time.Millisecond, 50*time.Millisecond, "payment did not reach confirmed status")
}

func readPaymentUsingExternalId(t *testing.T, sessionCookie *http.Cookie, externalId string) *publicapi.Payment {
	readPaymentUrl := fmt.Sprintf("%s?externalId=%s", paymentsUrl, externalId)

	req, err := http.NewRequest(http.MethodGet, readPaymentUrl, nil)
	if err != nil {
		t.Log(err.Error())
		return nil
	}

	req.AddCookie(sessionCookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log(err.Error())
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		t.Logf("no payment found for external id %s", externalId)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		t.Logf("unexpected status code: %d", resp.StatusCode)
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Log(err.Error())
		return nil
	}

	var paymentList = &publicapi.ListPaymentsOK{}
	err = paymentList.UnmarshalJSON(respBody)
	if err != nil {
		t.Log(err.Error())
		return nil
	}

	if !paymentList.Total.Set || paymentList.Total.Value == 0 {
		t.Logf("no payments found for external id %s", externalId)
		return nil
	}

	if paymentList.Total.Value > 1 {
		t.Logf("multiple payments found for external id %s", externalId)
		return nil
	}

	return &paymentList.Items[0]
}
