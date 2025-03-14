package ticketing

import "protohackers/6_speed/infra"

type Camera struct {
	*infra.IAmACamera
}

type Controller struct {
	cameras     []*infra.IAmACamera
	dispatchers []*infra.IAmADispatcher
	plates      []*infra.Plate
}

func MakeGlobal() *Controller {
	return &Controller{
		cameras:     make([]*infra.IAmACamera, 0),
		dispatchers: make([]*infra.IAmADispatcher, 0),
		plates:      make([]*infra.Plate, 0),
	}
}

func (g *Controller) AddCamera(cam *infra.IAmACamera) {
	g.cameras = append(g.cameras, cam)
}

func (g *Controller) AddDispatcher(dis *infra.IAmADispatcher) {
	g.dispatchers = append(g.dispatchers, dis)
}

func (g *Controller) AddPlates(dis *infra.IAmADispatcher) {

}
