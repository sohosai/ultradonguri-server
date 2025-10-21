package scene

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type SceneManager struct {
	obsClient *goobs.Client
	scenes    Scenes
	// isForceMutedがtrueなときに必ずシーンがMutedになっているわけではなく、
	// 単にSetMute実行時に内部で利用するためのフラグとして考えたほうが良い気がする
	isForceMutedFlag bool
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

func NewSceneManager(obsClient *goobs.Client, sceneNames SceneNames) (*SceneManager, error) {
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

	audioClient := &SceneManager{
		obsClient:        obsClient,
		scenes:           sceneUUIDs,
		isForceMutedFlag: false,
	}

	//初期シーンを設定する必要があるかどうかは後で考える
	audioClient.SetNormalScene()

	return audioClient, nil
}

// force_mute時のmute切り替えもこのメソッドが行う。
func (self *SceneManager) SetMute(state bool) error {
	if !state {
		// ミュートを解除する場合

		if self.isForceMutedFlag {
			return fmt.Errorf("cannot change mute state: force muted is active")
		}

		return self.SetNormalScene()
	}

	// ミュートする場合
	return self.SetMutedScene()
}

func (self *SceneManager) SetForceMuteFlag(state bool) {
	self.isForceMutedFlag = state
}

func (self *SceneManager) IsCm() (bool, error) {
	currentScene, err := self.GetCurrentScene()
	if err != nil {
		return false, err
	}

	return currentScene == self.scenes.CM, nil
}
