package action

import "fmt"

type factoryInterface interface {
	GenerateMission(m string) missionInterface
}

type ActionFactory struct {
}

func (a ActionFactory) GenerateMission(m string) missionInterface {
	switch m {
	case "shadow":
		return shadow{}
	case "upgrade":
		return upgrade{}
	case "reboot":
		return reboot{}
	default:
		return nil
	}
}

type missionInterface interface {
	DoMission()
}

type shadow struct{}

func (s shadow) DoMission() {
	fmt.Println("this is shadow")
}

type upgrade struct{}

func (u upgrade) DoMission() {
	fmt.Println("this is upgrade")
}

type reboot struct{}

func (r reboot) DoMission() {
	fmt.Println("this is reboot")
}
