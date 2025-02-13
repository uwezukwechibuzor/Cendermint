package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rs RestServices) IndexValidator(consHexAddr, moniker string) {
	validator, err := rs.ValidatorService.Index(consHexAddr, moniker)
	if err != nil {
		zap.L().Error("Error indexing validator", zap.Error(err))
		return
	} else {
		zap.L().Info("Validator successfully indexed", zap.String("ConsHexAddress", validator.ConsHexAddress))
	}
}

func (rs RestServices) IndexAbsentValidator(height int, consHexAddr string) {
	absentValidator, err := rs.AbsentValidatorService.Index(height, consHexAddr)
	if err != nil {
		zap.L().Error("Error indexing abscent validator", zap.Error(err))
		return
	} else {
		zap.L().Info("Absent validator successfully indexed", zap.String("ConsHexAddress", absentValidator.ConsHexAddress))
	}
}

func (rs RestServices) GetValidatorInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) []string {
	return rs.ValidatorService.GetValidatorInfo(cfg, currentBlockHeight, rd)
}
