package stripe

import (
	"net/url"
	"strconv"
)

// ISO 3-digit Currency Codes for major currencies (not the full list).
const (
	USD = "usd" // US Dollar ($)
	EUR = "eur" // Euro (€)
	GBP = "gbp" // British Pound Sterling (UK£)
	JPY = "jpy" // Japanese Yen (¥)
	CAD = "cad" // Canadian Dollar (CA$)
	HKD = "hkd" // Hong Kong Dollar (HK$)
	CNY = "cny" // Chinese Yuan (CN¥)
	AUD = "aud" // Australian Dollar (A$)
)

// Charge represents details about a credit card charge in Stripe.
//
// see https://stripe.com/docs/api#charge_object
type Charge struct {
	ID             string
	Description    string
	Amount         int
	Card           *Card
	Currency       string
	Created        UnixTime
	Customer       string
	Invoice        string
	Paid           bool
	Refunded       bool
	AmountRefunded int    `json:"amount_refunded"`
	FailureMessage string `json:"failure_message"`
	Disputed       bool
	Livemode       bool
}

// FeeDetails represents a single fee associated with a Charge.
type FeeDetails struct {
	Amount      int
	Currency    string
	Type        string
	Description string
	Application string
}

// ChargeParams encapsulates options for creating a new Charge.
type ChargeParams struct {
	// A positive integer in cents representing how much to charge the card.
	// The minimum amount is 50 cents.
	Amount int

	// 3-letter ISO code for currency. Currently, only 'usd' is supported.
	Currency string

	// (Optional) Either customer or card is required, but not both The ID of an
	// existing customer that will be charged in this request.
	Customer string

	// (Optional) Credit Card that should be charged.
	Card *CardParams

	// (Optional) Credit Card token that should be charged.
	Token string

	// An arbitrary string which you can attach to a charge object. It is
	// displayed when in the web interface alongside the charge. It's often a
	// good idea to use an email address as a description for tracking later.
	Description string
}

// ChargeClient encapsulates operations for creating, updating, deleting and
// querying charges using the Stripe REST API.
type ChargeClient struct{}

// Creates a new credit card Charge.
//
// see https://stripe.com/docs/api#create_charge
func (c *ChargeClient) Create(params *ChargeParams) (*Charge, error) {
	charge := Charge{}
	values := url.Values{
		"amount":      {strconv.Itoa(params.Amount)},
		"currency":    {params.Currency},
		"description": {params.Description},
	}

	// add optional credit card details, if specified
	if params.Card != nil {
		appendCardParams(values, params.Card)
	} else if len(params.Token) > 0 {
		values.Add("card", params.Token)
	} else {
		// if no credit card is provide we need to specify the customer
		values.Add("customer", params.Customer)
	}

	err := query("POST", "/charges", values, &charge)
	return &charge, err
}

// Retrieves the details of a charge with the given ID.
//
// see https://stripe.com/docs/api#retrieve_charge
func (c *ChargeClient) Retrieve(id string) (*Charge, error) {
	charge := Charge{}
	path := "/charges/" + url.QueryEscape(id)
	err := query("GET", path, nil, &charge)
	return &charge, err
}

// Refunds a charge for the full amount.
//
// see https://stripe.com/docs/api#refund_charge
func (c *ChargeClient) Refund(id string) (*Charge, error) {
	values := url.Values{}
	charge := Charge{}
	path := "/charges/" + url.QueryEscape(id) + "/refund"
	err := query("POST", path, values, &charge)
	return &charge, err
}

// Refunds a charge for the specified amount.
//
// see https://stripe.com/docs/api#refund_charge
func (c *ChargeClient) RefundAmount(id string, amt int) (*Charge, error) {
	values := url.Values{
		"amount": {strconv.Itoa(amt)},
	}
	charge := Charge{}
	path := "/charges/" + url.QueryEscape(id) + "/refund"
	err := query("POST", path, values, &charge)
	return &charge, err
}

// Returns a list of your Charges.
//
// see https://stripe.com/docs/api#list_charges
func (c *ChargeClient) List() ([]*Charge, error) {
	return c.list("", 10, 0)
}

// Returns a list of your Charges with the specified range.
//
// see https://stripe.com/docs/api#list_charges
func (c *ChargeClient) ListN(count int, offset int) ([]*Charge, error) {
	return c.list("", count, offset)
}

// Returns a list of your Charges with the given Customer ID.
//
// see https://stripe.com/docs/api#list_charges
func (c *ChargeClient) CustomerList(id string) ([]*Charge, error) {
	return c.list(id, 10, 0)
}

// Returns a list of your Charges with the given Customer ID and range.
//
// see https://stripe.com/docs/api#list_charges
func (c *ChargeClient) CustomerListN(id string, count int, offset int) ([]*Charge, error) {
	return c.list(id, count, offset)
}

func (c *ChargeClient) list(id string, count int, offset int) ([]*Charge, error) {
	// define a wrapper function for the Charge List, so that we can
	// cleanly parse the JSON
	type listChargesResp struct{ Data []*Charge }
	resp := listChargesResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	// query for customer id, if provided
	if id != "" {
		values.Add("customer", id)
	}

	err := query("GET", "/charges", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
