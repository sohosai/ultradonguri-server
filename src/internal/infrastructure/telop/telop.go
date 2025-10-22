package telop

import (
	"encoding/json"
	"log"
	"os"

	"github.com/samber/mo"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type Telop = mo.Either3[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop]

type TelopManager struct {
	telop      Telop
	backupPath string
}

type Backup struct {
	Type entities.TelopType `json:"type"`
	Data json.RawMessage    `json:"data"`
}

func NewTelopManager(backupPath string) *TelopManager {
	initialConversionTelop := entities.ConversionPost{}
	telop := mo.NewEither3Arg2[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](initialConversionTelop)

	return &TelopManager{
		telop:      telop,
		backupPath: backupPath,
	}
}

func (self *TelopManager) saveToFile() error {
	save_backup := func(type_ entities.TelopType, data json.RawMessage) error {
		backup := Backup{
			Type: type_,
			Data: data,
		}

		backup_json, err := json.Marshal(backup)
		if err != nil {
			return err
		}

		return os.WriteFile(self.backupPath, backup_json, 0666)
	}

	return match(self.telop,
		func(p entities.PerformancePost) error {
			data, err := json.Marshal(p)
			if err != nil {
				return err
			}

			return save_backup(entities.TelopTypePerformance, data)
		},
		func(c entities.ConversionPost) error {
			data, err := json.Marshal(c)
			if err != nil {
				return err
			}

			return save_backup(entities.TelopTypeConversion, data)
		},
		func(e entities.EmptyTelop) error {
			data, err := json.Marshal(e)
			if err != nil {
				return err
			}

			return save_backup(entities.TelopTypeEmpty, data)
		})
}

func (self *TelopManager) SetPerformanceTelop(performance entities.Performance) {
	defer self.saveToFile()

	performancePost := entities.PerformancePost{Performance: performance}

	self.telop = mo.NewEither3Arg1[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](performancePost)

	log.Printf("Telop changed: %v", performancePost)
}

func (self *TelopManager) SetMusicTelop(music entities.Music) {
	defer self.saveToFile()

	var performancePost entities.PerformancePost

	self.telop = match(self.telop, func(prevPerformance entities.PerformancePost) Telop {
		prevPerformance.Music = music
		performancePost = prevPerformance

		return mo.NewEither3Arg1[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](performancePost)
	}, func(_ entities.ConversionPost) Telop {
		performancePost = entities.PerformancePost{Music: music}

		return mo.NewEither3Arg1[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](performancePost)
	}, func(_ entities.EmptyTelop) Telop {
		performancePost = entities.PerformancePost{Music: music}

		return mo.NewEither3Arg1[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](performancePost)
	})

	self.telop = mo.NewEither3Arg1[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](performancePost)

	log.Printf("Telop changed: %v", performancePost)
}

func (self *TelopManager) SetConversionTelop(conversion entities.ConversionPost) {
	defer self.saveToFile()

	self.telop = mo.NewEither3Arg2[entities.PerformancePost, entities.ConversionPost, entities.EmptyTelop](conversion)

	log.Printf("Telop changed: %v", conversion)
}

func (self *TelopManager) GetCurrentTelopMessage() utils.Option[entities.TelopMessage] {
	return match(self.telop,
		func(performance entities.PerformancePost) utils.Option[entities.TelopMessage] {
			return utils.Some(entities.TelopMessage{
				Type:            entities.TelopTypePerformance,
				PerformanceData: &performance, // 参照だがヒープにコピーされた別のデータになる
				ConversionData:  nil,
			})
		},
		func(conversion entities.ConversionPost) utils.Option[entities.TelopMessage] {
			return utils.Some(entities.TelopMessage{
				Type:            entities.TelopTypeConversion,
				PerformanceData: nil,
				ConversionData:  &conversion, // 参照だがヒープにコピーされた別のデータになる
			})
		},
		func(_ entities.EmptyTelop) utils.Option[entities.TelopMessage] {
			return utils.None[entities.TelopMessage]()
		})
}

func (self *TelopManager) IsConversion() bool {
	return match(self.telop,
		func(_ entities.PerformancePost) bool {
			return false
		},
		func(_ entities.ConversionPost) bool {
			return true
		},
		func(_ entities.EmptyTelop) bool {
			return false
		})
}

func (self *TelopManager) ShouldBeMuted() bool {
	return match(self.telop,
		func(performance entities.PerformancePost) bool {
			return performance.Music.ShouldBeMuted
		},
		func(_ entities.ConversionPost) bool {
			return false
		},
		func(_ entities.EmptyTelop) bool {
			return false
		})
}

func match[T any](telop Telop,
	onPerformance func(entities.PerformancePost) T,
	onConversion func(entities.ConversionPost) T,
	onEmpty func(entities.EmptyTelop) T,
) T {
	var ret T

	_ = telop.Match(
		func(p entities.PerformancePost) Telop { ret = onPerformance(p); return telop },
		func(c entities.ConversionPost) Telop { ret = onConversion(c); return telop },
		func(e entities.EmptyTelop) Telop { ret = onEmpty(e); return telop },
	)

	return ret
}
