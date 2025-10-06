package telop

import (
	"example.com/donguri-back/client"
	"example.com/donguri-back/util"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
)

type TelopClient struct {
	SharedClient *client.SharedClient
	SceneUuid    string // 配信画面のSceneのUUID
	InputUuid    string // テロップのInputの名前
	SceneItemId  int    // 配信画面のScene上のテロップのSceneItemId
}

// 新しいTelopClientを作成する。作成後には以下のことを保証する
// - テロップのItemが存在すること
// - テロップのItemが配信画面のSceneのSceneItemになり、非表示になっていること
// - TelopClientのすべてのフィールドが有効な値を持つこと
func NewTelopClient(sharedClient *client.SharedClient, sceneName string, inputName string) (*TelopClient, error) {
	// Sceneが存在しているかどうかを確認し、存在する場合はuuidを取得する。存在しない場合はエラー
	sceneUuid := ""
	err := sharedClient.With(func(client *goobs.Client) error {
		uuid, err := util.FindSceneByName(client, sceneName)
		if err != nil {
			return err
		}
		sceneUuid = uuid
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Inputが存在しているかどうかを確認し、存在する場合はuuidを取得する。存在しない場合は新しく作成する
	inputUuid := ""
	err = sharedClient.With(func(client *goobs.Client) error {
		uuid, err := util.FindInputByName(client, inputName)

		if err != nil {
			kind := util.DetermineTextInputKind(client)

			uuid_, err := util.CreateInputToDummyScene(client, kind, inputName)
			if err != nil {
				return err
			}

			uuid = uuid_
		}

		inputUuid = uuid
		return nil
	})
	if err != nil {
		return nil, err
	}

	// InputをSceneに追加し非表示状態にする。SceneItemIdを取得する。
	sceneItemId := 0
	err = sharedClient.With(func(client *goobs.Client) error {
		res, err := client.SceneItems.CreateSceneItem(
			sceneitems.NewCreateSceneItemParams().
				WithSceneUuid(sceneUuid).
				WithSourceUuid(inputUuid).
				WithSceneItemEnabled(false))
		if err != nil {
			return err
		}
		sceneItemId = res.SceneItemId
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &TelopClient{
		SharedClient: sharedClient,
		SceneUuid:    sceneUuid,
		InputUuid:    inputUuid,
		SceneItemId:  sceneItemId,
	}, nil
}

func (self *TelopClient) SetText(text string) error {
	return self.SharedClient.With(func(client *goobs.Client) error {
		_, err := client.Inputs.SetInputSettings(
			inputs.NewSetInputSettingsParams().
				WithInputUuid(self.InputUuid).
				WithInputSettings(map[string]any{"text": text}).
				WithOverlay(true))
		return err
	})
}

func (self *TelopClient) SetVisible(state bool) error {
	return self.SharedClient.With(func(client *goobs.Client) error {
		_, err := client.SceneItems.SetSceneItemEnabled(
			sceneitems.NewSetSceneItemEnabledParams().
				WithSceneUuid(self.SceneUuid).
				WithSceneItemId(self.SceneItemId).
				WithSceneItemEnabled(state))
		return err
	})
}
