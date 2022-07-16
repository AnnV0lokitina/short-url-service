package service

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	"strconv"
	"strings"
)

const ipDelim = "."
const maskSubnet = "255.255.255.0"

func (s *Service) getIPSubnet(ipStr string) (string, error) {
	maskParts := strings.Split(maskSubnet, ipDelim)
	ipParts := strings.Split(ipStr, ipDelim)
	var maskedParts []string
	for i, ipPart := range ipParts {
		ipPartNum, err := strconv.Atoi(ipPart)
		if err != nil {
			return "", err
		}
		maskPartNum, err := strconv.Atoi(maskParts[i])
		if err != nil {
			return "", err
		}
		resultNum := ipPartNum & maskPartNum
		if resultNum > 0 {
			maskedParts = append(maskedParts, strconv.Itoa(resultNum))
		}
	}
	subnet := strings.Join(maskedParts, ipDelim)
	return subnet, nil
}

// GetStats Gets statistic of urls and users.
func (s *Service) GetStats(ctx context.Context, ipStr string) (entity.Stats, error) {
	if ipStr == "" {
		return entity.Stats{}, labelError.NewLabelError(labelError.TypeForbidden, errors.New("forbidden"))
	}
	subnet, err := s.getIPSubnet(ipStr)
	if err != nil {
		return entity.Stats{}, err
	}
	if subnet != s.trustedSubnet {
		return entity.Stats{}, labelError.NewLabelError(labelError.TypeForbidden, errors.New("forbidden"))
	}
	urls, users, err := s.repo.GetStats(ctx)
	if err != nil {
		return entity.Stats{}, err
	}
	return entity.Stats{
		Users: users,
		URLs:  urls,
	}, nil
}
