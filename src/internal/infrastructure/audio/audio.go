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
		shouldBeMuted: false,
		isForceMuted:  false,
	}

	//初期シーンを設定する必要があるかどうかは後で考える
	audioClient.SetNormalScene()

	return audioClient, nil
}

func (self *AudioClient) SetMute(state bool) error {
	if self.isForceMuted == true {
		return fmt.Errorf("cannot change mute state: force muted is active")
	}

	if state == true {
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
	_, sceneUUID, err := self.GetCurrentScene()

	// res, err := self.obsClient.Inputs.GetInputMute(inputs.NewGetInputMuteParams().WithInputUuid(self.inputUuid))
	if err != nil {
		return entities.MuteState{}, err
	}

	return entities.MuteState{IsMuted: sceneUUID == self.scenes.Muted}, nil
}

func (self *AudioClient) SetForceMute(state bool) error {
	self.isForceMuted = state
	return nil
}

func (self *AudioClient) SetShouldBeMuted(state bool) error {
	self.shouldBeMuted = state
	return nil
}
