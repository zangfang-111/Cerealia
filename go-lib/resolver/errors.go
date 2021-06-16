package resolver

import "github.com/robert-zaremba/errstack"

var (
	errSelfApprove = errstack.NewReq("You can't approve your own request")
	errChangeStage = errstack.NewReq("You can’t change deleted or closed stage")
)
