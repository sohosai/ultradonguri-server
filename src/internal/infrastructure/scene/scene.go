package scene

import (
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

func (self *SceneManager) SetNormalScene() error {
	//CMシーン中のシーン切り替えを許さない（CM解除以外の）
	//仕様にどうすべきか明記されていなかったため、変更の可能性あり

	if self.isForceMutedFlag {
		// force_mute中はNormalシーンに移行しない
		err := self.SetMutedScene()
		// return fmt.Errorf("Failed to switch scene to Normal: force_muted")
		return err
	}

	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.Normal,
	})

	return err
}

func (self *SceneManager) SetMutedScene() error {
	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.Muted,
	})
	return err
}

func (self *SceneManager) SetCMScene() error {
	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.CM,
	})
	return err
}

func (self *SceneManager) GetCurrentScene() (sceneUUID string, err error) {
	resp, err := self.obsClient.Scenes.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})
	if err != nil {
		return "", err
	}
	return resp.CurrentProgramSceneUuid, nil
}
