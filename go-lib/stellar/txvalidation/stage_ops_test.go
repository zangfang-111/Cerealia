package txvalidation

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

// stage close req create

const stageCloseReqCreateTX1 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzEzMwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAPc3RhZ2VfY2xvc2VSZXFzAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABA811D3+d1WtBtsn0AlE79223CRtbTjKM7Z1xXb8G2UaYZFx6NHeQG2gF1fCyOAJUp+n7j0RdkimqOhO73YtkmAA=="
const stageCloseReqCreateTX2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAPc3RhZ2VfY2xvc2VSZXFzAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABA5IPf64xxqSqUNVnrWrP86pcCaq4mPRgmb6ZuubiUpwz5r1Ej3ZSJQilZrFuB34b6d2a8NPagBLbabiSavmbjAQ=="
const stageCloseReqCreateTXInvalid = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAXc3RhZ2UuY2xvc2VSZXFzLmludmFsaWQAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAAAAAABF7klLwAAAEBYKDqnkswqzhKrh9U5RgtpA4CvBi5k2sy8Pmh3JBKiBGeWteESg9fYsWbi/Y8viQ+HSCJ9GKnL19WlG06nWJkI"
const stageCloseReqCreateTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAADOTg3AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAA9zdGFnZS5jbG9zZVJlcXMAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAAAAAABF7klLwAAAEBkcEocQ08It4pFEhEtG7flyRUQp5qGJKGKDemMhP78qLkODduAVrxA0SFgkS14Pg6d0EztM/TtaZo5WXBGgigO"

func (s *TxValidationSuite) TestValidateStageCloseReqCreateTXGreenPath(c *C) {
	eBuilder, se, err := ValidateStageCloseReqTX(stageCloseReqCreateTX1, 133, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)

	eBuilder, se, err = ValidateStageCloseReqTX(stageCloseReqCreateTX2, 987, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateStageCloseReqCreateTXBadID(c *C) {
	_, _, err := ValidateStageCloseReqTX(stageCloseReqCreateTX1, 123, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badData)
}

func (s *TxValidationSuite) TestValidateStageCloseReqCreateTXInvalid(c *C) {
	_, _, err := ValidateStageCloseReqTX(stageCloseReqCreateTXInvalid, 987, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageCloseReqTX(stageCloseReqCreateTXWithHash, 987, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)
}

// stage close req approve

const stageCloseReqApproveTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAPc3RhZ2VfY2xvc2VSZXFzAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAIYXBwcm92ZWQAAAAAAAAAARe5JS8AAABAeTWilp34TOPbLXKibka/P4EXy30f0uRgQVCiDj3/vYQqtUP0nN46UFlg/ervIYGq5lXuUSPW3oeoOgb7hz9iDw=="
const stageCloseReqApproveTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAADOTg3AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAA9zdGFnZS5jbG9zZVJlcXMAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhhcHByb3ZlZAAAAAAAAAABF7klLwAAAEAJRkY/KhmviPgpxjFP0w050ZXVAc76IILDrb0GMxEPTMejkJG51zALnXMemqEuS5ByO8ThoBOrDyT2+/tGCvUA"

func (s *TxValidationSuite) TestValidateStageCloseReqApproveTXGreenPath(c *C) {
	eBuilder, se, err := ValidateStageCloseReqTX(stageCloseReqApproveTX, 987, stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateStageCloseReqApproveTXBadID(c *C) {
	_, _, err := ValidateStageCloseReqTX(stageCloseReqApproveTX, 123, stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageCloseReqTX(stageCloseReqApproveTXWithHash, 987, stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)
}

// stage close req reject

const stageCloseReqRejectTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAPc3RhZ2VfY2xvc2VSZXFzAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAIcmVqZWN0ZWQAAAAAAAAAARe5JS8AAABAPR9Rrq7L9f1yGywXlB+1ZAj+gtv9V7mbfQfMZAZmqi6Py7imOK0/w/4wkkF1BscqDVr7Slce5/ZleGF60aKqAw=="
const stageCloseReqRejectTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAOqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAM5ODcAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAAD3N0YWdlLmNsb3NlUmVxcwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACHJlamVjdGVkAAAAAAAAAAEXuSUvAAAAQIokh0BGS42my8TDqFTaFLEbX4T7BeW+KmpG17sQW9OxH6O5Z42TJ5XTqSNXdWSKPI0p/snjrjkshv+46u/i5Ak="

func (s *TxValidationSuite) TestValidateStageCloseReqRejectTXGreenPath(c *C) {
	eBuilder, se, err := ValidateStageCloseReqTX(stageCloseReqRejectTX, 987, stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateStageCloseReqRejectTXBadID(c *C) {
	_, _, err := ValidateStageCloseReqTX(stageCloseReqRejectTX, 123, stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageCloseReqTX(stageCloseReqRejectTXWithHash, 987, stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, ErrorContains, badMemoErr)
	c.Assert(err, ErrorContains, badData)
}

// stage add req create

const stageAddReqCreateTX1 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2FkZAAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzEzMwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAAB3BlbmRpbmcAAAAAAAAAAAEXuSUvAAAAQAGqCrNdCYy/il6gA468545Kk8Fg5u/vkeorU3qgJiWfZry3JGPuHq+Hg0ZHHsLc1F1Ms1DeFfvf4+EoPlBi6wY="
const stageAddReqCreateTX2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2FkZAAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAAB3BlbmRpbmcAAAAAAAAAAAEXuSUvAAAAQAP2yDrVzueDIPFo/wKxdju9nQ1agSRWf4C7Y5A7nAj2C5/0DYdMYt/lEfvoG6ouMELJ9YysvRo3J+ftrA9I/AY="
const stageAddReqCreateTXInvalid = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAADXN0YWdlX2FkZC5iYWQAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAM5ODcAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAAAAAABF7klLwAAAEA0SbkfcHZ17f9j1fwZG+dHvTeB0Tv279FAUWFvF1GNkgrBodXOPxcYBpPJQ/c8vqVjN8kd7m/0B+7c34arTEUP"
const stageAddReqCreateTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2FkZAAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAAB3BlbmRpbmcAAAAAAAAAAAEXuSUvAAAAQLbpsl61n/k006/383E745HnTWjXbjLOzDhWOkRm/V9csnAiEtK7loqmXx1VANzg+/XxdUUMF1s5PUfGBhytUQQ="

func (s *TxValidationSuite) TestValidateStageAddReqCreateTXGreenPath(c *C) {
	eBuilder, se, err := ValidateStageAddReqTX(stageAddReqCreateTX1, 133, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateStageAddReqTX(stageAddReqCreateTX2, 987, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateStageAddReqCreateTXBadID(c *C) {
	_, _, err := ValidateStageAddReqTX(stageAddReqCreateTX1, 123, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageAddReqTX(stageAddReqCreateTXInvalid, 987, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageAddReqTX(stageAddReqCreateTXWithHash, 987, stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badMemoErr)
}

// stage Add req approve

const stageAddReqApproveTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAlzdGFnZV9hZGQAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAM5ODcAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhhcHByb3ZlZAAAAAAAAAABF7klLwAAAEBKuxHOvcyGqfrXhmFkwUtuck7aat2CmbPLQ889uehTkScJEVuViLkr/LjoKG+qUYOweYUaw94CU1P3b+zAz7gG"
const stageAddReqApproveTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2FkZAAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACGFwcHJvdmVkAAAAAAAAAAEXuSUvAAAAQC2Cd3j4MaM6OVqPu1l4inLufHs182rT0dKp5rf3c8zet/3P07jCSxIhoj71D+HEzDQMJ7VLybBfehRxprV7tgg="

func (s *TxValidationSuite) TestValidateStageAddReqApproveTXGreenPath(c *C) {
	eBuilder, se, err := ValidateStageAddReqTX(stageAddReqApproveTX, 987, stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateStageAddReqApproveTXBadID(c *C) {
	_, _, err := ValidateStageAddReqTX(stageAddReqApproveTX, 123, stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageAddReqTX(stageAddReqApproveTXWithHash, 987, stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, ErrorContains, badMemoErr)
}

// stage Add req reject

const stageAddReqRejectTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAlzdGFnZV9hZGQAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAM5ODcAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAhyZWplY3RlZAAAAAAAAAABF7klLwAAAECwKPoBRM93tjjm8SMZkIhxVQyyfD6sF230nRf+5rlqbA4vBBN01agBptuEiorwUJdfr4igeb8+CsTMkh7TSG0B"
const stageAddReqRejectTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlX2FkZAAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACHJlamVjdGVkAAAAAAAAAAEXuSUvAAAAQFNwPybCCLQyk1QDBCTzKpedh6pufTzV7c75IJjaHR1lm+P6cFBcsxrPdvTi2/XRwyDmcZPJ0EZldYH84fT8ngA="

func (s *TxValidationSuite) TestValidateStageAddReqRejectTXGreenPath(c *C) {
	eBuilder, se, err := ValidateStageAddReqTX(stageAddReqRejectTX, 987, stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateStageAddReqRejectTXBadID(c *C) {
	_, _, err := ValidateStageAddReqTX(stageAddReqRejectTX, 123, stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateStageAddReqTX(stageAddReqRejectTXWithHash, 987, stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, ErrorContains, badMemoErr)
}
