package mgo

import (
	"bytes"
	"sort"

	"gopkg.in/mgo.v2/bson"
)

// Bulk represents an operation that can be prepared with several
// orthogonal changes before being delivered to the server.
//
// MongoDB servers older than version 2.6 do not have proper support for bulk
// operations, so the driver attempts to map its API as much as possible into
// the functionality that works. In particular, in those releases updates and
// removals are sent individually, and inserts are sent in bulk but have
// suboptimal error reporting compared to more recent versions of the server.
// See the documentation of BulkErrorCase for details on that.
//
// Relevant documentation:
//
//   http://blog.mongodb.org/post/84922794768/mongodbs-new-bulk-api
//
type Bulk struct {
	c       *Collection
	opcount int
	actions []bulkAction
	ordered bool
}

type bulkOp int

const (
	bulkInsert bulkOp = iota + 1
	bulkUpdate
	bulkUpdateAll
	bulkRemove
)

type bulkAction struct {
	op   bulkOp
	docs []interface{}
	idxs []int
}

type bulkUpdateOp []interface{}
type bulkDeleteOp []interface{}

// BulkResult holds the results for a bulk operation.
type BulkResult struct {
	Matched  int
	Modified int // Available only for MongoDB 2.6+

	// Be conservative while we understand exactly how to report these
	// results in a useful and convenient way, and also how to emulate
	// them with prior servers.
	private bool
}

// BulkError holds an error returned from running a Bulk operation.
// Individual errors may be obtained and inspected via the Cases method.
type BulkError struct {
	ecases []BulkErrorCase
}

func (e *BulkError) Error() string {
	if len(e.ecases) == 0 {
		return "invalid BulkError instance: no errors"
	}
	if len(e.ecases) == 1 {
		return e.ecases[0].Err.Error()
	}
	msgs := make([]string, 0, len(e.ecases))
	seen := make(map[string]bool)
	for _, ecase := range e.ecases {
		msg := ecase.Err.Error()
		if !seen[msg] {
			seen[msg] = true
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) == 1 {
		return msgs[0]
	}
	var buf bytes.Buffer
	buf.WriteString("multiple errors in bulk operation:\n")
	for _, msg := range msgs {
		buf.WriteString("  - ")
		buf.WriteString(msg)
		buf.WriteByte('\n')
	}
	return buf.String()
}

type bulkErrorCases []BulkErrorCase

func (slice bulkErrorCases) Len() int           { return len(slice) }
func (slice bulkErrorCases) Less(i, j int) bool { return slice[i].Index < slice[j].Index }
func (slice bulkErrorCases) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// BulkErrorCase holds an individual error found while attempting a single change
// within a bulk operation, and the position in which it was enqueued.
//
// MongoDB servers older than version 2.6 do not have proper support for bulk
// operations, so the driver attempts to map its API as much as possible into
// the functionality that works. In particular, only the last error is reported
// for bulk inserts and without any positional information, so the Index
// field is set to -1 in these cases.
type BulkErrorCase struct {
	Index int // Position of operation that failed, or -1 if unknown.
	Err   error
}

// Cases returns all individual errors found while attempting the requested changes.
//
// See the documentation of BulkErrorCase for limitations in older MongoDB releases.
func (e *BulkError) Cases() []BulkErrorCase {
	return e.ecases
}

// Bulk returns a value to prepare the execution of a bulk operation.
func (c *Collection) Bulk() *Bulk {
	return &Bulk{c: c, ordered: true}
}

// Unordered puts the bulk operation in unordered mode.
//
// In unordered mode the indvidual operations may be sent
// out of order, which means latter operations may proceed
// even if prior ones have failed.
func (b *Bulk) Unordered() {
	b.ordered = false
}

func (b *Bulk) action(op bulkOp, opcount int) *bulkAction {
	var action *bulkAction
	if len(b.actions) > 0 && b.actions[len(b.actions)-1].op == op {
		action = &b.actions[len(b.actions)-1]
	} else if !b.ordered {
		for i := range b.actions {
			if b.actions[i]