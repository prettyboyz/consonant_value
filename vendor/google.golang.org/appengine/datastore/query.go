// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package datastore

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"google.golang.org/appengine/internal"
	pb "google.golang.org/appengine/internal/datastore"
)

type operator int

const (
	lessThan operator = iota
	lessEq
	equal
	greaterEq
	greaterThan
)

var operatorToProto = map[operator]*pb.Query_Filter_Operator{
	lessThan:    pb.Query_Filter_LESS_THAN.Enum(),
	lessEq:      pb.Query_Filter_LESS_THAN_OR_EQUAL.Enum(),
	equal:       pb.Query_Filter_EQUAL.Enum(),
	greaterEq:   pb.Query_Filter_GREATER_THAN_OR_EQUAL.Enum(),
	greaterThan: pb.Query_Filter_GREATER_THAN.Enum(),
}

// filter is a conditional filter on query results.
type filter struct {
	FieldName string
	Op        operator
	Value     interface{}
}

type sortDirection int

const (
	ascending sortDirection = iota
	descending
)

var sortDirectionToProto = map[sortDirection]*pb.Query_Order_Direction{
	ascending:  pb.Query_Order_ASCENDING.Enum(),
	descending: pb.Query_Order_DESCENDING.Enum(),
}

// order is a sort order on query results.
type order struct {
	FieldName string
	Direction sortDirection
}

// NewQuery creates a new Query for a specific entity kind.
//
// An empty kind means to return all entities, including entities created and
// managed by other App Engine features, and is called a kindless query.
// Kindless queries cannot include filters or sort orders on property values.
func NewQuery(kind string) *Query {
	return &Query{
		kind:  kind,
		limit: -1,
	}
}

// Query represents a datastore query.
type Query struct {
	kind       string
	ancestor   *Key
	filter     []filter
	order      []order
	projection []string

	distinct bool
	keysOnly bool
	eventual bool
	limit    int32
	offset   int32
	start    *pb.CompiledCursor
	end      *pb.CompiledCursor

	err error
}

func (q *Query) clone() *Query {
	x := *q
	// Copy the contents of the slice-typed fields to a new backing store.
	if len(q.filter) > 0 {
		x.filter = make([]filter, len(q.filter))
		copy(x.filter, q.filter)
	}
	if len(q.order) > 0 {
		x.order = make([]order, len(q.order))
		copy(x.order, q.order)
	}
	return &x
}

// Ancestor returns a derivative query with an ancestor filter.
// The ancestor should not be nil.
func (q *Query) Ancestor(ancestor *Key) *Query {
	q = q.clone()
	if ancestor == nil {
		q.err = errors.New("datastore: nil query ancestor")
		return q
	}
	q.ancestor = ancestor
	return q
}

// EventualConsistency returns a derivative query that returns eventually
// consistent results.
// It only has an effect on ancestor queries.
func (q *Query) EventualConsistency() *Query {
	q = q.clone()
	q.eventual = true
	return q
}

// Filter returns a derivative query with a field-based filter.
// The filterStr argument must be a field name followed by optional space,
// followed by an operator, one of ">", "<", ">=", "<=", or "=".
// Fields are compared against the provided value using the operator.
// Multiple filters are AND'ed together.
func (q *Query) Filter(filterStr string, value interface{}) *Query {
	q = q.clone()
	filterStr = strings.TrimSpace(filterStr)
	if len(filterStr) < 1 {
		q.err = errors.New("datastore: invalid filter: " + filterStr)
		return q
	}
	f := filter{
		FieldName: strings.TrimRight(filterStr, " ><=!"),
		Value:     value,
	}
	switch op := strings.TrimSpace(filterStr[len(f.FieldName):]); op {
	case "<=":
		f.Op = lessEq
	case ">=":
		f.Op = greaterEq
	case "<":
		f.Op = lessThan
	case ">":
		f.Op = greaterThan
	case "=":
		f.Op = equal
	default:
		q.err = fmt.Errorf("datastore: invalid operator %q in filter %q", op, filterStr)
		return q
	}
	q.filter = append(q.filter, f)
	return q
}

// Order returns a derivative query with a field-based sort order. Orders are
// applied in the order they are added. The default order is ascending; to sort
// in descending order prefix the fieldName with a minus sign (-).
func (q *Query) Order(fieldName string) *Query {
	q = q.clone()
	fieldName = strings.TrimSpace(fieldName)
	o := order{
		Direction: ascending,
		FieldName: fieldName,
	}
	if strings.HasPrefix(fieldName, "-") {
		o.Direction = descending
		o.FieldName = strings.TrimSpace(fieldName[1:])
	} else if strings.HasPrefix(fieldName, "+") {
		q.err = fmt.Errorf("datastore: invalid order: %q", fieldName)
		return q
	}
	if len(o.FieldName) == 0 {
		q.err = errors.New("datastore: empty order")
		return q
	}
	q.order = append(q.order, o)
	return q
}

// Project returns a derivative query that yields only the given fields. It
// cannot be used with KeysOnly.
func (q *Query) Project(fieldNames ...string) *Query {
	q = q.clone()
	q.projection = append([]string(nil), fieldNames...)
	return q
}

// Distinct returns a derivative query that yields de-duplicated entities with
// respect to the set of projected fields. It is only used for projection
// queries.
func (q *Query) Distinct() *Query {
	q = q.clone()
	q.distinct = true
	return q
}

// KeysOnly returns a derivative query that yields only keys, not keys and
// entities. It cannot be used with projection queries.
func (q *Query) KeysOnly() *Query {
	q = q.clone()
	q.keysOnly = true
	return q
}

// Limit returns a derivative query that has a limit on the number of results
// returned. A negative value means unlimited.
func (q *Query) Limit(limit int) *Query {
	q = q.clone()
	if limit < math.MinInt32 || limit > math.MaxInt32 {
		q.err = errors.New("datastore: query limit overflow")
		return q
	}
	q.limit = int32(limit)
	return q
}

// Offset returns a derivative query that has an offset of how many keys to
// skip over before returning results. A negative value is invalid.
func (q *Query) Offset(offset int) *Query {
	q = q.clone()
	if offset < 0 {
		q.err = errors.New("datastore: negative query offset")
		return q
	}
	if offset > math.MaxInt32 {
		q.err = errors.New("datastore: query offset overflow")
		return q
	}
	q.offset = int32(offset)
	return q
}

// Start returns a derivative query with the given start point.
func (q *Query) Start(c Cursor) *Query {
	q = q.clone()
	if c.cc == nil {
		q.err = errors.New("datastore: invalid cursor")
		return q
	}
	q.start = c.cc
	return q
}

// End returns a derivative query with the given end point.
func (q *Query) End(c Cursor) *Query {
	q = q.clone()
	if c.cc == nil {
		q.err = errors.New("datastore: invalid cursor")
		return q
	}
	q.end = c.cc
	return q
}

// toProto converts the query to a protocol buffer.
func (q *Query) toProto(dst *pb.Query, appID string) error {
	if len(q.projection) != 0 && q.keysOnly {
		return errors.New("datastore: query cannot both project and be keys-only")
	}
	dst.Reset()
	dst.App = proto.String(appID)
	if q.kind != "" {
		dst.Kind = proto.String(q.kind)
	}
	if q.ancestor != nil {
		dst.Ancestor = keyToProto(appID, q.ancestor)
		if q.eventual {
			dst.Strong = proto.Bool(false)
		}
	}
	if q.projection != nil {
		dst.PropertyName = q.projection
		if q.distinct {
			dst.GroupByPropertyName = q.projection
		}
	}
	if q.keysOnly {
		dst.KeysOnly = proto.Bool(true)
		dst.RequirePerfectPlan = proto.Bool(true)
	}
	for _, qf := range q.filter {
		if qf.FieldName == "" {
			return errors.New("datastore: empty query filter field name")
		}
		p, errStr := valueToProto(appID, qf.FieldName, reflect.ValueOf(qf.Value), false)
		if errStr != "" {
			return errors.New("datastore: bad query filter value type: " + errStr)
		}
		xf := &pb.Query_Filter{
			Op:       operatorToProto[qf.Op],
			Property: []*pb.Property{p},
		}
		if xf.Op == nil {
			return errors.New("datastore: unknown query filter operator")
		}
		dst.Filter = append(dst.Filter, xf)
	}
	for _, qo := range q.order {
		if qo.FieldName == "" {
			return errors.New("datastore: empty query order field name")
		}
		xo := &pb.Query_Order{
			Property:  proto.String(qo.FieldName),
			Direction: sortDirectionToProto[qo.Direction],
		}
		if xo.Direction == nil {
			return errors.New("datastore: unknown query order direction")
		}
		dst.Order = append(dst.Order, xo)
	}
	if q.limit >= 0 {
		dst.Limit = proto.Int32(q.limit)
	}
	if 