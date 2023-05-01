package entity

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewDomain(t *testing.T) {
	name := "test"
	domain := NewDomain(name)
	if domain.Name != name {
		t.Errorf("Expected domain name to be %s, but got %s", name, domain.Name)
	}
	if domain.SubDomains != nil {
		t.Errorf("Expected domain subdomains to be nil, but got %v", domain.SubDomains)
	}
	if domain.Components != nil {
		t.Errorf("Expected domain components to be nil, but got %v", domain.Components)
	}
}

func TestFindParent(t *testing.T) {
	root := &Domain{
		Name: "root",
		SubDomains: []*Domain{
			&Domain{
				Name: "sub1",
				SubDomains: []*Domain{
					&Domain{
						Name: "sub1.1",
					},
					&Domain{
						Name: "sub1.2",
					},
				},
			},
			&Domain{
				Name: "sub2",
				SubDomains: []*Domain{
					&Domain{
						Name: "sub2.1",
					},
					&Domain{
						Name: "sub2.2",
					},
				},
			},
		},
	}

	parentDomain := root.findDomain("sub1", root)
	if parentDomain == nil {
		t.Error("Failed to find parent domain")
	}
	if parentDomain.Name != "sub1" {
		t.Errorf("Expected parent domain to be sub1, but got %s", parentDomain.Name)
	}

	missingDomain := root.findDomain("missing", root)
	if missingDomain != nil {
		t.Error("Found a missing domain")
	}
}

func TestCompleteDomain(t *testing.T) {
	root := &Domain{
		Name: "root",
		SubDomains: []*Domain{
			&Domain{
				Name: "root.sub1",
			},
			&Domain{
				Name: "root.sub2",
			},
		},
	}

	path := []string{"root.sub1", "root.sub1.1"}

	root.completeDomain(root, path)

	parentDomain := root.findDomain("root.sub1", root)
	if parentDomain == nil {
		t.Error("Failed to find parent domain")
	}
	if parentDomain.Name != "root.sub1" {
		t.Errorf("Expected parent domain to be sub1, but got %s", parentDomain.Name)
	}

	subdomain := parentDomain.SubDomains[0]
	if subdomain == nil {
		t.Error("Failed to find subdomain")
	}
	if subdomain.Name != "root.sub1.1" {
		t.Errorf("Expected subdomain to be sub1.1, but got %s", subdomain.Name)
	}
}

func TestSplit(t *testing.T) {
	d := &Domain{
		Name: "example",
	}

	// Test valid input.
	input := "example/foo/bar"
	expected := []string{"example/foo", "example/foo/bar"}
	result, err := d.split(input)

	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result to be %v, but got %v", expected, result)
	}

	// Test invalid input.
	input = "test/foo/bar"
	expectedErr := errors.New("invalid domain path")
	result, err = d.split(input)

	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error to be %v, but got %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("Expected result to be nil, but got %v", result)
	}
}

func TestAddDomain(t *testing.T) {
	d := &Domain{
		Name: "example",
	}

	// Test adding a domain to an existing parent.
	nd := &Domain{
		Name: "example/foo",
	}
	err := d.AddDomain(nd)
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	if len(d.SubDomains) != 1 || d.SubDomains[0] != nd {
		t.Error("Expected domain to be added, but it was not")
	}

	// Test adding a domain to a non-existing parent.
	nd = &Domain{
		Name: "test/foo",
	}
	err = d.AddDomain(nd)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if err.Error() != "invalid domain path" {
		t.Errorf("Expected error to be 'invalid domain path', but got %v", err)
	}

	// Test adding a domain with a duplicate name.
	nd = &Domain{
		Name: "example/foo",
	}
	err = d.AddDomain(nd)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestDomain_GetDomain(t *testing.T) {
	// Create a new domain with some subdomains
	domain := NewDomain("root")
	subdomain1 := NewDomain("root/subdomain1")
	subdomain2 := NewDomain("root/subdomain2")
	domain.SubDomains = []*Domain{subdomain1, subdomain2}
	// Test case 1: Get an existing subdomain
	got1 := domain.GetDomain("root/subdomain1")
	if !reflect.DeepEqual(got1, subdomain1) {
		t.Errorf("Got: %v, want: %v", got1, subdomain1)
	}

	// Test case 2: Get a non-existing subdomain
	got2 := domain.GetDomain("root/subdomain3")
	if got2 == nil {
		t.Errorf("Got: %v, want non-nil", got2)
	}

	// Test case 3: Get a subdomain that needs to be added
	got3 := domain.GetDomain("root/subdomain3/subsubdomain1")
	want3 := NewDomain("root/subdomain3/subsubdomain1")
	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("Got: %v, want: %v", got3, want3)
	}

	// Test case 4: Get the root domain
	got4 := domain.GetDomain("root")
	if !reflect.DeepEqual(got4, domain) {
		t.Errorf("Got: %v, want: %v", got4, domain)
	}
}
