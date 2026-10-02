package scene

import (
	"fmt"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/sohosai/ultradonguri-server/internal/utils"
)

type SceneManager struct {
	obsClient  *goobs.Client
	scenes     Scenes
	sceneType  SceneType // ファイルバックアップのために今設定されているSceneTypeを保存しておく。setSceneというメソッドだけが触る
	backupPath string
}

// sceneのName or UUIDをまとめた型
type Scenes struct {
	Normal string
	CM     string
	Burari string
}

type SceneNames struct {
	Normal string
	CM     string
	Burari string
}

type SceneType = int

const (
	Normal SceneType = iota
	CM
	Burari
)

type Backup struct {
	SceneType SceneType `json:"scene_type"`
}

func NewSceneManager(obsClient *goobs.Client, sceneNames SceneNames, backupPath string) (*SceneManager, error) {
	return newSceneManager(obsClient, sceneNames, backupPath, Normal)
}

func RestoreSceneManager(obsClient *goobs.Client, sceneNames SceneNames, backupPath string) (*SceneManager, error) {
	backupRaw, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, err
	}

	savedInfo, err := utils.JsonStrictUnmarshal[Backup](backupRaw)
	if err != nil {
		return nil, err
	}

	sceneManager, err := newSceneManager(obsClient, sceneNames, backupPath, savedInfo.SceneType)
	if err != nil {
		return nil, err
	}

	return sceneManager, err
}

func newSceneManager(obsClient *goobs.Client, sceneNames SceneNames, backupPath string, initialScene SceneType) (*SceneManager, error) {
	// sceneNameからsceneUuidを取得する。取得できなければ、そのようなsceneNameのSceneが存在しないと判断してエラー
	// sceneNameはobsから容易に変更可能なので安全性のためにUUIDを用いる
	resolve := func(name string) (string, error) {
		uuid, err := utils.FindSceneByName(obsClient, name)
		if err != nil {
			return "", fmt.Errorf("failed to find scene named %s: %w", name, err)
		}
		return uuid, nil
	}

	normalUUID, err := resolve(sceneNames.Normal)
	if err != nil {
		return nil, err
	}
	cmUUID, err := resolve(sceneNames.CM)
	if err != nil {
		return nil, err
	}
	burariUUID, err := resolve(sceneNames.Burari)
	if err != nil {
		return nil, err
	}

	sceneUUIDs := Scenes{
		Normal: normalUUID,
		CM:     cmUUID,
		Burari: burariUUID,
	}

	sceneManager := &SceneManager{
		obsClient:  obsClient,
		scenes:     sceneUUIDs,
		backupPath: backupPath,
	}

	switch initialScene {
	case Normal:
		err = sceneManager.SetNormalScene()
		sceneManager.sceneType = Normal

	case CM:
		err = sceneManager.SetCMScene()
		sceneManager.sceneType = CM

	case Burari:
		err = sceneManager.SetBurariScene()
		sceneManager.sceneType = Burari

	default:
		err = sceneManager.SetNormalScene()
		sceneManager.sceneType = Normal
	}

	return sceneManager, err
}

func (self *SceneManager) IsCm() (bool, error) {
	currentScene, err := self.GetCurrentScene()
	if err != nil {
		return false, err
	}

	return currentScene == self.scenes.CM, nil
}
