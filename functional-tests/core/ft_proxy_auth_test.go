package hoverfly_test

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly

		username = "ft_user"
		password = "ft_password"
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	Context("with auth turned on", func() {

		BeforeEach(func() {
			hoverfly.Start("-db", "boltdb", "-add", "-username", username, "-password", password)
			hoverfly.Start("-db", "boltdb", "-auth", "true")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should return a 407 when trying to proxy without auth credentials", func() {
			resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
			Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
		})

		It("should return a 407 (no match in simulate mode) when trying to proxy with incorrect auth credentials", func() {
			resp := hoverfly.ProxyWithAuth(sling.New().Get("http://hoverfly.io"), "incorrect", "incorrect")
			Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
		})

		It("should return a 502 (no match in simulate mode) when trying to proxy with auth credentials", func() {
			resp := hoverfly.ProxyWithAuth(sling.New().Get("http://hoverfly.io"), username, password)
			Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
		})
	})
})
