package usecases

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	notiRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/notification/repositories"
	nhRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
	retirementRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
	userRepo "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/google/uuid"
)

type AssetUseCase interface {
	CreateAsset(asset entities.Asset) (*entities.Asset, error)
	GetAssetByID(id string) (*entities.Asset, error)
	GetAssetByUserID(userID string) ([]entities.Asset, error)
	UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error)
	DeleteAssetByID(id string, userID string, transfers []entities.TransferRequest) error
}

type AssetUseCaseImpl struct {
	assetrepo      repositories.AssetRepository
	userrepo       userRepo.UserRepository
	nhrepo         nhRepo.NhRepository
	retirementrepo retirementRepo.RetirementRepository
	notirepo       notiRepo.NotiRepository
}

func NewAssetUseCase(assetrepo repositories.AssetRepository, userrepo userRepo.UserRepository, nhrepo nhRepo.NhRepository, retirementrepo retirementRepo.RetirementRepository, notirepo notiRepo.NotiRepository) *AssetUseCaseImpl {
	return &AssetUseCaseImpl{
		assetrepo:      assetrepo,
		userrepo:       userrepo,
		nhrepo:         nhrepo,
		retirementrepo: retirementrepo,
		notirepo:       notirepo,
	}
}

func (u *AssetUseCaseImpl) CreateNotification(userID, assetName string) error {
	notification := &entities.Notification{
		ID:        fmt.Sprintf("notif-%d-%s", time.Now().UnixNano(), assetName),
		UserID:    userID,
		Message:   fmt.Sprintf("ทรัพย์สิน '%s' ถูกหยุดพักชั่วคราวเนื่องจากหมดเวลา", assetName),
		CreatedAt: time.Now(),
	}
	return u.notirepo.CreateNotification(notification)
}

func (u *AssetUseCaseImpl) UpdateAssetStatus(asset *entities.Asset, currentYear int) error {
	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return err
	}

	if endYear <= currentYear {
		asset.Status = "Paused"
		asset.LastCalculatedMonth = 0
		asset.MonthlyExpenses = 0
		return u.CreateNotification(asset.UserID, asset.Name)
	}

	if asset.Status == "Paused" {
		asset.LastCalculatedMonth = 0
		asset.MonthlyExpenses = 0
		return nil
	}

	if asset.TotalCost <= asset.CurrentMoney {
		asset.Status = "Completed"
	} else {
		asset.Status = "In_Progress"
	}
	return nil
}

func (u *AssetUseCaseImpl) CreateAsset(asset entities.Asset) (*entities.Asset, error) {
	if asset.TotalCost <= 0 {
		return nil, errors.New("totalcost must be greater than zero")
	}

	endYear, err := strconv.Atoi(asset.EndYear)
	if err != nil {
		return nil, err
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	if endYear <= currentYear {
		return nil, errors.New("end year must be greater than or equal to current year")
	}

	id, err := u.assetrepo.GetAssetNextID()
	if err != nil {
		return nil, err
	}

	asset.ID = id
	asset.Status = "In_Progress"
	monthlyExpenses := utils.CalculateMonthlyExpenses(&asset, currentYear, currentMonth)
	asset.MonthlyExpenses = monthlyExpenses
	asset.LastCalculatedMonth = currentMonth

	return u.assetrepo.CreateAsset(&asset)
}

func (u *AssetUseCaseImpl) GetAssetByID(id string) (*entities.Asset, error) {
	return u.assetrepo.GetAssetByID(id)
}

func (u *AssetUseCaseImpl) GetAssetByUserID(userID string) ([]entities.Asset, error) {
	assets, err := u.assetrepo.GetAssetByUserID(userID)
	if err != nil {
		return nil, err
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	for i := range assets {
		if err := u.UpdateAssetStatus(&assets[i], currentYear); err != nil {
			return nil, err
		}

		if assets[i].LastCalculatedMonth != currentMonth {
			assets[i].MonthlyExpenses = utils.CalculateMonthlyExpenses(&assets[i], currentYear, currentMonth)
			assets[i].LastCalculatedMonth = currentMonth
			if _, err = u.assetrepo.UpdateAssetByID(&assets[i]); err != nil {
				return nil, err
			}
		}
	}

	return assets, nil
}

func (u *AssetUseCaseImpl) UpdateAssetByID(id string, asset entities.Asset) (*entities.Asset, error) {
	existingAsset, err := u.assetrepo.GetAssetByID(id)
	if err != nil {
		return nil, err
	}

	if asset.TotalCost <= 0 {
		return nil, errors.New("totalcost must be greater than zero")
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	totalCostChanged := existingAsset.TotalCost != asset.TotalCost
	existingAsset.TotalCost = asset.TotalCost
	existingAsset.Name = asset.Name
	existingAsset.Type = asset.Type
	existingAsset.EndYear = asset.EndYear
	existingAsset.Status = asset.Status
	existingAsset.LastCalculatedMonth = currentMonth
	if asset.Status == "Paused" {
		existingAsset.LastCalculatedMonth = 0
		existingAsset.MonthlyExpenses = 0
	} else {
		if err := u.UpdateAssetStatus(existingAsset, currentYear); err != nil {
			return nil, err
		}

		if totalCostChanged || existingAsset.LastCalculatedMonth != currentMonth {
			existingAsset.LastCalculatedMonth = currentMonth
			existingAsset.MonthlyExpenses = utils.CalculateMonthlyExpenses(existingAsset, currentYear, currentMonth)
		}
	}

	return u.assetrepo.UpdateAssetByID(existingAsset)
}

func (u *AssetUseCaseImpl) DeleteAssetByID(id string, userID string, transfers []entities.TransferRequest) error {
	asset, err := u.assetrepo.GetAssetByID(id)
	if err != nil {
		return err
	}

	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	totalTransfer := 0.0
	for _, transfer := range transfers {
		totalTransfer += transfer.Amount
	}

	if totalTransfer > asset.CurrentMoney {
		return errors.New("transfer amount exceeds asset's current money")
	}

	for _, transfer := range transfers {
		switch transfer.Type {
		case "asset":
			selectedItem, err := u.assetrepo.FindAssetByNameandUserID(transfer.Name, userID)
			if err != nil {
				return err
			}

			if selectedItem.Status == "In_Progress" {
				selectedItem.CurrentMoney += transfer.Amount
				if selectedItem.CurrentMoney >= selectedItem.TotalCost {
					selectedItem.Status = "Completed"
					selectedItem.MonthlyExpenses = 0
					selectedItem.LastCalculatedMonth = 0
					notification := &entities.Notification{
						ID:        uuid.New().String(),
						UserID:    user.ID,
						Message:   fmt.Sprintf("สุดยอดมาก สินทรัพย์ : '%s' ได้เสร็จสิ้นแล้ว", selectedItem.Name),
						CreatedAt: time.Now(),
					}

					_ = u.notirepo.CreateNotification(notification)
				}

				_, err = u.assetrepo.UpdateAssetByID(selectedItem)
				if err != nil {
					return err
				}
			} else {
				return errors.New("cannot update completed or paused asset")
			}

		case "house":
			if user.House.NursingHouseID != "00001" || user.House.Status != "Completed" {
				user.House.CurrentMoney += transfer.Amount
				requiredMoney := (user.RetirementPlan.ExpectLifespan - user.RetirementPlan.RetirementAge) * 12 * user.House.NursingHouse.Price
				if user.House.CurrentMoney >= float64(requiredMoney) {
					user.House.Status = "Completed"
					user.House.MonthlyExpenses = 0
					user.House.LastCalculatedMonth = 0
					notification := &entities.Notification{
						ID:        uuid.New().String(),
						UserID:    user.ID,
						Message:   fmt.Sprintf("สุดยอดมาก บ้านพัก : '%s' ได้เสร็จสิ้นแล้ว", user.House.NursingHouse.Name),
						CreatedAt: time.Now(),
					}

					_ = u.notirepo.CreateNotification(notification)
				}
			} else {
				return errors.New("cannot update completed nursing house")
			}

		case "retirementplan":
			user.RetirementPlan.CurrentSavings += transfer.Amount
			allMoney := user.RetirementPlan.CurrentSavings + user.RetirementPlan.CurrentTotalInvestment
			if allMoney >= user.RetirementPlan.LastRequiredFunds {
				user.RetirementPlan.Status = "Completed"
				user.RetirementPlan.LastMonthlyExpenses = 0
				user.RetirementPlan.LastMonthlyExpenses = 0
				notification := &entities.Notification{
					ID:        uuid.New().String(),
					UserID:    user.ID,
					Message:   fmt.Sprintf("สุดยอดมาก แผนเกษียณ : '%s' ของคุณได้ถึงเป้าแล้ว", user.RetirementPlan.PlanName),
					CreatedAt: time.Now(),
				}

				_ = u.notirepo.CreateNotification(notification)
			}

		default:
			continue
		}
	}
	return u.assetrepo.DeleteAssetByID(asset.ID)
}
