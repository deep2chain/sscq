package mint

const (
	BlkTime         = 5
	AvgDaysPerMonth = 30
	DayinSecond     = 24 * 3600
	AvgBlksPerMonth = AvgDaysPerMonth * DayinSecond / BlkTime

	ValidatorNumbers    = 2          // the number of validators
	ValidatorProvisions = float64(1) // 1 for each validator
	// Ignore validaotorProvisin because it is set as small as enough to be neglected.
	// if you want to set it bigger as your , you should care about this part again.
	// Including this part will decrease interoperability of the source code.
	// ValidatorTotalProvisions = ValidatorProvisions * ValidatorNumbers // 1 for each validator
	ValidatorTotalProvisions = 0

	// IssuerAmount = float64(1000000) // this is for test. 0 for production, 1000000 for test
	IssuerAmount = 0 // 0 for production

	FixedMineProvision  = float64(189000000)
	MineTotalProvisions = FixedMineProvision - ValidatorTotalProvisions - IssuerAmount // ~36,000,000 for 40 years

	// this is for export case,that's,this is activated if there exporting accounts exist.
	UserProvisions = float64(21000000) // if not, this should be set as zero

	CurrentProvisions          = UserProvisions + ValidatorTotalProvisions + IssuerAmount // ~60,000,000 at genesis
	CurrentProvisionsAsSatoshi = int64(CurrentProvisions * sscq2satoshi)                  // ~60,000,000 at genesis
	TotalLiquid                = MineTotalProvisions + CurrentProvisions                  // 96,000,000
	TotalLiquidAsSatoshi       = int64(TotalLiquid * sscq2satoshi)                        // 96,000,000 * 100,000,000

	sscq2satoshi = 100000000 // 1 sscq = 10 ** 8 satoshi

)
