// Code generated by protoc-gen-go.
// source: google.golang.org/appengine/internal/datastore/datastore_v3.proto
// DO NOT EDIT!

/*
Package datastore is a generated protocol buffer package.

It is generated from these files:
	google.golang.org/appengine/internal/datastore/datastore_v3.proto

It has these top-level messages:
	Action
	PropertyValue
	Property
	Path
	Reference
	User
	EntityProto
	CompositeProperty
	Index
	CompositeIndex
	IndexPostfix
	IndexPosition
	Snapshot
	InternalHeader
	Transaction
	Query
	CompiledQuery
	CompiledCursor
	Cursor
	Error
	Cost
	GetRequest
	GetResponse
	PutRequest
	PutResponse
	TouchRequest
	TouchResponse
	DeleteRequest
	DeleteResponse
	NextRequest
	QueryResult
	AllocateIdsRequest
	AllocateIdsResponse
	CompositeIndices
	AddActionsRequest
	AddActionsResponse
	BeginTransactionRequest
	CommitResponse
*/
package datastore

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Property_Meaning int32

const (
	Property_NO_MEANING       Property_Meaning = 0
	Property_BLOB             Property_Meaning = 14
	Property_TEXT             Property_Meaning = 15
	Property_BYTESTRING       Property_Meaning = 16
	Property_ATOM_CATEGORY    Property_Meaning = 1
	Property_ATOM_LINK        Property_Meaning = 2
	Property_ATOM_TITLE       Property_Meaning = 3
	Property_ATOM_CONTENT     Property_Meaning = 4
	Property_ATOM_SUMMARY     Property_Meaning = 5
	Property_ATOM_AUTHOR      Property_Meaning = 6
	Property_GD_WHEN          Property_Meaning = 7
	Property_GD_EMAIL         Property_Meaning = 8
	Property_GEORSS_POINT     Property_Meaning = 9
	Property_GD_IM            Property_Meaning = 10
	Property_GD_PHONENUMBER   Property_Meaning = 11
	Property_GD_POSTALADDRESS Property_Meaning = 12
	Property_GD_RATING        Property_Meaning = 13
	Property_BLOBKEY          Property_Meaning = 17
	Property_ENTITY_PROTO     Property_Meaning = 19
	Property_INDEX_VALUE      Property_Meaning = 18
)

var Property_Meaning_name = map[int32]string{
	0:  "NO_MEANING",
	14: "BLOB",
	15: "TEXT",
	16: "BYTESTRING",
	1:  "ATOM_CATEGORY",
	2:  "ATOM_LINK",
	3:  "ATOM_TITLE",
	4:  "ATOM_CONTENT",
	5:  "ATOM_SUMMARY",
	6:  "ATOM_AUTHOR",
	7:  "GD_WHEN",
	8:  "GD_EMAIL",
	9:  "GEORSS_POINT",
	10: "GD_IM",
	11: "GD_PHONENUMBER",
	12: "GD_POSTALADDRESS",
	13: "GD_RATING",
	17: "BLOBKEY",
	19: "ENTITY_PROTO",
	18: "INDEX_VALUE",
}
var Property_Meaning_value = map[string]int32{
	"NO_MEANING":       0,
	"BLOB":             14,
	"TEXT":             15,
	"BYTESTRING":       16,
	"ATOM_CATEGORY":    1,
	"ATOM_LINK":        2,
	"ATOM_TITLE":       3,
	"ATOM_CONTENT":     4,
	"ATOM_SUMMARY":     5,
	"ATOM_AUTHOR":      6,
	"GD_WHEN":          7,
	"GD_EMAIL":         8,
	"GEORSS_POINT":     9,
	"GD_IM":            10,
	"GD_PHONENUMBER":   11,
	"GD_POSTALADDRESS": 12,
	"GD_RATING":        13,
	"BLOBKEY":          17,
	"ENTITY_PROTO":     19,
	"INDEX_VALUE":      18,
}

func (x Property_Meaning) Enum() *Property_Meaning {
	p := new(Property_Meaning)
	*p = x
	return p
}
func (x Property_Meaning) String() string {
	return proto.EnumName(Property_Meaning_name, int32(x))
}
func (x *Property_Meaning) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Property_Meaning_value, data, "Property_Meaning")
	if err != nil {
		return err
	}
	*x = Property_Meaning(value)
	return nil
}

type Property_FtsTokenizationOption int32

const (
	Property_HTML Property_FtsTokenizationOption = 1
	Property_ATOM Property_FtsTokenizationOption = 2
)

var Property_FtsTokenizationOption_name = map[int32]string{
	1: "HTML",
	2: "ATOM",
}
var Property_FtsTokenizationOption_value = map[string]int32{
	"HTML": 1,
	"ATOM": 2,
}

func (x Property_FtsTokenizationOption) Enum() *Property_FtsTokenizationOption {
	p := new(Property_FtsTokenizationOption)
	*p = x
	return p
}
func (x Property_FtsTokenizationOption) String() string {
	return proto.EnumName(Property_FtsTokenizationOption_name, int32(x))
}
func (x *Property_FtsTokenizationOption) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Property_FtsTokenizationOption_value, data, "Property_FtsTokenizationOption")
	if err != nil {
		return err
	}
	*x = Property_FtsTokenizationOption(value)
	return nil
}

type EntityProto_Kind int32

const (
	EntityProto_GD_CONTACT EntityProto_Kind = 1
	EntityProto_GD_EVENT   EntityProto_Kind = 2
	EntityProto_GD_MESSAGE EntityProto_Kind = 3
)

var EntityProto_Kind_name = map[int32]string{
	1: "GD_CONTACT",
	2: "GD_EVENT",
	3: "GD_MESSAGE",
}
var EntityProto_Kind_value = map[string]int32{
	"GD_CONTACT": 1,
	"GD_EVENT":   2,
	"GD_MESSAGE": 3,
}

func (x EntityProto_Kind) Enum() *EntityProto_Kind {
	p := new(EntityProto_Kind)
	*p = x
	return p
}
func (x EntityProto_Kind) String() string {
	return proto.EnumName(EntityProto_Kind_name, int32(x))
}
func (x *EntityProto_Kind) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EntityProto_Kind_value, data, "EntityProto_Kind")
	if err != nil {
		return err
	}
	*x = EntityProto_Kind(value)
	return nil
}

type Index_Property_Direction int32

const (
	Index_Property_ASCENDING  Index_Property_Direction = 1
	Index_Property_DESCENDING Index_Property_Direction = 2
)

var Index_Property_Direction_name = map[int32]string{
	1: "ASCENDING",
	2: "DESCENDING",
}
var Index_Property_Direction_value = map[string]int32{
	"ASCENDING":  1,
	"DESCENDING": 2,
}

func (x Index_Property_Direction) Enum() *Index_Property_Direction {
	p := new(Index_Property_Direction)
	*p = x
	return p
}
func (x Index_Property_Direction) String() string {
	return proto.EnumName(Index_Property_Direction_name, int32(x))
}
func (x *Index_Property_Direction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Index_Property_Direction_value, data, "Index_Property_Direction")
	if err != nil {
		return err
	}
	*x = Index_Property_Direction(value)
	return nil
}

type CompositeIndex_State int32

const (
	CompositeIndex_WRITE_ONLY CompositeIndex_State = 1
	CompositeIndex_READ_WRITE CompositeIndex_State = 2
	CompositeIndex_DELETED    CompositeIndex_State = 3
	CompositeIndex_ERROR      CompositeIndex_State = 4
)

var CompositeIndex_State_name = map[int32]string{
	1: "WRITE_ONLY",
	2: "READ_WRITE",
	3: "DELETED",
	4: "ERROR",
}
var CompositeIndex_State_value = map[string]int32{
	"WRITE_ONLY": 1,
	"READ_WRITE": 2,
	"DELETED":    3,
	"ERROR":      4,
}

func (x CompositeIndex_State) Enum() *CompositeIndex_State {
	p := new(CompositeIndex_State)
	*p = x
	return p
}
func (x CompositeIndex_State) String() string {
	return proto.EnumName(CompositeIndex_State_name, int32(x))
}
func (x *CompositeIndex_State) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CompositeIndex_State_value, data, "CompositeIndex_State")
	if err != nil {
		return err
	}
	*x = CompositeIndex_State(value)
	return nil
}

type Snapshot_Status int32

const (
	Snapshot_INACTIVE Snapshot_Status = 0
	Snapshot_ACTIVE   Snapshot_Status = 1
)

var Snapshot_Status_name = map[int32]string{
	0: "INACTIVE",
	1: "ACTIVE",
}
var Snapshot_Status_value = map[string]int32{
	"INACTIVE": 0,
	"ACTIVE":   1,
}

func (x Snapshot_Status) Enum() *Snapshot_Status {
	p := new(Snapshot_Status)
	*p = x
	return p
}
func (x Snapshot_Status) String() string {
	return proto.EnumName(Snapshot_Status_name, int32(x))
}
func (x *Snapshot_Status) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Snapshot_Status_value, data, "Snapshot_Status")
	if err != nil {
		return err
	}
	*x = Snapshot_Status(value)
	return nil
}

type Query_Hint int32

const (
	Query_ORDER_FIRST    Query_Hint = 1
	Query_ANCESTOR_FIRST Query_Hint = 2
	Query_FILTER_FIRST   Query_Hint = 3
)

var Query_Hint_name = map[int32]string{
	1: "ORDER_FIRST",
	2: "ANCESTOR_FIRST",
	3: "FILTER_FIRST",
}
var Query_Hint_value = map[string]int32{
	"ORDER_FIRST":    1,
	"ANCESTOR_FIRST": 2,
	"FILTER_FIRST":   3,
}

func (x Query_Hint) Enum() *Query_Hint {
	p := new(Query_Hint)
	*p = x
	return p
}
func (x Query_Hint) String() string {
	return proto.EnumName(Query_Hint_name, int32(x))
}
func (x *Query_Hint) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Query_Hint_value, data, "Query_Hint")
	if err != nil {
		return err
	}
	*x = Query_Hint(value)
	return nil
}

type Query_Filter_Operator int32

const (
	Query_Filter_LESS_THAN             Query_Filter_Operator = 1
	Query_Filter_LESS_THAN_OR_EQUAL    Query_Filter_Operator = 2
	Query_Filter_GREATER_THAN          Query_Filter_Operator = 3
	Query_Filter_GREATER_THAN_OR_EQUAL Query_Filter_Operator = 4
	Query_Filter_EQUAL                 Query_Filter_Operator = 5
	Query_Filter_IN                    Query_Filter_Operator = 6
	Query_Filter_EXISTS                Query_Filter_Operator = 7
)

var Query_Filter_Operator_name = map[int32]string{
	1: "LESS_THAN",
	2: "LESS_THAN_OR_EQUAL",
	3: "GREATER_THAN",
	4: "GREATER_THAN_OR_EQUAL",
	5: "EQUAL",
	6: "IN",
	7: "EXISTS",
}
var Query_Filter_Operator_value = map[string]int32{
	"LESS_THAN":             1,
	"LESS_THAN_OR_EQUAL":    2,
	"GREATER_THAN":          3,
	"GREATER_THAN_OR_EQUAL": 4,
	"EQUAL":                 5,
	"IN":                    6,
	"EXISTS":                7,
}

func (x Query_Filter_Operator) Enum() *Query_Filter_Operator {
	p := new(Query_Filter_Operator)
	*p = x
	return p
}
func (x Query_Filter_Operator) String() string {
	return proto.EnumName(Query_Filter_Operator_name, int32(x))
}
func (x *Query_Filter_Operator) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Query_Filter_Operator_value, data, "Query_Filter_Operator")
	if err != nil {
		return err
	}
	*x = Query_Filter_Operator(value)
	return nil
}

type Query_Order_Direction int32

const (
	Query_Order_ASCENDING  Query_Order_Direction = 1
	Query_Order_DESCENDING Query_Order_Direction = 2
)

var Query_Order_Direction_name = map[int32]string{
	1: "ASCENDING",
	2: "DESCENDING",
}
var Query_Order_Direction_value = map[string]int32{
	"ASCENDING":  1,
	"DESCENDING": 2,
}

func (x Query_Order_Direction) Enum() *Query_Order_Direction {
	p := new(Query_Order_Direction)
	*p = x
	return p
}
func (x Query_Order_Direction) String() string {
	return proto.EnumName(Query_Order_Direction_name, int32(x))
}
func (x *Query_Order_Direction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Query_Order_Direction_value, data, "Query_Order_Direction")
	if err != nil {
		return err
	}
	*x = Query_Order_Direction(value)
	return nil
}

type Error_ErrorCode int32

const (
	Error_BAD_REQUEST                  Error_ErrorCode = 1
	Error_CONCURRENT_TRANSACTION       Error_ErrorCode = 2
	Error_INTERNAL_ERROR               Error_ErrorCode = 3
	Error_NEED_INDEX                   Error_ErrorCode = 4
	Error_TIMEOUT                      Error_ErrorCode = 5
	Error_PERMISSION_DENIED            Error_ErrorCode = 6
	Error_BIGTABLE_ERROR               Error_ErrorCode = 7
	Error_COMMITTED_BUT_STILL_APPLYING Error_ErrorCode = 8
	Error_CAPABILITY_DISABLED          Error_ErrorCode = 9
	Error_TRY_ALTERNATE_BACKEND        Error_ErrorCode = 10
	Error_SAFE_TIME_TOO_OLD            Error_ErrorCode = 11
)

var Error_ErrorCode_name = map[int32]string{
	1:  "BAD_REQUEST",
	2:  "CONCURRENT_TRANSACTION",
	3:  "INTERNAL_ERROR",
	4:  "NEED_INDEX",
	5:  "TIMEOUT",
	6:  "PERMISSION_DENIED",
	7:  "BIGTABLE_ERROR",
	8:  "COMMITTED_BUT_STILL_APPLYING",
	9:  "CAPABILITY_DISABLED",
	10: "TRY_ALTERNATE_BACKEND",
	11: "SAFE_TIME_TOO_OLD",
}
var Error_ErrorCode_value = map[string]int32{
	"BAD_REQUEST":                  1,
	"CONCURRENT_TRANSACTION":       2,
	"INTERNAL_ERROR":               3,
	"NEED_INDEX":                   4,
	"TIMEOUT":                      5,
	"PERMISSION_DENIED":            6,
	"BIGTABLE_ERROR":               7,
	"COMMITTED_BUT_STILL_APPLYING": 8,
	"CAPABILITY_DISABLED":          9,
	"TRY_ALTERNATE_BACKEND":        10,
	"SAFE_TIME_TOO_OLD":            11,
}

func (x Error_ErrorCode) Enum() *Error_ErrorCode {
	p := new(Error_ErrorCode)
	*p = x
	return p
}
func (x Error_ErrorCode) String() string {
	return proto.EnumName(Error_ErrorCode_name, int32(x))
}
func (x *Error_ErrorCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Error_ErrorCode_value, data, "Error_ErrorCode")
	if err != nil {
		return err
	}
	*x = Error_ErrorCode(value)
	return nil
}

type PutRequest_AutoIdPolicy int32

const (
	PutRequest_CURRENT    PutRequest_AutoIdPolicy = 0
	PutRequest_SEQUENTIAL PutRequest_AutoIdPolicy = 1
)

var PutRequest_AutoIdPolicy_name = map[int32]string{
	0: "CURRENT",
	1: "SEQUENTIAL",
}
var PutRequest_AutoIdPolicy_value = map[string]int32{
	"CURRENT":    0,
	"SEQUENTIAL": 1,
}

func (x PutRequest_AutoIdPolicy) Enum() *PutRequest_AutoIdPolicy {
	p := new(PutRequest_AutoIdPolicy)
	*p = x
	return p
}
func (x PutRequest_AutoIdPolicy) String() string {
	return proto.EnumName(PutRequest_AutoIdPolicy_name, int32(x))
}
func (x *PutRequest_AutoIdPolicy) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PutRequest_AutoIdPolicy_value, data, "PutRequest_AutoIdPolicy")
	if err != nil {
		return err
	}
	*x = PutRequest_AutoIdPolicy(value)
	return nil
}

type Action struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()    {}

type PropertyValue struct {
	Int64Value       *int64                        `protobuf:"varint,1,opt,name=int64Value" json:"int64Value,omitempty"`
	BooleanValue     *bool                         `protobuf:"varint,2,opt,name=booleanValue" json:"booleanValue,omitempty"`
	StringValue      *string                       `protobuf:"bytes,3,opt,name=stringValue" json:"stringValue,omitempty"`
	DoubleValue      *float64                      `protobuf:"fixed64,4,opt,name=doubleValue" json:"doubleValue,omitempty"`
	Pointvalue       *PropertyValue_PointValue     `protobuf:"group,5,opt,name=PointValue" json:"pointvalue,omitempty"`
	Uservalue        *PropertyValue_UserValue      `protobuf:"group,8,opt,name=UserValue" json:"uservalue,omitempty"`
	Referencevalue   *PropertyValue_ReferenceValue `protobuf:"group,12,opt,name=ReferenceValue" json:"referencevalue,omitempty"`
	XXX_unrecognized []byte                        `json:"-"`
}

func (m *PropertyValue) Reset()         { *m = PropertyValue{} }
func (m *PropertyValue) String() string { return proto.CompactTextString(m) }
func (*PropertyValue) ProtoMessage()    {}

func (m *PropertyValue) GetInt64Value() int64 {
	if m != nil && m.Int64Value != nil {
		return *m.Int64Value
	}
	return 0
}

func (m *PropertyValue) GetBooleanValue() bool {
	if m != nil && m.BooleanValue != nil {
		return *m.BooleanValue
	}
	return false
}

func (m *PropertyValue) GetStringValue() string {
	if m != nil && m.StringValue != nil {
		return *m.StringValue
	}
	return ""
}

func (m *PropertyValue) GetDoubleValue() float64 {
	if m != nil && m.DoubleValue != nil {
		return *m.DoubleValue
	}
	return 0
}

func (m *PropertyValue) GetPointvalue() *PropertyValue_PointValue {
	if m != nil {
		return m.Pointvalue
	}
	return nil
}

func (m *PropertyValue) GetUservalue() *PropertyValue_UserValue {
	if m != nil {
		return m.Uservalue
	}
	return nil
}

func (m *PropertyValue) GetReferencevalue() *PropertyValue_ReferenceValue {
	if m != nil {
		return m.Referencevalue
	}
	return nil
}

type PropertyValue_PointValue struct {
	X                *float64 `protobuf:"fixed64,6,req,name=x" json:"x,omitempty"`
	Y                *float64 `protobuf:"fixed64,7,req,name=y" json:"y,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PropertyValue_PointValue) Reset()         { *m = PropertyValue_PointValue{} }
func (m *PropertyValue_PointValue) String() string { return proto.CompactTextString(m) }
func (*PropertyValue_PointValue) ProtoMessage()    {}

func (m *PropertyValue_PointValue) GetX() float64 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *PropertyValue_PointValue) GetY() float64 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

type PropertyValue_UserValue struct {
	Email             *string `protobuf:"bytes,9,req,name=email" json:"email,omitempty"`
	AuthDomain        *string `protobuf:"bytes,10,req,name=auth_domain" json:"auth_domain,omitempty"`
	Nickname          *string `protobuf:"bytes,11,opt,name=nickname" json:"nickname,omitempty"`
	FederatedIdentity *string `protobuf:"bytes,21,opt,name=federated_identity" json:"federated_identity,omitempty"`
	FederatedProvider *string `protobuf:"bytes,22,opt,name=federated_provider" json:"federated_provider,omitempty"`
	XXX_unrecognized  []byte  `json:"-"`
}

func (m *PropertyValue_UserValue) Reset()         { *m = PropertyValue_UserValue{} }
func (m *PropertyValue_UserValue) String() string { return proto.CompactTextString(m) }
func (*PropertyValue_UserValue) ProtoMessage()    {}

func (m *PropertyValue_UserValue) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *PropertyValue_UserValue) GetAuthDomain() string {
	if m != nil && m.AuthDomain != nil {
		return *m.AuthDomain
	}
	return ""
}

func (m *PropertyValue_UserValue) GetNickname() string {
	if m != nil && m.Nickname != nil {
		return *m.Nickname
	}
	return ""
}

func (m *PropertyValue_UserValue) GetFederatedIdentity() string {
	if m != nil && m.FederatedIdentity != nil {
		return *m.FederatedIdentity
	}
	return ""
}

func (m *PropertyValue_UserValue) GetFederatedProvider() string {
	if m != nil && m.FederatedProvider != nil {
		return *m.FederatedProvider
	}
	return ""
}

type PropertyValue_ReferenceValue struct {
	App              *string                                     `protobuf:"bytes,13,req,name=app" json:"app,omitempty"`
	NameSpace        *string                                     `protobuf:"bytes,20,opt,name=name_space" json:"name_space,omitempty"`
	Pathelement      []*PropertyValue_ReferenceValue_PathElement `protobuf:"group,14,rep,name=PathElement" json:"pathelement,omitempty"`
	XXX_unrecognized []byte                                      `json:"-"`
}

func (m *PropertyValue_ReferenceValue) Reset()         { *m = PropertyValue_ReferenceValue{} }
func (m *PropertyValue_ReferenceValue) String() string { return proto.CompactTextString(m) }
func (*PropertyValue_ReferenceValue) ProtoMessage()    {}

func (m *PropertyValue_ReferenceValue) GetApp() string {
	if m != nil && m.App != nil {
		return *m.App
	}
	return ""
}

func (m *PropertyValue_ReferenceValue) GetNameSpace() string {
	if m != nil && m.NameSpace != nil {
		return *m.NameSpace
	}
	return ""
}

func (m *PropertyValue_ReferenceValue) GetPathelement() []*PropertyValue_ReferenceValue_PathElement {
	if m != nil {
		return m.Pathelement
	}
	return nil
}

type PropertyValue_ReferenceValue_PathElement struct {
	Type             *string `protobuf:"bytes,15,req,name=type" json:"type,omitempty"`
	Id               *int64  `protobuf:"varint,16,opt,name=id" json:"id,omitempty"`
	Name             *string `protobuf:"bytes,17,opt,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PropertyValue_ReferenceValue_PathElement) Reset() {
	*m = PropertyValue_ReferenceValue_PathElement{}
}
func (m *PropertyValue_ReferenceValue_PathElement) String() string { return proto.CompactTextString(m) }
func (*PropertyValue_ReferenceValue_PathElement) ProtoMessage()    {}

func (m *PropertyValue_ReferenceValue_PathElement) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *PropertyValue_ReferenceValue_PathElement) GetId() int64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *PropertyValue_ReferenceValue_PathElement) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type Property struct {
	Meaning               *Property_Meaning               `protobuf:"varint,1,opt,name=meaning,enum=appengine.Property_Meaning,def=0" json:"meaning,omitempty"`
	MeaningUri            *string                         `protobuf:"bytes,2,opt,name=meaning_uri" json:"meaning_uri,omitempty"`
	Name                  *string                         `protobuf:"bytes,3,req,name=name" json:"name,omitempty"`
	Value                 *PropertyValue                  `protobuf:"bytes,5,req,name=value" json:"value,omitempty"`
	Multiple              *bool                           `protobuf:"varint,4,req,name=multiple" json:"multiple,omitempty"`
	Searchable            *bool                           `protobuf:"varint,6,opt,name=searchable,def=0" json:"searchable,omitempty"`
	FtsTokenizationOption *Property_FtsTokenizationOption `protobuf:"varint,8,opt,name=fts_tokenization_option,enum=appengine.Property_FtsTokenizationOption" json:"fts_tokenization_option,o