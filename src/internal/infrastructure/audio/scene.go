package audio

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/scenes"
)

func (self *AudioClient) SetNormalScene() error {
	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.Normal,
	})
	return err
}

func (self *AudioClient) SetMutedScene() error {
	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.Muted,
	})
	return err
}

func (self *AudioClient) SetCMScene() error {
	if !self.isConversion {
		return fmt.Errorf("cannot change force_mute state: it's not conversion now")
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
