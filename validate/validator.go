package validate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yiwen101/CardWizards/temp"
)

type Validator interface {
	validate(ctx context.Context, c *app.RequestContext) error
}

type validatorImplement struct{}

func NewValidator() Validator {
	return &validatorImplement{}
}

func (v *validatorImplement) validate(ctx context.Context, c *app.RequestContext) error {
	serviceName := c.Param("serviceName")
	methodName := c.Param("methodName")

	if err := validSerivceAndMethod(serviceName, methodName); err != nil {
		c.SetStatusCode(http.StatusBadRequest)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid route, error message is: %s", err))
		return err
	}

	j, err := treatJsonBody(ctx, c)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Internal Server Error in opening the json body, error message is: %s", err))
		return err
	}
	desc, err := temp.DescsManager.GetFunctionDescriptor(serviceName, methodName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Internal Server Error in getting the function descriptor, error message is: %s", err))
		return err
	}

	err = validateBody(desc, j)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid body, error message is: %s", err))
		return err
	}

	return nil
}
