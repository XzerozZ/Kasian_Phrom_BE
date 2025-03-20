package utils

import (
	"fmt"
	"time"

	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

var assetMessages = []string{
	"🎯 สำเร็จแล้ว! คุณสะสมเงินครบเป้าหมาย %s เก่งมาก!",
	"🌟 ยอดเยี่ยม! สินทรัพย์ %s ถึงเป้าหมายแล้ว ก้าวต่อไปรออยู่!",
	"💪 เยี่ยมจริงๆ! สินทรัพย์ %s สำเร็จแล้ว คุณทำได้ดีมาก!",
	"✨ สุดปัง! สินทรัพย์ %s เข้าเป้าแล้วจ้า #ชีวิตรุ่ง",
	"🔥 เลเวลอัพ! สินทรัพย์ %s คอมพลีทแล้ว เก่งมากนะ",
	"✨ ความสำเร็จ! เป้าหมาย %s บรรลุแล้ว — ก้าวสำคัญของคุณ",
}

var houseMessages = []string{
	"🏠 ชัยชนะ! แผนบ้านพักของคุณ บรรลุเป้าหมายแล้ว ภูมิใจในตัวคุณ!",
	"🏆 วิน! แผนบ้านพักของคุณ สำเร็จแล้ว #เงินทองต้องวางแผน",
	"🏡 เป้าหมายสำเร็จ! แผนบ้านพักของคุณ เสร็จสมบูรณ์ — อนาคตที่มั่นคงรออยู่",
}

var retirementPlanMessages = []string{
	"🏠 ชัยชนะ! แผนเกษียณ %s บรรลุเป้าหมายแล้ว ภูมิใจในตัวคุณ!",
	"🏆 วิน! แผนเกษียณ %s สำเร็จแล้ว #เงินทองต้องวางแผน",
	"🏡 เป้าหมายสำเร็จ! แผนเกษียณ %s เสร็จสมบูรณ์ — อนาคตที่มั่นคงรออยู่",
}

var loanMessages = []string{
	"🎉 หมดหนี้! คุณชำระ%sครบถ้วนแล้ว อิสรภาพทางการเงินใกล้เข้ามา!",
	"💸 ฟรีแล้ว! ปลดหนี้ %s เรียบร้อย อิสระทางการเงินมาแล้วจ้า",
	"🔓 ปลดล็อคสำเร็จ! หนี้ %s ชำระครบถ้วน — ก้าวสู่อิสรภาพทางการเงิน",
}

var alertLoanMessages = []string{
	"🔔 ยังมีรายการหนี้ชำระของเดือนก่อนที่ยังไม่เสร็จสมบูรณ์",
	"🌟 เพื่อประโยชน์ของคุณ: ขอแจ้งว่ายังมีรายการชำระเดือนที่แล้วรออยู่",
	"🌼 ขอเรียนเตือนด้วยความห่วงใย: การชำระเดือนที่ผ่านมายังรอการดำเนินการ",
	"📋 แจ้งเตือนสถานะ: การชำระหนี้ของเดือนที่แล้วยังรอการดำเนินการ",
	"🍀 เพื่อสุขภาพทางการเงินที่ดี: รายการชำระหนี้ประจำเดือนที่ผ่านมายังรออยู่",
	"⭐ เพื่อการเงินที่ราบรื่น: การชำระหนี้เดือนที่ผ่านมายังไม่เสร็จสมบูรณ์",
}

var alertAssetMessages = []string{
	"✨ เพื่อการวางแผนที่ดี: มีสินทรัพย์ของคุณอยู่ในโหมดรอดำเนินการ",
	"🔔 ข้อมูลสำคัญ: มีสินทรัพย์ต้องการการอัปเดต (ถึงวันที่ครบกำหนดแล้ว)",
	"📝 แจ้งให้ทราบ: สินทรัพย์ของคุณถูกพักการทำงานชั่วคราวเนื่องจากถึงวันที่กำหนด",
}

func SuccessNotification(itemType, userID, itemName, objectID string, balance float64) *entities.Notification {
	switch itemType {
	case "asset":
		notification := &entities.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Message:   fmt.Sprintf(assetMessages[rand.Intn(len(assetMessages))], itemName),
			Type:      itemType,
			ObjectID:  objectID,
			Balance:   balance,
			CreatedAt: time.Now(),
		}

		return notification
	case "house":
		notification := &entities.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Message:   houseMessages[rand.Intn(len(houseMessages))],
			Type:      itemType,
			ObjectID:  objectID,
			Balance:   balance,
			CreatedAt: time.Now(),
		}

		return notification
	case "retirementplan":
		notification := &entities.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Message:   fmt.Sprintf(retirementPlanMessages[rand.Intn(len(retirementPlanMessages))], itemName),
			Type:      itemType,
			ObjectID:  objectID,
			Balance:   balance,
			CreatedAt: time.Now(),
		}

		return notification
	case "loan":
		notification := &entities.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Message:   fmt.Sprintf(loanMessages[rand.Intn(len(loanMessages))], itemName),
			Type:      itemType,
			ObjectID:  objectID,
			Balance:   balance,
			CreatedAt: time.Now(),
		}

		return notification
	default:
		return nil
	}
}

func AlertNoti(itemType, userID, itemName, objectID string, balance float64) *entities.Notification {
	switch itemType {
	case "asset":
		notification := &entities.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Message:   alertAssetMessages[rand.Intn(len(alertAssetMessages))],
			Type:      itemType,
			ObjectID:  objectID,
			CreatedAt: time.Now(),
		}

		return notification
	case "loan":
		notification := &entities.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Message:   alertLoanMessages[rand.Intn(len(alertLoanMessages))],
			Type:      itemType,
			ObjectID:  objectID,
			Balance:   balance,
			CreatedAt: time.Now(),
		}

		return notification
	default:
		return nil
	}
}
