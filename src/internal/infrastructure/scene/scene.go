package scene

import (
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

func (self *SceneManager) SetNormalScene() error {
	//CMシーン中のシーン切り替えを許さない（CM解除以外の）
	//仕様にどうすべきか明記されていなかったため、変更の可能性あり

	err := self.setScene(Normal, self.scenes.Normal)

	return err
}

func (self *SceneManager) SetCMScene() error {
	// isConversionの管理はSceneManagerの責任ではないので外から受け取る
	err := self.setScene(CM, self.scenes.CM)
	return err
}

func (self *SceneManager) SetBurariScene() error {
	err := self.setScene(Burari, self.scenes.Burari)
	return err
}

func (self *SceneManager) setScene(sceneType SceneType, sceneUuid string) error {
	_, err := self.obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &sceneUuid,
	})
	if err != nil {
		return err
	}

	self.sceneType = sceneType
	return self.saveToFile()
}

func (self *SceneManager) GetCurrentScene() (sceneUUID string, err error) {
	resp, err := self.obsClient.Scenes.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})
	if err != nil {
		return "", err
	}
	return resp.CurrentProgramSceneUuid, nil
}
