package validate

import (
	"github.com/cloudwego/hertz/pkg/app"
	desc "github.com/yiwen101/CardWizards/common/descriptor"
)

type Validator interface {
	ValidateBody(c *app.RequestContext, servinceName, methodName string) error
}

func NewValidatorFor(serviceName, methodName string) (Validator, error) {
	dm, err := desc.GetDescriptorManager()
	if err != nil {
		return nil, err
	}

	desc, err := dm.GetFunctionDescriptor(serviceName, methodName)
	if err != nil {
		//c.String(http.StatusInternalServerError, fmt.Sprintf("Internal Server Error in getting the function descriptor, error message is: %s", err))
		return nil, err
	}
	return &validatorImplement{functionDescriptor: desc}, nil
}

func (v *validatorImplement) ValidateBody(c *app.RequestContext, serviceName, methodName string) error {
	j, err := treatJsonBody(c)
	if err != nil {
		//c.String(http.StatusInternalServerError, fmt.Sprintf("Internal Server Error in opening the json body, error message is: %s", err))
		return err
	}

	err = validateBody(v.functionDescriptor, j)
	if err != nil {
		//c.String(http.StatusBadRequest, fmt.Sprintf("Invalid body, error message is: %s", err))
		return err
	}

	return nil
}
