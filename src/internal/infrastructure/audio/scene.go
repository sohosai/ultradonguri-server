package audio

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/scenes"
)

func (self *AudioClient) SetNormalScene() error {
	//CMシーン中のシーン切り替えを許さない（CM解除以外の）
	//仕様にどうすべきか明記されていなかったため、変更の可能性あり
	current_scene, err := self.GetCurrentScene()
	if err != nil {
		return err
	}
	//CMから出れない！？！？消しましょう
	if current_scene == self.scenes.CM {
		return fmt.Errorf("cannot change scene: it's CM scene")
	}
	_, err = self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.Normal,
	})
	return err
}

func (self *AudioClient) SetMutedScene() error {
	current_scene, err := self.GetCurrentScene()
	if err != nil {
		return err
	}
	if current_scene == self.scenes.CM {
		return fmt.Errorf("cannot change scene: it's CM scene")
	}
	_, err = self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.Muted,
	})
	return err
}

func (self *AudioClient) SetCMScene() error {
	if !self.isConversion {
		return fmt.Errorf("cannot change force_mute state: it's not conversion now")
	} else if self.isForceMuted {
		//force_mute中もCMにできない
		return fmt.Errorf("cannot change force_mute state: force_muted")
	}

	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.CM,
	})
	return err
}

// SceneNameが必要なら返り値に含めてもよい
// func (self *AudioClient) GetCurrentScene() (sceneName string, sceneUUID string, err error) {
func (self *AudioClient) GetCurrentScene() (sceneUUID string, err error) {
	resp, err := self.obsClient.Scenes.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})
	if err != nil {
		return "", err
	}
	// return resp.CurrentProgramSceneName, resp.CurrentProgramSceneUuid, nil
	return resp.CurrentProgramSceneUuid, nil
}
