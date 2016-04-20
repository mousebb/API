package middleware

import "github.com/curt-labs/API/models/customer"

// BuildDataContext Generates the DataContext property based off the requesting
// API credentials.
func (ctx *APIContext) BuildDataContext(k string, t string, requireSudo bool) error {
	dtx, err := customer.NewContext(ctx.Session, k, t, requireSudo)
	if err != nil {
		return err
	}

	ctx.DataContext = dtx

	return nil
}
