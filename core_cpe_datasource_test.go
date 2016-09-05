package main

import (
	"testing"
	"time"

	"github.com/MustWin/baremetal-sdk-go"
	"github.com/MustWin/terraform-Oracle-BareMetal-Provider/client/mocks"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/stretchr/testify/suite"
)

type DatasourceCoreCpeTestSuite struct {
	suite.Suite
	Client       *mocks.BareMetalClient
	Config       string
	Provider     terraform.ResourceProvider
	Providers    map[string]terraform.ResourceProvider
	ResourceName string
}

func (s *DatasourceCoreCpeTestSuite) SetupTest() {
	s.Client = &mocks.BareMetalClient{}
	s.Provider = Provider(func(d *schema.ResourceData) (interface{}, error) {
		return s.Client, nil
	})

	s.Providers = map[string]terraform.ResourceProvider{
		"baremetal": s.Provider,
	}
	s.Config = `
    data "baremetal_core_cpes" "s" {
      compartment_id = "compartmentid"
    }
  `
	s.Config += testProviderConfig
	s.ResourceName = "data.baremetal_core_cpes.s"

}

func (s *DatasourceCoreCpeTestSuite) TestCpeList() {

	s.Client.On(
		"ListCpes",
		"compartmentid",
		[]baremetal.Options{},
	).Return(
		&baremetal.ListCpes{

			Cpes: []baremetal.Cpe{
				baremetal.Cpe{
					ID:            "id1",
					CompartmentID: "compartmentid",
					DisplayName:   "name",
					IPAddress:     "10.10.10.2",
					TimeCreated:   baremetal.Time{Time: time.Now()},
				},

				baremetal.Cpe{
					ID:            "id2",
					CompartmentID: "compartmentid",
					DisplayName:   "name",
					IPAddress:     "10.10.10.3",
					TimeCreated:   baremetal.Time{Time: time.Now()},
				},
			},
		},
		nil,
	)

	resource.UnitTest(s.T(), resource.TestCase{
		PreventPostDestroyRefresh: true,
		Providers:                 s.Providers,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: s.Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "compartment_id", "compartmentid"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.0.ip_address", "10.10.10.2"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.0.id", "id1"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.1.ip_address", "10.10.10.3"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.1.id", "id2"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.#", "2"),
				),
			},
		},
	},
	)

	s.Client.AssertCalled(s.T(), "ListCpes", "compartmentid", []baremetal.Options{})

}

func (s *DatasourceCoreCpeTestSuite) TestCpePagedList() {

	s.Client.On(
		"ListCpes",
		"compartmentid",
		[]baremetal.Options{},
	).Return(
		&baremetal.ListCpes{
			ResourceContainer: baremetal.ResourceContainer{
				NextPage: "nextpage",
			},
			Cpes: []baremetal.Cpe{
				baremetal.Cpe{
					ID:            "id1",
					CompartmentID: "compartmentid",
					DisplayName:   "name",
					IPAddress:     "10.10.10.2",
					TimeCreated:   baremetal.Time{Time: time.Now()},
				},

				baremetal.Cpe{
					ID:            "id2",
					CompartmentID: "compartmentid",
					DisplayName:   "name",
					IPAddress:     "10.10.10.3",
					TimeCreated:   baremetal.Time{Time: time.Now()},
				},
			},
		},
		nil,
	)

	s.Client.On(
		"ListCpes",
		"compartmentid",
		[]baremetal.Options{baremetal.Options{Page: "nextpage"}},
	).Return(
		&baremetal.ListCpes{
			Cpes: []baremetal.Cpe{
				baremetal.Cpe{
					ID:            "id3",
					CompartmentID: "compartmentid",
					DisplayName:   "name",
					IPAddress:     "10.10.10.4",
					TimeCreated:   baremetal.Time{Time: time.Now()},
				},

				baremetal.Cpe{
					ID:            "id4",
					CompartmentID: "compartmentid",
					DisplayName:   "name",
					IPAddress:     "10.10.10.5",
					TimeCreated:   baremetal.Time{Time: time.Now()},
				},
			},
		},
		nil,
	)

	resource.UnitTest(s.T(), resource.TestCase{
		PreventPostDestroyRefresh: true,
		Providers:                 s.Providers,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: s.Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "compartment_id", "compartmentid"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.0.ip_address", "10.10.10.2"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.0.id", "id1"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.1.ip_address", "10.10.10.3"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.1.id", "id2"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.2.id", "id3"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.3.id", "id4"),
					resource.TestCheckResourceAttr(s.ResourceName, "cpes.#", "4"),
				),
			},
		},
	},
	)

	s.Client.AssertCalled(s.T(), "ListCpes", "compartmentid", []baremetal.Options{})

}

func TestDatasourceCoreCpeTestSuite(t *testing.T) {
	suite.Run(t, new(DatasourceCoreCpeTestSuite))
}
