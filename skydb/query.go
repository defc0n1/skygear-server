package skydb

// SortOrder denotes an the order of Records returned from a Query.
type SortOrder int

// A list of SordOrder, their meaning is self descriptive.
const (
	Ascending SortOrder = iota
	Descending
	Asc  = Ascending
	Desc = Descending
)

// Sort specifies the order of a collection of Records returned from a Query.
//
// Record order can be sorted w.r.t. a record field or a value returned
// from a predefined function.
type Sort struct {
	KeyPath string
	Func    Func
	Order   SortOrder
}

// Operator denotes how the result of a predicate is determined from
// its subpredicates or subexpressions.
//go:generate stringer -type=Operator
type Operator int

// A list of Operator.
const (
	And Operator = iota + 1
	Or
	Not
	Equal
	GreaterThan
	LessThan
	GreaterThanOrEqual
	LessThanOrEqual
	NotEqual
	Like
	ILike
)

// Checks whether the Operator is a compound operator, meaning the
// operator combine the results of other subpredicates.
func (op Operator) IsCompound() bool {
	switch op {
	default:
		return false
	case And, Or, Not:
		return true
	}
}

// Checks whether the Operator determines the result of a predicate
// by comparing two subexpressions.
func (op Operator) IsBinary() bool {
	switch op {
	default:
		return false
	case Equal, GreaterThan, LessThan, GreaterThanOrEqual, LessThanOrEqual, NotEqual, Like, ILike:
		return true
	}
}

// Returns the type of an Expression.
type ExpressionType int

// A list of ExpressionTypes.
const (
	Literal ExpressionType = iota
	KeyPath
	Function
)

// An Expression represents value to be compared against.
type Expression struct {
	Type  ExpressionType
	Value interface{}
}

// Predicate is a representation of used in query for filtering records.
type Predicate struct {
	Operator Operator
	Children []interface{}
}

// GetSubPredicates returns Predicate.Children as []Predicate.
//
// This method is only valid when Operator is either And, Or and Not. Caller
// is responsible to check for this preconditions. Otherwise the method
// will panic.
func (p Predicate) GetSubPredicates() (ps []Predicate) {
	for _, childPred := range p.Children {
		ps = append(ps, childPred.(Predicate))
	}
	return
}

// GetExpressions returns Predicate.Children as []Expression.
//
// This method is only valid when Operator is binary operator. Caller
// is responsible to check for this preconditions. Otherwise the method
// will panic.
func (p Predicate) GetExpressions() (ps []Expression) {
	for _, childPred := range p.Children {
		ps = append(ps, childPred.(Expression))
	}
	return
}

// Query specifies the type, predicate and sorting order of Database
// query.
// ReadableBy is a temp solution for ACL before a full predicate implemented.
type Query struct {
	Type         string
	Predicate    *Predicate
	Sorts        []Sort
	ReadableBy   string
	ComputedKeys map[string]Expression
	DesiredKeys  []string
	Limit        uint64
	Offset       uint64
}

// Func is a marker interface to denote a type being a function in skydb.
//
// skydb's function receives zero or more arguments and returns a DataType
// as a result. Result data type is currently omitted in this interface since
// skygear doesn't use it internally yet. In the future it can be utilized to
// provide more extensive type checking at handler level.
type Func interface {
	Args() []interface{}
}

// DistanceFunc represents a function that calculates distance between
// a user supplied location and a Record's field
type DistanceFunc struct {
	Field    string
	Location *Location
}

// Args implements the Func interface
func (f *DistanceFunc) Args() []interface{} {
	return []interface{}{f.Field, f.Location}
}