package utils_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		birthDate   string
		expectedAge int
		expectError bool
	}{
		{
			name:        "ปกติ - อายุ 30 ปี",
			birthDate:   now.AddDate(-30, 0, 0).Format("02-01-2006"),
			expectedAge: 30,
			expectError: false,
		},
		{
			name:        "รูปแบบวันที่ไม่ถูกต้อง",
			birthDate:   "2000-01-02",
			expectedAge: 0,
			expectError: true,
		},
		{
			name:        "ยังไม่ถึงวันเกิดปีนี้",
			birthDate:   now.AddDate(-30, 1, 0).Format("02-01-2006"),
			expectedAge: 29,
			expectError: false,
		},
		{
			name:        "เกิดวันนี้",
			birthDate:   now.Format("02-01-2006"),
			expectedAge: 0,
			expectError: false,
		},
		{
			name:        "วันเกิดในอนาคต",
			birthDate:   now.AddDate(1, 0, 0).Format("02-01-2006"),
			expectedAge: -1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age, err := utils.CalculateAge(tt.birthDate)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAge, age)
			}
		})
	}
}

func TestCalculateRetirementFunds(t *testing.T) {
	tests := []struct {
		name        string
		plan        *entities.RetirementPlan
		age         int
		expectError bool
		checkFunds  func(t *testing.T, funds float64)
	}{
		{
			name: "คำนวณปกติ",
			plan: &entities.RetirementPlan{
				ExpectedMonthlyExpenses: 30000,
				ExpectedInflation:       3,
				RetirementAge:           60,
				ExpectLifespan:          80,
			},
			age:         30,
			expectError: false,
			checkFunds: func(t *testing.T, funds float64) {
				assert.Greater(t, funds, float64(0))
			},
		},
		{
			name: "อายุเท่ากับอายุเกษียณ",
			plan: &entities.RetirementPlan{
				RetirementAge:  60,
				ExpectLifespan: 80,
			},
			age:         60,
			expectError: true,
			checkFunds:  func(t *testing.T, funds float64) {},
		},
		{
			name: "อายุคาดหมายเท่ากับอายุเกษียณ",
			plan: &entities.RetirementPlan{
				RetirementAge:  60,
				ExpectLifespan: 60,
			},
			age:         30,
			expectError: true,
			checkFunds:  func(t *testing.T, funds float64) {},
		},
		{
			name: "ค่าใช้จ่ายรายเดือนเป็น 0",
			plan: &entities.RetirementPlan{
				ExpectedMonthlyExpenses: 0,
				ExpectedInflation:       3,
				RetirementAge:           60,
				ExpectLifespan:          80,
			},
			age:         30,
			expectError: false,
			checkFunds: func(t *testing.T, funds float64) {
				assert.Equal(t, float64(0), funds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funds, err := utils.CalculateRetirementFunds(tt.plan, tt.age)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.checkFunds(t, funds)
			}
		})
	}
}

func TestCalculateMonthlyExpenses(t *testing.T) {
	currentYear := time.Now().Year()
	currentMonth := int(time.Now().Month())

	tests := []struct {
		name             string
		asset            *entities.Asset
		currentYear      int
		currentMonth     int
		expectedExpenses float64
	}{
		{
			name: "สถานะ Completed",
			asset: &entities.Asset{
				Status:       "Completed",
				TotalCost:    60000,
				CurrentMoney: 60000,
			},
			expectedExpenses: 0,
		},
		{
			name: "กำลังดำเนินการปกติ",
			asset: &entities.Asset{
				Status:       "In_Progress",
				TotalCost:    66000,
				CurrentMoney: 11000,
				EndYear:      "2026",
			},
			expectedExpenses: 5000,
		},
		{
			name: "เงินปัจจุบันเท่ากับค่าใช้จ่ายทั้งหมด",
			asset: &entities.Asset{
				Status:       "In_Progress",
				TotalCost:    60000,
				CurrentMoney: 60000,
				EndYear:      "2026",
			},
			expectedExpenses: 0,
		},
		{
			name: "ปีสิ้นสุดผิดรูปแบบ",
			asset: &entities.Asset{
				Status:       "In_Progress",
				TotalCost:    60000,
				CurrentMoney: 10000,
				EndYear:      "invalid",
			},
			expectedExpenses: 0,
		},
		{
			name: "ระยะเวลาสั้นกว่า 1 เดือน",
			asset: &entities.Asset{
				Status:       "In_Progress",
				TotalCost:    60000,
				CurrentMoney: 10000,
				EndYear:      strconv.Itoa(currentYear),
			},
			currentYear:      currentYear,
			currentMonth:     12,
			expectedExpenses: 0,
		},
		{
			name: "ปีแรกเดือนมกรา",
			asset: &entities.Asset{
				Status:       "In_Progress",
				TotalCost:    60000,
				CurrentMoney: 10000,
				EndYear:      strconv.Itoa(currentYear + 2),
			},
			currentYear:      currentYear,
			currentMonth:     1,
			expectedExpenses: 2174,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenses := utils.CalculateMonthlyExpenses(tt.asset, currentYear, currentMonth)
			assert.Equal(t, tt.expectedExpenses, expenses,
				"For asset with total cost %.2f and current money %.2f",
				tt.asset.TotalCost, tt.asset.CurrentMoney)
		})
	}
}

func TestCalculateAllAssetsMonthlyExpenses(t *testing.T) {
	tests := []struct {
		name             string
		user             *entities.User
		expectedExpenses float64
		expectError      bool
	}{
		{
			name: "มี assets หลายสถานะ",
			user: &entities.User{
				Assets: []entities.Asset{
					{
						Status:              "In_Progress",
						TotalCost:           60000,
						CurrentMoney:        10000,
						EndYear:             "2026",
						LastCalculatedMonth: int(time.Now().Month()),
						MonthlyExpenses:     1042,
					},
					{
						Status:       "Completed",
						TotalCost:    30000,
						CurrentMoney: 30000,
					},
				},
			},
			expectedExpenses: 1042,
			expectError:      false,
		},
		{
			name: "assets ทั้งหมดเสร็จสิ้น",
			user: &entities.User{
				Assets: []entities.Asset{
					{
						Status:       "Completed",
						TotalCost:    60000,
						CurrentMoney: 60000,
					},
					{
						Status:       "Completed",
						TotalCost:    30000,
						CurrentMoney: 30000,
					},
				},
			},
			expectedExpenses: 0,
			expectError:      false,
		},
		{
			name: "ไม่มี assets",
			user: &entities.User{
				Assets: []entities.Asset{},
			},
			expectedExpenses: 0,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenses, err := utils.CalculateAllAssetsMonthlyExpenses(tt.user)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExpenses, expenses)
			}
		})
	}
}

func TestCalculateAllAssetSavings(t *testing.T) {
	tests := []struct {
		name          string
		user          *entities.User
		method        string
		expectedTotal float64
	}{
		{
			name: "method All กับ assets หลายสถานะ",
			user: &entities.User{
				Assets: []entities.Asset{
					{
						Status:       "In_Progress",
						CurrentMoney: 10000,
						TotalCost:    20000,
					},
					{
						Status:       "Completed",
						CurrentMoney: 30000,
						TotalCost:    30000,
					},
				},
			},
			method:        "All",
			expectedTotal: 40000,
		},
		{
			name: "method Plan กับ assets หลายสถานะ",
			user: &entities.User{
				Assets: []entities.Asset{
					{
						Status:       "In_Progress",
						CurrentMoney: 10000,
						TotalCost:    20000,
					},
					{
						Status:       "Completed",
						CurrentMoney: 30000,
						TotalCost:    30000,
					},
				},
			},
			method:        "Plan",
			expectedTotal: 40000,
		},
		{
			name: "ไม่มี assets",
			user: &entities.User{
				Assets: []entities.Asset{},
			},
			method:        "All",
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := utils.CalculateAllAssetSavings(tt.user, tt.method)
			assert.Equal(t, tt.expectedTotal, total)
		})
	}
}
