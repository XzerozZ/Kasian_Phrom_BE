package utils

import (
	"errors"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"gorm.io/gorm"
)

func DistributeSavingMoney(amount float64, assets []entities.Asset, house entities.SelectedHouse, plan entities.RetirementPlan, db *gorm.DB) error {
	count := len(assets)
	if house.NursingHouseID != "00001" {
		count++
	}

	count++
	if count == 0 {
		return errors.New("no assets or accounts to distribute money to")
	}

	portion := amount / float64(count)

	for i := range assets {
		(assets)[i].CurrentMoney += portion
		if err := db.Save(&(assets)[i]).Error; err != nil {
			return err
		}
	}

	if house.NursingHouseID != "00001" {
		house.CurrentMoney += portion
		if err := db.Save(house).Error; err != nil {
			return err
		}
	}

	plan.CurrentSavings += portion
	if err := db.Save(plan).Error; err != nil {
		return err
	}

	return nil
}

func WithdrawSavingMoney(amount float64, assets []entities.Asset, house entities.SelectedHouse, plan entities.RetirementPlan, db *gorm.DB) error {
	count := len(assets)
	if house.NursingHouseID != "00001" {
		count++
	}

	count++
	if count == 0 {
		return errors.New("no assets or accounts to withdraw money from")
	}

	portion := amount / float64(count)

	for i := range assets {
		if (assets)[i].CurrentMoney < portion {
			return errors.New("insufficient funds in asset")
		}

		(assets)[i].CurrentMoney -= portion
		if err := db.Save(&(assets)[i]).Error; err != nil {
			return err
		}
	}

	if house.NursingHouseID != "00001" {
		if house.CurrentMoney < portion {
			return errors.New("insufficient funds in selected house")
		}

		house.CurrentMoney -= portion
		if err := db.Save(house).Error; err != nil {
			return err
		}
	}

	if plan.CurrentSavings < portion {
		return errors.New("insufficient funds in savings")
	}

	plan.CurrentSavings -= portion
	if err := db.Save(plan).Error; err != nil {
		return err
	}

	return nil
}
