package do

import (
	"github.com/digitalocean/godo"
	"github.com/digitalocean/godo/util"
)

// DropletIPTable is a table of interface IPS.
type DropletIPTable map[InterfaceType]string

// InterfaceType is a an interface type.
type InterfaceType string

const (
	// InterfacePublic is a public interface.
	InterfacePublic InterfaceType = "public"
	// InterfacePrivate is a private interface.
	InterfacePrivate InterfaceType = "private"
)

// Droplet is a wrapper for godo.Droplet
type Droplet struct {
	*godo.Droplet
}

// IPs returns a map of interface.s
func (d *Droplet) IPs() DropletIPTable {
	t := DropletIPTable{}
	for _, in := range d.Networks.V4 {
		switch in.Type {
		case "public":
			t[InterfacePublic] = in.IPAddress
		case "private":
			t[InterfacePrivate] = in.IPAddress
		}
	}

	return t
}

// Droplets is a slice of Droplet.
type Droplets []Droplet

// Kernel is a wrapper for godo.Kernel
type Kernel struct {
	*godo.Kernel
}

// Kernels is a slice of Kernel.
type Kernels []Kernel

// DropletsService is an interface for interacting with DigitalOcean's droplet api.
type DropletsService interface {
	List() (Droplets, error)
	Get(int) (*Droplet, error)
	Create(*godo.DropletCreateRequest, bool) (*Droplet, error)
	CreateMultiple(*godo.DropletMultiCreateRequest) (Droplets, error)
	Delete(int) error
	Kernels(int) (Kernels, error)
	Snapshots(int) (Images, error)
	Backups(int) (Images, error)
	Actions(int) (Actions, error)
	Neighbors(int) (Droplets, error)
}

type dropletsService struct {
	client *godo.Client
}

var _ DropletsService = &dropletsService{}

// NewDropletsService builds a DropletsService instance.
func NewDropletsService(client *godo.Client) DropletsService {
	return &dropletsService{
		client: client,
	}
}

func (ds *dropletsService) List() (Droplets, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Droplets.List(opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Droplets, len(si))
	for i := range si {
		a := si[i].(godo.Droplet)
		list[i] = Droplet{Droplet: &a}
	}

	return list, nil
}

func (ds *dropletsService) Get(id int) (*Droplet, error) {
	d, _, err := ds.client.Droplets.Get(id)
	if err != nil {
		return nil, err
	}

	return &Droplet{Droplet: d}, nil
}

func (ds *dropletsService) Create(dcr *godo.DropletCreateRequest, wait bool) (*Droplet, error) {
	d, resp, err := ds.client.Droplets.Create(dcr)
	if err != nil {
		return nil, err
	}

	if wait {
		var action *godo.LinkAction
		for _, a := range resp.Links.Actions {
			if a.Rel == "create" {
				action = &a
				break
			}
		}

		if action != nil {
			_ = util.WaitForActive(ds.client, action.HREF)
			doDroplet, err := ds.Get(d.ID)
			if err != nil {
				return nil, err
			}
			d = doDroplet.Droplet
		}
	}

	return &Droplet{Droplet: d}, nil
}

func (ds *dropletsService) CreateMultiple(dmcr *godo.DropletMultiCreateRequest) (Droplets, error) {
	godoDroplets, _, err := ds.client.Droplets.CreateMultiple(dmcr)
	if err != nil {
		return nil, err
	}

	var droplets Droplets
	for _, d := range godoDroplets {
		droplets = append(droplets, Droplet{Droplet: &d})
	}

	return droplets, nil
}

func (ds *dropletsService) Delete(id int) error {
	_, err := ds.client.Droplets.Delete(id)
	return err
}

func (ds *dropletsService) Kernels(id int) (Kernels, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Droplets.Kernels(id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Kernels, len(si))
	for i := range si {
		a := si[i].(godo.Kernel)
		list[i] = Kernel{Kernel: &a}
	}

	return list, nil
}

func (ds *dropletsService) Snapshots(id int) (Images, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Droplets.Snapshots(id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Images, len(si))
	for i := range si {
		a := si[i].(godo.Image)
		list[i] = Image{Image: &a}
	}

	return list, nil
}

func (ds *dropletsService) Backups(id int) (Images, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Droplets.Backups(id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Images, len(si))
	for i := range si {
		a := si[i].(godo.Image)
		list[i] = Image{Image: &a}
	}

	return list, nil
}

func (ds *dropletsService) Actions(id int) (Actions, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Droplets.Actions(id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Actions, len(si))
	for i := range si {
		a := si[i].(godo.Action)
		list[i] = Action{Action: &a}
	}

	return list, nil
}

func (ds *dropletsService) Neighbors(id int) (Droplets, error) {
	list, _, err := ds.client.Droplets.Neighbors(id)
	if err != nil {
		return nil, err
	}

	var droplets Droplets
	for _, d := range list {
		droplets = append(droplets, Droplet{Droplet: &d})
	}

	return droplets, nil
}
