package model

// BuildStages converts TradeStageTemplate stages to trade stages
func (tt TradeTemplate) BuildStages() []TradeStage {
	var stages []TradeStage
	ss := tt.Stages
	for i := range tt.Stages {
		stages = append(stages, NewTradeStage(
			ss[i].Name, ss[i].Description, -1, ss[i].Owner))
	}
	return stages
}

// NewTradeStage constructs TradeStage
func NewTradeStage(name, description string, addReqIdx int, owner TradeActor) TradeStage {
	return TradeStage{
		Name:        name,
		Description: description,
		AddReqIdx:   addReqIdx,
		Owner:       owner,
		Docs:        []TradeStageDoc{},
		DelReqs:     []ApproveReq{},
		CloseReqs:   []ApproveReq{},
	}
}
