// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/product-api/client-sdk/models"
)

// ListAllReader is a Reader for the ListAll structure.
type ListAllReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAllReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAllOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListAllOK creates a ListAllOK with default headers values
func NewListAllOK() *ListAllOK {
	return &ListAllOK{}
}

/*ListAllOK handles this case with default header values.

A list of products
*/
type ListAllOK struct {
	Payload []*models.Product
}

func (o *ListAllOK) Error() string {
	return fmt.Sprintf("[GET /products][%d] listAllOK  %+v", 200, o.Payload)
}

func (o *ListAllOK) GetPayload() []*models.Product {
	return o.Payload
}

func (o *ListAllOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}