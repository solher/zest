package internalerrors

type InternalError struct {
	Description string `json:"description"`
}

func (internalError *InternalError) Error() string {
	return internalError.Description
}

var (
	RessourceNotFound      = &InternalError{Description: "The specified ressource was not found."}
	DatabaseError          = &InternalError{Description: "A database error occured."}
	InsufficentPermissions = &InternalError{Description: "You do not have sufficient permissions."}
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
