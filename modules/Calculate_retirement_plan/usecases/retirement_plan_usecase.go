package usecases

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/repositories"

	user_repo "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"math"

	"fmt"
	"github.com/gofiber/fiber/v2"
)

type RetUseCase interface {
	CreateRet(userID string, ctx *fiber.Ctx) (*entities.RetirementPlan, error)
	GetRetByID(id string) (*entities.RetirementPlan, error)
	GetRetNextID() (string, error)
}

type RetUseCaseImpl struct {
	retrepo  repositories.RetRepository
	finrepo  repositories.FinRepository
	assrepo  repositories.AssRepository
	userrepo user_repo.UserRepository
}

func NewRetUseCase(
	retrepo repositories.RetRepository,
	finrepo repositories.FinRepository,
	assrepo repositories.AssRepository,
	userrepo user_repo.UserRepository,
) *RetUseCaseImpl {
	return &RetUseCaseImpl{
		retrepo:  retrepo,
		finrepo:  finrepo,
		assrepo:  assrepo,
		userrepo: userrepo,
	}
}

// create and calculate
func (u *RetUseCaseImpl) CreateRet(userID string, ctx *fiber.Ctx) (*entities.RetirementPlan, error) {
	fmt.Println("เข้ามาใน usecase")
	id, err := u.retrepo.GetRetNextID()
	if err != nil {
		return nil, err
	}
    fmt.Println("ผ่าน next id")
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	username := user.Username
    fmt.Println("ได้ username")
    fmt.Println(username)

	financial, err := u.finrepo.GetFinByUserID(userID)
	if err != nil {
		return nil, err
	}

    fmt.Println("ได้ financial")
    fmt.Println(financial.Age)

	assets, err := u.assrepo.GetAssByUsername(username)
	if err != nil {
		return nil, err
	}

    fmt.Println("ได้ array Asset")
    fmt.Println(assets)

	// คำนวณค่าใช้จ่ายทั้งหมดหลังเกษียณ (ปรับเงินเฟ้อ)
	totalAssetCost := 0.0
	for _, asset := range assets {
		totalAssetCost += asset.TotalMoney
	}

	// ปรับค่าใช้จ่ายหลังเกษียณให้รวมเงินเฟ้อแบบทบต้น (Future Value Formula)
	yearsAfterRetirement := float64(financial.ExpLifespan - financial.RetirementAge)
	inflationFactor := math.Pow(1+financial.ExpectedInflation, yearsAfterRetirement)

	totalRetirementExpense := (financial.ExpectedMonthlyExpenses * 12 * yearsAfterRetirement) * inflationFactor

	yearsToRetirement := float64(financial.RetirementAge - financial.Age)
	assetInflationFactor := math.Pow(1+financial.ExpectedInflation, yearsToRetirement)
	adjustedTotalAssetCost := totalAssetCost * assetInflationFactor

	totalRetirementExpense += adjustedTotalAssetCost // รวม Asset เข้าไปด้วย รวมเงินเฟ้อราคาบ้านแพงขึ้นแล้ว

	// คำนวณจำนวนปีและเดือนก่อนเกษียณ

	monthsToRetirement := yearsToRetirement * 12

	// ใช้ Future Value of Annuity Formula (FVIFA) เพื่อคำนวณเงินที่ต้องออมรายเดือนโดยคำนึงถึงดอกเบี้ย
	r := financial.AnnualSavingsReturn / 12              // ดอกเบี้ยรายเดือน
	fvifa := (math.Pow(1+r, monthsToRetirement) - 1) / r // FVIFA Formula

	requiredMonthlySavings := totalRetirementExpense / fvifa

	// เซ็ตค่าใน RetirementPlan
    fmt.Println("ก่อน pointer ของ retirementPlan")
	// var retirementPlan *entities.RetirementPlan // บัคบันทัดนี้
    retirementPlan := &entities.RetirementPlan{} // สร้าง instance ใหม่

	retirementPlan.TotalRetirementExpense = totalRetirementExpense
	retirementPlan.RequiredMonthlySavings = requiredMonthlySavings
	retirementPlan.TotalAssetCost = totalAssetCost

	retirementPlan.ID = id

    fmt.Println("หลัง pointer ของ retirementPlan")

	// var createdRet entities.RetirementPlan
	retirementPlan, err = u.retrepo.CreateRet(retirementPlan)
	if err != nil {
		return nil, err
	}

	return retirementPlan, nil
}

func (u *RetUseCaseImpl) GetRetByID(id string) (*entities.RetirementPlan, error) {
	return u.retrepo.GetRetByID(id)
}

func (u *RetUseCaseImpl) GetRetNextID() (string, error) {
	return u.retrepo.GetRetNextID()
}
