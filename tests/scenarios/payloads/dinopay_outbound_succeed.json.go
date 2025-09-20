package payloads

const DinopayOutboundSucceedJSON = `
{
  "id": "%s",
  "amount": 100,
  "currency": "USD",
  "gateway": "dinopay",
  "debtor": {
    "institutionName": "dinopay",
    "institutionId": "dinopay",
    "currency": "USD",
    "accountDetails": {
      "accountType": "dinopay",
      "accountHolder": "Richard Roe",
      "accountNumber": "1200079635"
    }
  },
  "beneficiary": {
    "institutionName": "dinopay",
    "institutionId": "dinopay",
    "currency": "USD",
    "accountDetails": {
      "accountType": "dinopay",
      "accountHolder": "Richard Roe",
      "accountNumber": "1200079635"
    }
  }
}
`
