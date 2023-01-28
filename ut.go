package dst

import (
	"errors"

	"gorm.io/gorm"
)

type Hdl func() error

// 依赖于hdl会自动处理掉where命中数据,最终where里不包含命中记录
func CleanLoop(query *gorm.DB, dest interface{}, batchSize int, hdl Hdl, thresholds ...int) error {
	threshold := 10000
	if len(thresholds) > 0 && thresholds[0] > 1 {
		threshold = thresholds[0]
	}
	count := 0
	for {
		if count > threshold {
			return errors.New(`dead loop`)
		}
		count++
		result := query.Limit(batchSize).Find(dest)
		if result.Error != nil {
			return result.Error
		}
		if int(result.RowsAffected) == 0 {
			break
		}
		if err := hdl(); err != nil {
			return err
		}
	}
	return nil
}
