package audio

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

func SwitchScene(client *goobs.Client, sceneUUID string) error {
	_, err := client.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneUuid: &sceneUUID,
	})
	return err
}

func GetCurrentScene(client *goobs.Client) (sceneName string, sceneUUID string, err error) {
	resp, err := client.Scenes.GetCurrentProgramScene(&scenes.GetCurrentProgramSceneParams{})
	if err != nil {
		return "", "", err
	}
	return resp.CurrentProgramSceneName, resp.CurrentProgramSceneUuid, nil
}
