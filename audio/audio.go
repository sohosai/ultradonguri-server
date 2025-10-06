package audio

import (
	"fmt"

	"example.com/donguri-back/util"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
)

type AudioClient struct {
	obsClient   *goobs.Client
	sceneUuid   string
	inputUuid   string
	sceneItemId int
}

func NewAudioClient(obsClient *goobs.Client, sceneName string, inputName string) (*AudioClient, error) {
	// sceneUuid := ""
	// inputUuid := ""
	// sceneItemId := 0

	// sceneNameからsceneUuidを取得する。取得できなければ、そのようなsceneNameのSceneが存在しないと判断してエラー
	sceneUuid, err := util.FindSceneByName(obsClient, sceneName)
	if err != nil {
		return nil, fmt.Errorf("Failed to find scene named %s: %s", sceneName, err)
	}

	// inputNameからinputUuidを取得する。取得できなければ、そのようなinputNameのInputが存在しないと判断してエラー
	inputUuid, err := util.FindInputByName(obsClient, inputName)
	if err != nil {
		return nil, fmt.Errorf("Failed to find input named %s: %s", inputName, err)
	}

	// 最初からオーディオのInputはSceneItemになっているはずなので、それを確認し、なっていなかったらエラー
	// 同一のSceneに同一のInputが複数のSceneItemとして登録されうるが、この場合一番始めに見つかったものが返されるのだろうか。
	// 一旦一番始めに見つかったものを使う
	sceneItemId, err := obsClient.SceneItems.GetSceneItemId(
		sceneitems.NewGetSceneItemIdParams().
			WithSceneUuid(sceneUuid).
			WithSceneName(sceneName).
			WithSourceName(inputName))
	if err != nil {
		return nil, fmt.Errorf("Failed to find scene item id: %s", err)
	}

	return &AudioClient{
		obsClient:   obsClient,
		sceneUuid:   sceneUuid,
		inputUuid:   inputUuid,
		sceneItemId: sceneItemId.SceneItemId,
	}, nil
}

func (self *AudioClient) SetMute(state bool) error {
	_, err := self.obsClient.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().
			WithInputUuid(self.inputUuid).
			WithInputMuted(state))

	return err
}

func (self *AudioClient) Mute() error {
	return self.SetMute(true)
}

func (self *AudioClient) UnMute() error {
	return self.SetMute(false)
}

func (self *AudioClient) GetMute() (bool, error) {

	res, err := self.obsClient.Inputs.GetInputMute(inputs.NewGetInputMuteParams().WithInputUuid(self.inputUuid))
	if err != nil {
		return false, err
	}

	return res.InputMuted, nil
}
