package txvalidation

import "bitbucket.org/cerealia/apps/go-lib/model"

var path1 = model.TradeStageDocPath{
	StageIdx:    133,
	StageDocIdx: 242,
}

var path2 = model.TradeStageDocPath{
	StageIdx:     123,
	StageDocIdx:  456,
	StageDocHash: "abf12341523aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
}

var path3 = model.TradeStageDocPath{
	StageIdx:     987,
	StageDocIdx:  13,
	StageDocHash: "abf12341523aaaaaaccccccccccccccccccccccccccccccccccccccccccccccc",
}

var path4TextHash = model.TradeStageDocPath{
	StageIdx:     123,
	StageDocIdx:  456,
	StageDocHash: "hello world",
}

var user1 = &model.User{
	ID:              "test-id",
	DefaultWalletID: "my-wallet",
	StaticWallets: map[string]model.StaticWallet{
		"my-wallet": model.StaticWallet{
			PubKey: stageActionTXSigner,
		},
	},
}
