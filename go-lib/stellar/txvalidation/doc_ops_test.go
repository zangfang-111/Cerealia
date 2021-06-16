package txvalidation

import (
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/validation"
	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func (s *TxValidationSuite) TestValidateDataDocAddGreenPath(c *C) {
	vb := &validation.Builder{}
	d := dataMap{
		"aaa": {
			"idx":       "133:242",
			"entity":    model.TxTradeEntityStageDoc.String(),
			"operation": model.ApprovalPending.String(),
		}}

	outcome := validateDataDocOp(vb, d,
		"aaa",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, true)

	d["aaa"]["operation"] = model.ApprovalApproved.String()
	outcome = validateDataDocOp(vb, d,
		"aaa",
		path1,
		model.ApprovalApproved)
	c.Check(outcome, Equals, true)
}

func (s *TxValidationSuite) TestValidateDataDocAddBadIdx(c *C) {
	vb := &validation.Builder{}
	outcome := validateDataDocOp(vb,
		dataMap{
			"aaa": {
				"idx":       "2:242",
				"entity":    model.TxTradeEntityStageDoc.String(),
				"operation": model.ApprovalPending.String(),
			}},
		"aaa",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
	outcome = validateDataDocOp(vb,
		dataMap{
			"aaa": {
				"idx":       "133:1",
				"entity":    model.TxTradeEntityStageDoc.String(),
				"operation": model.ApprovalPending.String(),
			}},
		"aaa",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
}

func (s *TxValidationSuite) TestValidateDataDocAddBadValues(c *C) {
	vb := &validation.Builder{}
	outcome := validateDataDocOp(vb,
		dataMap{
			"aaa": {
				"idx":       "133:242",
				"entity":    "stage.doc.BAD",
				"operation": model.ApprovalPending.String(),
			}},
		"aaa",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
	outcome = validateDataDocOp(vb,
		dataMap{
			"bbb": {
				"idx":       "133:242",
				"entity":    model.TxTradeEntityStageDoc.String(),
				"operation": "pending.Bad",
			}},
		"bbb",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
}

func (s *TxValidationSuite) TestValidateDataDocAddNonmatchingTradeAccKey(c *C) {
	// trade account key must match the key that data is set to
	vb := &validation.Builder{}
	d := dataMap{
		"aaa": {
			"idx":       "133:242",
			"entity":    model.TxTradeEntityStageDoc.String(),
			"operation": model.ApprovalPending.String(),
		}}
	outcome := validateDataDocOp(vb, d,
		"bbb",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
	outcome = validateDataDocOp(vb, d,
		"ccc",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
	outcome = validateDataDocOp(vb,
		dataMap{
			"bbb": {
				"idx":       "133:242",
				"entity":    model.TxTradeEntityStageDoc.String(),
				"operation": model.ApprovalPending.String(),
			},
			"ccc": {
				"idx":       "133:242",
				"entity":    model.TxTradeEntityStageDoc.String(),
				"operation": model.ApprovalPending.String(),
			}},
		"ccc",
		path1,
		model.ApprovalPending)
	c.Check(outcome, Equals, false)
}

// doc add

const docAddTXValid = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABkAAEYqOEodhjAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADxX/YK/PkjdSVdbDuhbryh4jCPKIueGEyDKG2X+YuaL0AAAAEAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2RvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAABzEyMzo0NTYAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACmV4cGlyZVRpbWUAAAAAAAEAAAAKMTYwOTM3MjgwMAAAAAAAAAAAAAEXuSUvAAAAQOnyN+Hzg4kVm6EKzNkrDU11y1SeyKRQ1WbXwYsvPwOSDeCcTaYFSJR4jVBsZ+BfWDP7akm8aW3T9nBfNw/VmQs="
const docAddTxValidSubmitted = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADxX/YK/PkjdSVdbDuhbryh4jCPKIueGEyDKG2X+YuaL0AAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2RvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAABzEyMzo0NTYAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAlzdWJtaXR0ZWQAAAAAAAAAAAAAARe5JS8AAABAUCBBfp4Pn7E5ff7rAyovXXdVwa0/2hdnald4XaRQJMgCS5nZiJBrzDzSzc7qTu+iAavB6hCZ58hcUR7s4ug5BQ=="
const docAddTXValidDocHash = "c57fd82bf3e48dd49575b0ee85baf28788c23ca22e7861320ca1b65fe62e68bd"
const docAddTXBadData = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAPFf9gr8+SN1JV1sO6FuvKHiMI8oi54YTIMobZf5i5ovQAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAANCQUQAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlLmRvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABAQcRE5u1C48FvtV+BXEZfcOLb0pasa14LTVuNNyc8Z7bErwD3WQnGPRRL3ex0TxybmiMCTqUmjsdRd3EHLbepCQ=="
const docAddTXValid2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABkAAEYqOEodhjAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADxX/YK/PkjdSaqqqqqqqqqqrCPKIueGEyDKG2X+YuaL0AAAAEAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2RvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAABjk4NzoxMwAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACmV4cGlyZVRpbWUAAAAAAAEAAAAKMTYwOTM3MjgwMAAAAAAAAAAAAAEXuSUvAAAAQH4HKUk7BHYL4H3cH5Z41I+tuq93VC5YOT0GJBJwzrMXxuh5As6f+ClrTpdIe5KyS6VsMaGeMa5DbNXn1HiFRQo="
const docAddTxValidSubmitted2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADxX/YK/PkjdSaqqqqqqqqqqrCPKIueGEyDKG2X+YuaL0AAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2RvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAABjk4NzoxMwAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAlzdWJtaXR0ZWQAAAAAAAAAAAAAARe5JS8AAABAVwwnS20wYu0u5+Yt4jA6eQi/hETx47wkWM/a/u1GFUTnGBcriZ96ALOkuhVYdx5PLkfY7f5p3IW8Qf12k9Q8BA=="
const docAddTXValid2DocHash = "c57fd82bf3e48dd49aaaaaaaaaaaaaaaaac23ca22e7861320ca1b65fe62e68bd"
const docAddTXTextHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAPFf9gr8+SN1JV1sO6FuvKHiMI8oi54YTIMobZf5i5ovQAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAANiYWQAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlLmRvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAJY29uZmlybWVkAAAAAAAAAAAAAAEXuSUvAAAAQLoFFzOW6R1SN1sNMmE+FrCRgxtpuC4fESIsDYWurPwDMWXkMrStW6c2sRJlSjo+F4ONVWb3muqC4zOIAqN3TQ8="

func (s *TxValidationSuite) TestValidateDocAddTXGreenPath(c *C) {
	eBuilder, se, err := ValidateDocAddExpireTX(docAddTXValid, path2, stageActionValidTrade, user1, docAddTXValidDocHash, s.sampleExpireTime, s.sampleExpireTime.Add(-time.Second))
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateDocAddExpireTX(docAddTXValid2, path3, stageActionValidTrade, user1, docAddTXValid2DocHash, s.sampleExpireTime, s.sampleExpireTime.Add(-time.Second))
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateDocAddTX(docAddTxValidSubmitted, path2, stageActionValidTrade, user1, docAddTXValidDocHash)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateDocAddTX(docAddTxValidSubmitted2, path3, stageActionValidTrade, user1, docAddTXValid2DocHash)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateDocAddTXNegative(c *C) {
	// invalid doc hash checked with a well-formed hash from tx
	_, _, err := ValidateDocAddExpireTX(docAddTXValid, path2, stageActionValidTrade, user1, "invalid-doc-hash", s.sampleExpireTime, s.sampleExpireTime.Add(-time.Second))
	c.Assert(err, ErrorContains, badMemoErr)

	// invalid data
	_, _, err = ValidateDocAddExpireTX(docAddTXBadData, path1, stageActionValidTrade, user1, docAddTXValidDocHash, s.sampleExpireTime, s.sampleExpireTime.Add(-time.Second))
	c.Assert(err, ErrorContains, badData)

	// bad hash type
	_, _, err = ValidateDocRejectTX(docAddTXTextHash, path4TextHash, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)
}

// doc approval

const docConfirmTXValid = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAHMTIzOjQ1NgAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2VfZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhhcHByb3ZlZAAAAAAAAAABF7klLwAAAECVSS5Aae3DnEypSarM9qSBxm8kkMeoQkzntAcbeINT3Pff6KxIXkcEN9PSTq/iZ99CV3/soIRcBUGkBHhkTJ0N"
const docConfirmTXValid2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqszMzMzMzMzMzMzMzMzMzMzMzMzMzMzMwAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAGOTg3OjEzAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2VfZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhhcHByb3ZlZAAAAAAAAAABF7klLwAAAEDrBLaO3rl/e2j4WkZgAUXFZgnOiOybgZcNmAttOiUnpwL7IGZXzOkjPpmkYsLYfUnWGQMNpDKSLV7hcVDuGtUP"
const docConfirmTXBadData = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAPFf9gr8+SN1JV1sO6FuvKHiMI8oi54YTIMobZf5i5ovQAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAANiYWQAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlLmRvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAJY29uZmlybWVkAAAAAAAAAAAAAAEXuSUvAAAAQLoFFzOW6R1SN1sNMmE+FrCRgxtpuC4fESIsDYWurPwDMWXkMrStW6c2sRJlSjo+F4ONVWb3muqC4zOIAqN3TQ8="
const docConfirmTXWithInvalidHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAGOTg3OjEzAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2UuZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhhcHByb3ZlZAAAAAAAAAABF7klLwAAAEB+CZyrQJX4nyAHB2bilZBzEVSEIPuWiSDnXHDuDU92KQkLl54JHYe9/gbBEYXQwpfQL1AfPAnv+e5jdo97MsgM"
const docApproveTXWithTextHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAABAAAAC2hlbGxvIHdvcmxkAAAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxMjM6NDU2AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAlzdGFnZS5kb2MAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACGFwcHJvdmVkAAAAAAAAAAEXuSUvAAAAQP9YhazbdGFP/AVIg0pZWGgFwVCQc9Xhv3zSww0OizCH7FWOprUXDxC88IkPt6omvXMDOdYJZPTjzzytVGjl2gM="

func (s *TxValidationSuite) TestValidateDocConfirmTXGreenPath(c *C) {
	eBuilder, se, err := ValidateDocApproveTX(docConfirmTXValid, path2, stageActionValidTrade, user1)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateDocApproveTX(docConfirmTXValid2, path3, stageActionValidTrade, user1)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateDocConfirmTXNegative(c *C) {
	// InvalidData
	_, _, err := ValidateDocApproveTX(docConfirmTXBadData, path1, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)

	// invalid hash
	_, _, err = ValidateDocApproveTX(docConfirmTXWithInvalidHash, path3, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)

	// invalid type
	_, _, err = ValidateDocRejectTX(docApproveTXWithTextHash, path4TextHash, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)
}

// doc rejection

const docRejectTXValid1 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqszMzMzMzMzMzMzMzMzMzMzMzMzMzMzMwAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAGOTg3OjEzAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2VfZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhyZWplY3RlZAAAAAAAAAABF7klLwAAAEAQM2GyDm+0PGKSHyWm+902TlHgOXqTea/CNuR7mWuyAiYthugd4a5BtQ9tbOw7BXfEcAVG3wCY/qTMIzMw7z4C"
const docRejectTXValid2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAHMTIzOjQ1NgAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2VfZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhyZWplY3RlZAAAAAAAAAABF7klLwAAAECe/J6W0hmeqljK8ZfJcbb3luvDZS3hy5jxQNFNutJeSU8za49KneZVW3QeFVWUQxhZbpTVq0Pd/i+2GH74oBIK"
const docRejectTXInvalidData = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAIMTIzOjQ1NjEAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2UuZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhyZWplY3RlZAAAAAAAAAABF7klLwAAAEDvqKWeJpT0xBSfWENdv4G/RN9zaeL+/Mq6fKGfLSc3SOMDjDOaF48Dp+GJ6jJiCxnVmNItu6PgoFZNFxS3WXEL"
const docRejectTXWithInvalidHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAOqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxMjM6NDU2AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAlzdGFnZS5kb2MAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACHJlamVjdGVkAAAAAAAAAAEXuSUvAAAAQJcLgfX7bl98MBOFLj4uwDoUL4VwrUnmSsKyH3w4kXTkUI1Reic0Yy/I4YIFTPBLeDFyoRIqyjbJhJTF6ON5SQc="
const docRejectTXTextHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAABAAAAC2hlbGxvIHdvcmxkAAAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxMjM6NDU2AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAlzdGFnZS5kb2MAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACHJlamVjdGVkAAAAAAAAAAEXuSUvAAAAQBo4EoIQkrRKTlK75Ii84TwrWe1tlQSdWwov4MrTwKv0hFNrVpoxYVZSLxKzGLh14lNiq/7tMfZ2PimB15gVnQ4="

func (s *TxValidationSuite) TestValidateDocRejectTXGreenPath(c *C) {
	eBuilder, se, err := ValidateDocRejectTX(docRejectTXValid1, path3, stageActionValidTrade, user1)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateDocRejectTX(docRejectTXValid2, path2, stageActionValidTrade, user1)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateDocRejectTXNegative(c *C) {
	// bad data
	_, _, err := ValidateDocRejectTX(docRejectTXInvalidData, path2, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)

	// bad hash
	_, _, err = ValidateDocRejectTX(docRejectTXWithInvalidHash, path2, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)

	// bad type
	_, _, err = ValidateDocRejectTX(docRejectTXTextHash, path4TextHash, stageActionValidTrade, user1)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)
}
