package stripe

import (
	"net/url"
	"strconv"
)

// Plan Intervals
const (
	IntervalMonth = "month"
	IntervalYear  = "year"
)

// Plan holds details about pricing information for different products and
// feature levels on your site. For example, you might have a $10/month plan
// for basic features and a different $20/month plan for premium features.
//
// see https://stripe.com/docs/api#plan_object
type Plan struct {
	ID              string
	Name            string
	Amount          int
	Interval        string
	Currency        string
	TrialPeriodDays int `json:"trial_period_days"`
	Livemode        bool
}

// PlanClient encapsulates operations for creating, updating, deleting and
// querying plans using the Stripe REST API.
type PlanClient struct{}

// PlanParams encapsulates options for creating a new Plan.
type PlanParams struct {
	// Unique string of your choice that will be used to identify this plan
	// when subscribing a customer.
	ID string

	// A positive integer in cents (or 0 for a free plan) representing how much
	// to charge (on a recurring basis)
	Amount int

	// 3-letter ISO code for currency. Currently, only 'usd' is supported.
	Currency string

	// Specifies billing frequency. Either month or year.
	Interval string

	// Name of the plan, to be displayed on invoices and in the web interface.
	Name string

	// (Optional) Specifies a trial period in (an integer number of) days. If
	// you include a trial period, the customer won't be billed for the first
	// time until the trial period ends. If the customer cancels before the
	// trial period is over, she'll never be billed at all.
	TrialPeriodDays int
}

// Creates a new Plan.
//
// see https://stripe.com/docs/api#create_plan
func (c *PlanClient) Create(params *PlanParams) (*Plan, error) {
	plan := Plan{}
	values := url.Values{
		"id":       {params.ID},
		"name":     {params.Name},
		"amount":   {strconv.Itoa(params.Amount)},
		"interval": {params.Interval},
		"currency": {params.Currency},
	}

	// trial_period_days is optional, add if specified
	if params.TrialPeriodDays != 0 {
		values.Add("trial_period_days", strconv.Itoa(params.TrialPeriodDays))
	}

	err := query("POST", "/plans", values, &plan)
	return &plan, err
}

// Retrieves the plan with the given ID.
//
// see https://stripe.com/docs/api#retrieve_plan
func (c *PlanClient) Retrieve(id string) (*Plan, error) {
	plan := Plan{}
	path := "/plans/" + url.QueryEscape(id)
	err := query("GET", path, nil, &plan)
	return &plan, err
}

// Updates the name of a plan. Other plan details (price, interval, etc.) are,
// by design, not editable.
//
// see https://stripe.com/docs/api#update_plan
func (c *PlanClient) Update(id string, newName string) (*Plan, error) {
	values := url.Values{"name": {newName}}
	plan := Plan{}
	path := "/plans/" + url.QueryEscape(id)
	err := query("POST", path, values, &plan)
	return &plan, err
}

// Deletes a plan with the given ID.
//
// see https://stripe.com/docs/api#delete_plan
func (c *PlanClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/plans/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of your Plans.
//
// see https://stripe.com/docs/api#list_Plans
func (c *PlanClient) List() ([]*Plan, error) {
	return c.ListN(10, 0)
}

// Returns a list of your Plans at the specified range.
//
// see https://stripe.com/docs/api#list_Plans
func (c *PlanClient) ListN(count int, offset int) ([]*Plan, error) {
	// define a wrapper function for the Plan List, so that we can
	// cleanly parse the JSON
	type listPlanResp struct{ Data []*Plan }
	resp := listPlanResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	err := query("GET", "/plans", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
