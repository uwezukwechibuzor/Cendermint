package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type ConsensusService struct {
	DB *sql.DB
}

func (css *ConsensusService) GetConsensusDump(cfg config.Config, rpc *types.RPCData) {
	var cs types.ConsensusState
	var vSetsResult map[string][]string = make(map[string][]string)

	res, err := utils.HttpQuery(constants.RPCAddr + "/dump_consensus_state")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &cs)

	conspubMonikerMap := getConspubMonikerMap(cfg, rpc)
	// cs.Result.Validatorset.Validators is already sorted based on voting power
	for index, validator := range cs.Result.Validatorset.Validators {
		var prevote, precommit string

		// get moniker
		validator.Moniker = conspubMonikerMap[validator.ConsPubKey.Key]
		if cs.Result.RoundState.Votes[0].Prevotes[index] != "nil-Vote" {
			prevote = "✅"
		} else {
			prevote = "❌"
		}

		if cs.Result.RoundState.Votes[0].Precommits[index] != "nil-Vote" {
			precommit = "✅"
		} else {
			precommit = "❌"
		}

		// populate the map => [ConsAddr][]string{ConsPubKey, VotingPower, ProposerPriority, prevote, precommit, moniker}
		vSetsResult[validator.ConsAddr] = []string{validator.ConsPubKey.Key, validator.VotingPower, validator.ProposerPriority, prevote, precommit, validator.Moniker}
	}

	rpc.ConsensusState = cs
	rpc.Validatorsets = vSetsResult

	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Consensus", "height("+rpc.ConsensusState.Result.Height+") "+"round("+strconv.FormatInt(rpc.ConsensusState.Result.Round, 10)+") "+"step("+strconv.FormatInt(rpc.ConsensusState.Result.Step, 10)+")"))
	prevoteParsed := utils.ParseConsensusOutput(rpc.ConsensusState.Result.Votes[0].PrevotesBitArray, "\\= (.*)", 1)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Prevote bit array", fmt.Sprintf("%.2f", prevoteParsed)))
	precommitParsed := utils.ParseConsensusOutput(rpc.ConsensusState.Result.Votes[0].PrecommitsBitArray, "\\= (.*)", 1)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Precommit bit array", fmt.Sprintf("%.2f", precommitParsed)))
	zap.L().Info("", zap.Bool("Success", true), zap.String("# of validators from RPC", fmt.Sprint(len(rpc.Validatorsets))))
}

func getConspubMonikerMap(cfg config.Config, rpc *types.RPCData) map[string]string {
	var v types.RpcValidators
	var vResult map[string]string = make(map[string]string)

	route := rest.GetValidatorsRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route + "?status=BOND_STATUS_BONDED&pagination.limit=300")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &v)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	for _, validator := range v.Validators {
		// populate the map => [conspub] -> (moniker)
		vResult[validator.ConsPubKey.Key] = validator.Description.Moniker
	}
	return vResult
}
