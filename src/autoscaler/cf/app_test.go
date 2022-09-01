package cf_test

import (
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/cf"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/models"
	. "code.cloudfoundry.org/app-autoscaler/src/autoscaler/testhelpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	"net/http"
)

var _ = Describe("Cf client App", func() {
	BeforeEach(login)

	Describe("GetApp", func() {
		When("get app succeeds", func() {
			BeforeEach(func() {
				fakeCC.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/v3/apps/test-app-id"),
						VerifyHeaderKV("Authorization", "Bearer test-access-token"),
						RespondWith(http.StatusOK, LoadFile("testdata/app.json"), http.Header{"Content-Type": []string{"application/json"}}),
					),
				)
			})

			It("returns correct state", func() {
				app, err := cfc.GetApp("test-app-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(app).To(Equal(&cf.App{
					Guid:      "663e9a25-30ba-4fb4-91fa-9b784f4a8542",
					Name:      "autoscaler-1--0cde0e473e3e47f4",
					State:     "STOPPED",
					CreatedAt: ParseDate("2022-07-21T13:42:30Z"),
					UpdatedAt: ParseDate("2022-07-21T14:30:17Z"),
					Relationships: cf.Relationships{
						Space: &cf.Space{
							Data: cf.SpaceData{
								Guid: "3dfc4a10-6e70-44f8-989d-b3842f339e3b",
							},
						},
					},
				}))
			})
		})
	})

	Describe("GetAppAndProcesses", func() {

		When("get app & process return ok", func() {
			BeforeEach(func() {
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id/processes", CombineHandlers(
					RespondWith(http.StatusOK, LoadFile("testdata/app_processes.json"), http.Header{"Content-Type": []string{"application/json"}}),
				))
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id", CombineHandlers(
					RespondWith(http.StatusOK, LoadFile("testdata/app.json"), http.Header{"Content-Type": []string{"application/json"}}),
				))
			})

			It("returns correct state", func() {
				appAndProcess, err := cfc.GetAppAndProcesses("test-app-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(appAndProcess).To(Equal(&cf.AppAndProcesses{
					App: &cf.App{
						Guid:      "663e9a25-30ba-4fb4-91fa-9b784f4a8542",
						Name:      "autoscaler-1--0cde0e473e3e47f4",
						State:     "STOPPED",
						CreatedAt: ParseDate("2022-07-21T13:42:30Z"),
						UpdatedAt: ParseDate("2022-07-21T14:30:17Z"),
						Relationships: cf.Relationships{
							Space: &cf.Space{
								Data: cf.SpaceData{
									Guid: "3dfc4a10-6e70-44f8-989d-b3842f339e3b",
								},
							},
						},
					},
					Processes: cf.Processes{
						{
							Guid:       "6a901b7c-9417-4dc1-8189-d3234aa0ab82",
							Type:       "web",
							Instances:  5,
							MemoryInMb: 256,
							DiskInMb:   1024,
							CreatedAt:  ParseDate("2016-03-23T18:48:22Z"),
							UpdatedAt:  ParseDate("2016-03-23T18:48:42Z"),
						},
						{
							Guid:       "3fccacd9-4b02-4b96-8d02-8e865865e9eb",
							Type:       "worker",
							Instances:  1,
							MemoryInMb: 256,
							DiskInMb:   1024,
							CreatedAt:  ParseDate("2016-03-23T18:48:22Z"),
							UpdatedAt:  ParseDate("2016-03-23T18:48:42Z"),
						}},
				}))
			})
		})

		When("get app returns 500 and get process return ok", func() {
			BeforeEach(func() {
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id/processes", CombineHandlers(
					RespondWithJSONEncoded(http.StatusInternalServerError, models.CfInternalServerError),
				))
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id", CombineHandlers(
					RespondWith(http.StatusOK, LoadFile("testdata/app.json"), http.Header{"Content-Type": []string{"application/json"}}),
				))
			})

			It("should error", func() {
				appAndProcesses, err := cfc.GetAppAndProcesses("test-app-id")
				Expect(appAndProcesses).To(BeNil())
				Expect(err).To(MatchError(MatchRegexp(`get state&instances failed: failed GetAppProcesses 'test-app-id': failed getting page 1: failed getting cf.Response\[.*cf.Process\]:.*'UnknownError'`)))
			})
		})

		When("get processes return OK get app returns 500", func() {
			BeforeEach(func() {
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id/processes", CombineHandlers(
					RespondWith(http.StatusOK, LoadFile("testdata/app_processes.json"), http.Header{"Content-Type": []string{"application/json"}}),
				))
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id", CombineHandlers(
					RespondWithJSONEncoded(http.StatusInternalServerError, models.CfInternalServerError),
				))
			})

			It("should error", func() {
				appAndProcesses, err := cfc.GetAppAndProcesses("test-app-id")
				Expect(appAndProcesses).To(BeNil())
				Expect(err).To(MatchError(MatchRegexp("get state&instances failed: failed getting app 'test-app-id':.*'UnknownError'")))
			})
		})

		When("get processes return 500 and get app returns 500", func() {
			BeforeEach(func() {
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id/processes", CombineHandlers(
					RespondWithJSONEncoded(http.StatusInternalServerError, models.CfInternalServerError),
				))
				fakeCC.RouteToHandler("GET", "/v3/apps/test-app-id", CombineHandlers(
					RespondWithJSONEncoded(http.StatusInternalServerError, models.CfInternalServerError),
				))
			})

			It("should error", func() {
				appAndProcesses, err := cfc.GetAppAndProcesses("test-app-id")
				Expect(appAndProcesses).To(BeNil())
				Expect(err).To(MatchError(MatchRegexp(`get state&instances failed: .*'UnknownError'`)))
			})
		})
	})

	Describe("ScaleAppWebProcess", func() {
		JustBeforeEach(func() {
			err = cfc.ScaleAppWebProcess("test-app-id", 6)
		})

		When("scaling web app succeeds", func() {
			BeforeEach(func() {
				fakeCC.AppendHandlers(
					CombineHandlers(
						VerifyRequest("POST", "/v3/apps/test-app-id/processes/web/actions/scale"),
						VerifyHeaderKV("Authorization", "Bearer test-access-token"),
						VerifyJSON(`{"instances":6}`),
						RespondWith(http.StatusAccepted, LoadFile("scale_response.yml")),
					),
				)
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("scaling endpoint return 500", func() {
			BeforeEach(func() {
				setCfcClient(3)
				fakeCC.RouteToHandler("POST",
					"/v3/apps/test-app-id/processes/web/actions/scale",
					RespondWithJSONEncoded(http.StatusInternalServerError, models.CfInternalServerError))
			})

			It("should error correctly", func() {
				Expect(err).To(MatchError(MatchRegexp("failed scaling app 'test-app-id' to 6: POST request failed:.*'UnknownError'.*")))
			})
		})

	})

})
