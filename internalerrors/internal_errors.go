package internalerrors

type InternalError struct {
	Description string `json:"description"`
}

func (internalError *InternalError) Error() string {
	return internalError.Description
}

var (
	DatabaseError = &InternalError{Description: "A database error occured."}
	NotFound      = &InternalError{Description: "The specified resource was not found or you do not have sufficient permissions."}
	Undefined     = &InternalError{Description: "Undefined error."}
)

type ViolatedConstraint struct {
	InternalError
}

func (violatedConstraint *ViolatedConstraint) Error() string {
	return violatedConstraint.Description
}

func NewViolatedConstraint(description string) *ViolatedConstraint {
	return &ViolatedConstraint{InternalError: InternalError{Description: description}}
}
