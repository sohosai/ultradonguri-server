package scene

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type SceneManager struct {
	obsClient *goobs.Client
	scenes    Scenes
	sceneType SceneType // ファイルバックアップのために今設定されているSceneTypeを保存しておく。setSceneというメソッドだけが触る
	// isForceMutedがtrueなときに必ずシーンがMutedになっているわけではなく、
	// 単にSetMute実行時に内部で利用するためのフラグとして考えたほうが良い気がする
	isForceMutedFlag bool
	backupPath       string
}

// sceneのName or UUIDをまとめた型
type Scenes struct {
	Normal string
	Muted  string
	CM     string
}

type SceneNames struct {
	Normal string
	Muted  string
	CM     string
}

type SceneType = int

const (
	Normal SceneType = iota
	Muted
	CM
)

type Backup struct {
	SceneType        SceneType `json:"scene_type"`
	IsForceMutedFlag bool      `json:"is_force_muted"`
}

func NewSceneManager(obsClient *goobs.Client, sceneNames SceneNames, backupPath string) (*SceneManager, error) {
	return newSceneManager(obsClient, sceneNames, backupPath, Normal)
}

func RestoreSceneManager(obsClient *goobs.Client, sceneNames SceneNames, backupPath string) (*SceneManager, error) {
	backupRaw, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, err
	}

	savedInfo, err := utils.JsonStrictUnmarshal[Backup](backupRaw)
	if err != nil {
		return nil, err
	}

	sceneManager, err := newSceneManager(obsClient, sceneNames, backupPath, savedInfo.SceneType)
	if err != nil {
		return nil, err
	}

	sceneManager.isForceMutedFlag = savedInfo.IsForceMutedFlag

	return sceneManager, err
}

func newSceneManager(obsClient *goobs.Client, sceneNames SceneNames, backupPath string, initialScene SceneType) (*SceneManager, error) {
	// sceneNameからsceneUuidを取得する。取得できなければ、そのようなsceneNameのSceneが存在しないと判断してエラー
	// sceneNameはobsから容易に変更可能なので安全性のためにUUIDを用いる
	resolve := func(name string) (string, error) {
		uuid, err := utils.FindSceneByName(obsClient, name)
		if err != nil {
			return "", fmt.Errorf("failed to find scene named %s: %w", name, err)
		}
		return uuid, nil
	}

	normalUUID, err := resolve(sceneNames.Normal)
	if err != nil {
		return nil, err
	}
	mutedUUID, err := resolve(sceneNames.Muted)
	if err != nil {
		return nil, err
	}
	cmUUID, err := resolve(sceneNames.CM)
	if err != nil {
		return nil, err
	}

	sceneUUIDs := Scenes{
		Normal: normalUUID,
		Muted:  mutedUUID,
		CM:     cmUUID,
	}

	sceneManager := &SceneManager{
		obsClient:        obsClient,
		scenes:           sceneUUIDs,
		isForceMutedFlag: false,
		backupPath:       backupPath,
	}

	switch initialScene {
	case Normal:
		err = sceneManager.SetNormalScene()
		sceneManager.sceneType = Normal

	case Muted:
		err = sceneManager.SetMutedScene()
		sceneManager.sceneType = Muted

	case CM:
		err = sceneManager.SetCMScene()
		sceneManager.sceneType = CM

	default:
		err = sceneManager.SetNormalScene()
		sceneManager.sceneType = Normal
	}

	return sceneManager, err
}

// force_mute時のmute切り替えもこのメソッドが行う。
func (self *SceneManager) SetMute(state bool) error {
	if !state {
		// ミュートを解除する場合

		if self.isForceMutedFlag {
			return fmt.Errorf("cannot change mute state: force muted is active")
			// return nil
		}

		return self.SetNormalScene()
	}

	// ミュートする場合
	return self.SetMutedScene()
}

func (self *SceneManager) SetForceMuteFlag(state bool) {
	self.isForceMutedFlag = state

	self.saveToFile()
}

func (self *SceneManager) IsCm() (bool, error) {
	currentScene, err := self.GetCurrentScene()
	if err != nil {
		return false, err
	}

	return currentScene == self.scenes.CM, nil
}

func (self *SceneManager) saveToFile() error {
	info := Backup{
		SceneType:        self.sceneType,
		IsForceMutedFlag: self.isForceMutedFlag,
	}

	data, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(self.backupPath, data, 0o600)
}
