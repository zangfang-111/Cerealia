package txvalidation

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

// trade close req create

const tradeCloseReqCreateTX1 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAAD3RyYWRlX2Nsb3NlUmVxcwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxOTkzMTM0AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABAJcwXw60r4Mb+ljZBoX2V2+RMjH/0WIW+Kt3en7f0wUr2baWUMzFhzxrMmVmx2+nZNRwYCFN/o6G1fy0Yzo5gCg=="
const tradeCloseReqCreateTX2 = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAAD3RyYWRlX2Nsb3NlUmVxcwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxOTkzMTM1AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABAez2hiBLaN4wgq8Ff3TKom985PSg6k8+EAw55Na9r6aFvpBQE+d/U7THWHiNe237kdgBn72vSccyRjJ2YSwLmAg=="
const tradeCloseReqCreateTXInvalid = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAAA2lkeAAAAAABAAAAAzk4NwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAXc3RhZ2UuY2xvc2VSZXFzLmludmFsaWQAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAAAAAABF7klLwAAAEBYKDqnkswqzhKrh9U5RgtpA4CvBi5k2sy8Pmh3JBKiBGeWteESg9fYsWbi/Y8viQ+HSCJ9GKnL19WlG06nWJkI"
const tradeCloseReqCreateTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAAD3RyYWRlX2Nsb3NlUmVxcwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxOTkzMTM0AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABAowJQz8QZhBo5/xS8KUQ7MdYJvNmD/TXM3g2MaVMgkaGRDgV52lP4yb85Tyyh3CQY2j6ZVsq25Z35thJkR1ssAg=="

func (s *TxValidationSuite) TestValidateTradeCloseReqCreateTXGreenPath(c *C) {
	eBuilder, se, err := ValidateTradeCloseReqTX(tradeCloseReqCreateTX1, "1993134", stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	eBuilder, se, err = ValidateTradeCloseReqTX(tradeCloseReqCreateTX2, "1993135", stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateTradeCloseReqCreateTXBadID(c *C) {
	_, _, err := ValidateTradeCloseReqTX(tradeCloseReqCreateTX1, "199313", stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badData)

	// bad tx
	_, _, err = ValidateTradeCloseReqTX(tradeCloseReqCreateTXInvalid, "1993134", stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badData)

	// has hash
	_, _, err = ValidateTradeCloseReqTX(tradeCloseReqCreateTXWithHash, "1993134", stageActionValidTrade, user1, model.ApprovalPending)
	c.Assert(err, ErrorContains, badMemoErr)
}

// trade close req approve

const tradeCloseReqApproveTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAA90cmFkZV9jbG9zZVJlcXMAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAHMTk5MzEzNAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACGFwcHJvdmVkAAAAAAAAAAEXuSUvAAAAQDpwQE/Rmm/Ydbnr+TIAgWEXIghwcMBqTktLK5Hjh9YaKxw74Obw5iP1m3jo18Ih0k8+Dlf1u+A6DnqtZNf03AI="
const tradeCloseReqApproveTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAAD3RyYWRlX2Nsb3NlUmVxcwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxOTkzMTM0AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAIYXBwcm92ZWQAAAAAAAAAARe5JS8AAABAO3MtVokIMlrW+hnjW2i5awlIgxUSEMbaEXXnuJXvpqNt5rKxmQpcmJfaP7ekHSIwLkT6ykW7lbLOs2YeTjWbCg=="

func (s *TxValidationSuite) TestValidateTradeCloseReqApproveTXGreenPath(c *C) {
	eBuilder, se, err := ValidateTradeCloseReqTX(tradeCloseReqApproveTX, "1993134", stageActionValidTrade, user1, model.ApprovalApproved)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
	c.Assert(err, IsNil)
}

func (s *TxValidationSuite) TestValidateTradeCloseReqApproveTXNegative(c *C) {
	_, _, err := ValidateTradeCloseReqTX(tradeCloseReqApproveTX, "123", stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateTradeCloseReqTX(tradeCloseReqApproveTXWithHash, "1993134", stageActionValidTrade, user1, model.ApprovalApproved)
	c.Assert(err, ErrorContains, badMemoErr)
}

// trade close req reject

const tradeCloseReqRejectTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAA90cmFkZV9jbG9zZVJlcXMAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAAHMTk5MzEzNAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAACHJlamVjdGVkAAAAAAAAAAEXuSUvAAAAQJNB5ZV6hFVDUYU2mueoh5SymRvfuVCVUKZ2x0tQrGuJUrLd+KHemBGenWRRLlg7zMAM7P03dSgY9KGEicJA7gA="
const tradeCloseReqRejectTXWithHash = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAAFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADq/EjQVI6qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqoAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAAD3RyYWRlX2Nsb3NlUmVxcwAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAcxOTkzMTM0AAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAIcmVqZWN0ZWQAAAAAAAAAARe5JS8AAABAPMW7NOme/42Cd0rBXJHALBd9dkAFQfjJdA3POFxhiggtQH1YOJnkCtqDiUqhYahZr1TTwY/qNsNAARXLs64oAw=="

func (s *TxValidationSuite) TestValidateTradeCloseReqRejectTXGreenPath(c *C) {
	eBuilder, se, err := ValidateTradeCloseReqTX(tradeCloseReqRejectTX, "1993134", stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, IsNil)
	c.Check(eBuilder, NotNil)
	c.Check(se, NotNil)
}

func (s *TxValidationSuite) TestValidateTradeCloseReqRejectTXNegative(c *C) {
	_, _, err := ValidateTradeCloseReqTX(tradeCloseReqRejectTX, "123", stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, ErrorContains, badData)

	_, _, err = ValidateTradeCloseReqTX(tradeCloseReqRejectTXWithHash, "1993134", stageActionValidTrade, user1, model.ApprovalRejected)
	c.Assert(err, ErrorContains, badMemoErr)
}
