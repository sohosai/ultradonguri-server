package audio

import (
	"fmt"

	"example.com/donguri-back/client"
	"example.com/donguri-back/util"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
)

type AudioClient struct {
	sharedClient *client.SharedClient
	sceneUuid    string
	inputUuid    string
	sceneItemId  int
}

func NewAudioClient(sharedClient *client.SharedClient, sceneName string, inputName string) (*AudioClient, error) {
	sceneUuid := ""
	inputUuid := ""
	sceneItemId := 0

	if err := sharedClient.With(func(client *goobs.Client) error {
		// sceneNameからsceneUuidを取得する。取得できなければ、そのようなsceneNameのSceneが存在しないと判断してエラー
		sceneUuid_, err := util.FindSceneByName(client, sceneName)
		if err != nil {
			return fmt.Errorf("Failed to find scene named %s: %s", sceneName, err)
		}

		// inputNameからinputUuidを取得する。取得できなければ、そのようなinputNameのInputが存在しないと判断してエラー
		inputUuid_, err := util.FindInputByName(client, inputName)
		if err != nil {
			return fmt.Errorf("Failed to find input named %s: %s", inputName, err)
		}

		// 最初からオーディオのInputはSceneItemになっているはずなので、それを確認し、なっていなかったらエラー
		// 同一のSceneに同一のInputが複数のSceneItemとして登録されうるが、この場合一番始めに見つかったものが返されるのだろうか。
		// 一旦一番始めに見つかったものを使う
		sceneItemId_, err := client.SceneItems.GetSceneItemId(
			sceneitems.NewGetSceneItemIdParams().
				WithSceneUuid(sceneUuid).
				WithSceneName(sceneName).
				WithSourceName(inputName))
		if err != nil {
			return fmt.Errorf("Failed to find scene item id: %s", err)
		}

		sceneUuid = sceneUuid_
		inputUuid = inputUuid_
		sceneItemId = sceneItemId_.SceneItemId

		return nil
	}); err != nil {
		return nil, err
	}

	return &AudioClient{
		sharedClient: sharedClient,
		sceneUuid:    sceneUuid,
		inputUuid:    inputUuid,
		sceneItemId:  sceneItemId,
	}, nil
}

func (self *AudioClient) SetMute(state bool) error {
	return self.sharedClient.With(func(client *goobs.Client) error {
		if _, err := client.Inputs.SetInputMute(
			inputs.NewSetInputMuteParams().
				WithInputUuid(self.inputUuid).
				WithInputMuted(state)); err != nil {
			return err
		}

		return nil
	})
}

func (self *AudioClient) Mute() error {
	return self.SetMute(true)
}

func (self *AudioClient) UnMute() error {
	return self.SetMute(false)
}

func (self *AudioClient) GetMute() (bool, error) {
	isMuted := false

	if err := self.sharedClient.With(func(client *goobs.Client) error {
		res, err := client.Inputs.GetInputMute(inputs.NewGetInputMuteParams().WithInputUuid(self.inputUuid))
		if err != nil {
			return err
		}

		isMuted = res.InputMuted

		return nil
	}); err != nil {
		return false, err
	}

	return isMuted, nil
}
