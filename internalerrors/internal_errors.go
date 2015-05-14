package internalerrors

type InternalError struct {
	Description string `json:"description"`
}

func (internalError *InternalError) Error() string {
	return internalError.Description
}

var (
	RessourceNotFound  = &InternalError{Description: "The specified ressource was not found."}
	ViolatedConstraint = &InternalError{Description: "A database constraint was violated."}
	DatabaseError      = &InternalError{Description: "A database error occured."}
)
