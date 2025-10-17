package audio

import (
	"github.com/andreykaipov/goobs"
	// "github.com/andreykaipov/goobs/api/requests/inputs"　//後で使うかも

	// "github.com/andreykaipov/goobs/api/requests/sceneitems"　　//後で使うかも
	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
	// "github.com/sohosai/ultradonguri-server/internal/utils"　　//後で使うかも
)

type AudioClient struct {
	obsClient       *goobs.Client
	normalSceneUuid string //複数になる
	mutedSceneUuid  string
	cmSceneUuid     string
	isForceMuted    bool
}

func NewAudioClient(obsClient *goobs.Client, sceneName string, inputName string) (*AudioClient, error) {
	// sceneUuid := ""
	// inputUuid := ""
	// sceneItemId := 0

	// sceneNameからsceneUuidを取得する。取得できなければ、そのようなsceneNameのSceneが存在しないと判断してエラー
	// sceneUuid, err := utils.FindSceneByName(obsClient, sceneName)
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to find scene named %s: %s", sceneName, err)
	// }

	// inputNameからinputUuidを取得する。取得できなければ、そのようなinputNameのInputが存在しないと判断してエラー
	// inputUuid, err := utils.FindInputByName(obsClient, inputName)
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to find input named %s: %s", inputName, err)
	// }

	// 最初からオーディオのInputはSceneItemになっているはずなので、それを確認し、なっていなかったらエラー
	// 同一のSceneに同一のInputが複数のSceneItemとして登録されうるが、この場合一番始めに見つかったものが返されるのだろうか。
	// 一旦一番始めに見つかったものを使う
	// sceneItemId, err := obsClient.SceneItems.GetSceneItemId(
	// 	sceneitems.NewGetSceneItemIdParams().
	// 		WithSceneUuid(sceneUuid).
	// 		WithSceneName(sceneName).
	// 		WithSourceName(inputName))
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to find scene item id: %s", err)
	// }

	return &AudioClient{
		obsClient:       obsClient,
		normalSceneUuid: "hoge", //複数になる
		mutedSceneUuid:  "hoge",
		cmSceneUuid:     "hoge",
		isForceMuted:    false,
		// sceneUuid:   sceneUuid,
		// inputUuid:   inputUuid,
		// sceneItemId: sceneItemId.SceneItemId,
	}, nil
}

func (self *AudioClient) SetMute(state bool) error {
	// // fmt.Println("Input UUID:", self.inputUuid)
	// // resp, err := self.obsClient.Inputs.SetInputMute(
	// _, err := self.obsClient.Inputs.SetInputMute(
	// 	inputs.NewSetInputMuteParams().
	// 		// WithInputUuid(self.inputUuid).
	// 		WithInputMuted(state))

	// // fmt.Println("Response from SetMute:", resp)
	// return err
	return nil //すべての実装が変わる
}

func (self *AudioClient) Mute() error {
	return self.SetMute(true)
}

func (self *AudioClient) UnMute() error {
	return self.SetMute(false)
}

func (self *AudioClient) GetMute() (entities.MuteState, error) {

	// res, err := self.obsClient.Inputs.GetInputMute(inputs.NewGetInputMuteParams().WithInputUuid(self.inputUuid))
	// if err != nil {
	// 	return entities.MuteState{}, err
	// }

	// return entities.MuteState{IsMuted: res.InputMuted}, nil
	return entities.MuteState{IsMuted: false}, nil //常にfalseなので変更が必要
}
