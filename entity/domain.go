package entity

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

type Domain struct {
	Name       string
	SubDomains []*Domain
	Components []DomainComponent
	Interfaces []*Interface
}

func (d *Domain) GetDomain(p string) *Domain {
	if td := d.findDomain(p, d); td != nil {
		return td
	}

	nd := NewDomain(p)
	if err := d.AddDomain(nd); err != nil {
		fmt.Println(err)
		return nil
	}

	return nd
}

func (d *Domain) AddDomain(nd *Domain) error {
	if fd := d.findDomain(nd.Name, d); fd != nil {
		return nil
	}

	parent := path.Dir(nd.Name)
	if parent == d.Name {
		d.SubDomains = append(d.SubDomains, nd)
	} else {
		if strings.HasPrefix(parent, d.Name) {
			sp, err := d.split(parent)
			if err != nil {
				return err
			}
			d.completeDomain(d, sp)
			if p := d.findDomain(parent, d); p != nil {
				p.SubDomains = append(p.SubDomains, nd)
			}
		} else {
			return errors.New("invalid domain path")
		}
	}

	return nil
}

func (d *Domain) split(dp string) ([]string, error) {
	if strings.Contains(dp, d.Name) {
		uri := dp[len(d.Name):]
		turi := strings.TrimPrefix(uri, "/")
		turi = strings.TrimSuffix(turi, "/")
		segments := strings.Split(turi, "/")

		var paths [][]string
		for i := 1; i <= len(segments); i++ {
			paths = append(paths, segments[:i])
		}

		var dps []string
		for _, p := range paths {
			dps = append(dps, path.Join(d.Name, strings.Join(p, "/")))
		}

		return dps, nil
	}
	return nil, errors.New("invalid domain path")
}

func (d *Domain) completeDomain(root *Domain, path []string) {
	if len(path) == 0 {
		return
	}

	var subdomain *Domain
	for _, sub := range root.SubDomains {
		if sub.Name == path[0] {
			subdomain = sub
			break
		}
	}

	if subdomain == nil {
		subdomain = &Domain{
			Name:       path[0],
			SubDomains: []*Domain{},
		}
		root.SubDomains = append(root.SubDomains, subdomain)
	}

	d.completeDomain(subdomain, path[1:])
}

func (d *Domain) findDomain(name string, currentDomain *Domain) *Domain {
	if currentDomain.Name == name {
		return currentDomain
	}
	for _, sd := range d.SubDomains {
		if p := sd.findDomain(name, sd); p != nil {
			return p
		}
	}
	return nil
}

func NewDomain(name string) *Domain {
	return &Domain{
		Name:       name,
		SubDomains: nil,
		Components: nil,
	}
}
