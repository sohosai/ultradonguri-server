package audio

import (
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
	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &self.scenes.CM,
	})
	return err
}

func (self *AudioClient) GetCurrentScene() (sceneName string, sceneUUID string, err error) {
	resp, err := self.obsClient.Scenes.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})
	if err != nil {
		return "", "", err
	}
	return resp.CurrentProgramSceneName, resp.CurrentProgramSceneUuid, nil
}
