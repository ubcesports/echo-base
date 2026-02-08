package models

import (
	"fmt"
	"time"

	"github.com/ubcesports/echo-base/internal/utils"
)

type MembershipTier interface {
	GetName() string
	GetExpiryDate() (*time.Time, error)
	GetSessionDurationMs() int64
	HasDailyLimit() bool
	IsExpired(expiryDate *time.Time) (bool, error)
}

type TierConstructor func() MembershipTier

var tierRegistry = make(map[int]TierConstructor)

func RegisterTier(tierNumber int, constructor TierConstructor) {
	tierRegistry[tierNumber] = constructor
}

func NewMembershipTier(tierNumber int) (MembershipTier, error) {
	constructor, exists := tierRegistry[tierNumber]
	if !exists {
		return nil, fmt.Errorf("invalid tier number: %d (must be 0-3)", tierNumber)
	}
	return constructor(), nil
}

type NoMembership struct{}

func init() {
	RegisterTier(0, func() MembershipTier { return &NoMembership{} })
}

func (t *NoMembership) GetName() string {
	return "No Membership"
}

func (t *NoMembership) GetExpiryDate() (*time.Time, error) {
	return nil, nil
}

func (t *NoMembership) GetSessionDurationMs() int64 {
	return 0
}

func (t *NoMembership) HasDailyLimit() bool {
	return false
}

func (t *NoMembership) IsExpired(expiryDate *time.Time) (bool, error) {
	return false, nil
}

type Tier1 struct{}

func init() {
	RegisterTier(1, func() MembershipTier { return &Tier1{} })
}

func (t *Tier1) GetName() string {
	return "Tier 1"
}

func (t *Tier1) GetExpiryDate() (*time.Time, error) {
	return utils.GetNextMayFirst()
}

func (t *Tier1) GetSessionDurationMs() int64 {
	return 60 * 60 * 1000
}

func (t *Tier1) HasDailyLimit() bool {
	return true
}

func (t *Tier1) IsExpired(expiryDate *time.Time) (bool, error) {
	return utils.IsDateExpired(expiryDate)
}

type Tier2 struct{}

func init() {
	RegisterTier(2, func() MembershipTier { return &Tier2{} })
}

func (t *Tier2) GetName() string {
	return "Tier 2"
}

func (t *Tier2) GetExpiryDate() (*time.Time, error) {
	return utils.GetNextMayFirst()
}

func (t *Tier2) GetSessionDurationMs() int64 {
	return 2 * 60 * 60 * 1000
}

func (t *Tier2) HasDailyLimit() bool {
	return false
}

func (t *Tier2) IsExpired(expiryDate *time.Time) (bool, error) {
	return utils.IsDateExpired(expiryDate)
}

type PremierTier struct{}

func init() {
	RegisterTier(3, func() MembershipTier { return &PremierTier{} })
}

func (t *PremierTier) GetName() string {
	return "Premier"
}

func (t *PremierTier) GetExpiryDate() (*time.Time, error) {
	return utils.GetNextMayFirst()
}

func (t *PremierTier) GetSessionDurationMs() int64 {
	return 5 * 60 * 60 * 1000
}

func (t *PremierTier) HasDailyLimit() bool {
	return false
}

func (t *PremierTier) IsExpired(expiryDate *time.Time) (bool, error) {
	return utils.IsDateExpired(expiryDate)
}
