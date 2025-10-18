package audio

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type AudioClient struct {
	obsClient *goobs.Client
	scenes    Scenes
	// normalSceneUuid string
	// mutedSceneUuid  string
	// cmSceneUuid     string
	isConversion  bool
	shouldBeMuted bool
	isForceMuted  bool
}

// sceneのName or UUIDをまとめた型
type Scenes struct {
	Normal string
	Muted  string
	CM     string
}

func NewAudioClient(obsClient *goobs.Client, scenes Scenes) (*AudioClient, error) {
	// sceneUuid := ""
	// inputUuid := ""
	// sceneItemId := 0

	// sceneNameからsceneUuidを取得する。取得できなければ、そのようなsceneNameのSceneが存在しないと判断してエラー
	// sceneNameはobsから容易に変更可能なので安全性のためにUUIDを用いる
	resolve := func(name string) (string, error) {
		uuid, err := utils.FindSceneByName(obsClient, name)
		if err != nil {
			return "", fmt.Errorf("failed to find scene named %s: %w", name, err)
		}
		return uuid, nil
	}

	normalUUID, err := resolve(scenes.Normal)
	if err != nil {
		return nil, err
	}
	mutedUUID, err := resolve(scenes.Muted)
	if err != nil {
		return nil, err
	}
	cmUUID, err := resolve(scenes.CM)
	if err != nil {
		return nil, err
	}

	sceneUUIDs := Scenes{
		Normal: normalUUID,
		Muted:  mutedUUID,
		CM:     cmUUID,
	}

	audioClient := &AudioClient{
		obsClient:     obsClient,
		scenes:        sceneUUIDs,
		isConversion:  false,
		shouldBeMuted: false,
		isForceMuted:  false,
	}

	//初期シーンを設定する必要があるかどうかは後で考える
	audioClient.SetNormalScene()

	return audioClient, nil
}

func (self *AudioClient) SetMute(state bool) error {
	if self.isForceMuted && !state {
		return fmt.Errorf("cannot change mute state: force muted is active")
	}

	if state {
		err := self.SetMutedScene()
		if err != nil {
			return err
		}
	} else {
		err := self.SetNormalScene()
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *AudioClient) Mute() error {
	err := self.SetMute(true)
	return err
}

func (self *AudioClient) UnMute() error {
	err := self.SetMute(false)
	return err
}

func (self *AudioClient) GetMute() (entities.MuteState, error) {

	// sceneName, sceneUUID, err := GetCurrentScene(self.obsClient)
	sceneUUID, err := self.GetCurrentScene()

	if err != nil {
		return entities.MuteState{}, err
	}

	return entities.MuteState{IsMuted: sceneUUID == self.scenes.Muted}, nil
}

func (self *AudioClient) SetForceMute(state bool) error {
	currentSceneUuid, err := self.GetCurrentScene()
	if err != nil {
		return err
	}

	if currentSceneUuid == self.scenes.CM {
		return fmt.Errorf("cannot change force_mute state: it's CM scene")
	}

	//解除は慎重に行う
	if !state {
		if self.shouldBeMuted {
			return fmt.Errorf("cannot change force_mute state: There is music playing that needs to be muted")
		}

		//isForceMutedを先に書き換えないとと変更不可になってしまうため後でシーン変更
		self.isForceMuted = state
		err = self.SetMute(state)
		if err != nil {
			return err
		}
	} else {
		//isForceMutedを先に書き換えると変更不可になってしまうため先にシーン変更
		err = self.SetMute(state)
		if err != nil {
			return err
		}
		self.isForceMuted = state
	}

	//isForceMutedを先に書き換えると変更不可になってしまうため先にシーン変更
	err = self.SetMute(state)
	if err != nil {
		return err
	}
	self.isForceMuted = state

	return nil
}

func (self *AudioClient) SetShouldBeMuted(state bool) error {
	self.shouldBeMuted = state

	err := self.SetMute(state)
	return err

}

func (self *AudioClient) SetIsConversion(state bool) error {
	self.isConversion = state
	// conversionでは基本的に音声ありだが、force_mute中は音声なしになる
	if state {
		err := self.SetMute(!state)
		return err
	}
	return nil
}
