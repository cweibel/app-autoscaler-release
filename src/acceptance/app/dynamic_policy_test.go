package app_test

import (
	"acceptance"
	. "acceptance/helpers"
	"fmt"
	"sync"
	"time"

	cfh "github.com/cloudfoundry/cf-test-helpers/v2/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AutoScaler dynamic policy", func() {
	var (
		policy         string
		err            error
		doneChan       chan bool
		doneAcceptChan chan bool
		ticker         *time.Ticker
		maxHeapLimitMb int
		wg             sync.WaitGroup
	)

	const minimalMemoryUsage = 28 // observed minimal memory usage by the test app

	JustBeforeEach(func() {
		appName = CreateTestApp(cfg, "dynamic-policy", initialInstanceCount)

		appGUID, err = GetAppGuid(cfg, appName)
		Expect(err).NotTo(HaveOccurred())
		StartApp(appName, cfg.CfPushTimeoutDuration())
		instanceName = CreatePolicy(cfg, appName, appGUID, policy)

		wg = sync.WaitGroup{}
	})
	BeforeEach(func() {
		maxHeapLimitMb = cfg.NodeMemoryLimit - minimalMemoryUsage
	})

	AfterEach(AppAfterEach)

	Context("when scaling by memoryused", func() {

		Context("There is a scale out and scale in policy", func() {
			var heapToUse int
			BeforeEach(func() {
				heapToUse = min(maxHeapLimitMb, 200)
				expectedAverageUsageAfterScaling := float64(heapToUse)/2 + minimalMemoryUsage
				policy = GenerateDynamicScaleOutAndInPolicy(1, 2, "memoryused", int64(0.9*expectedAverageUsageAfterScaling), int64(heapToUse))
				initialInstanceCount = 1
			})

			It("should scale out and then back in.", func() {
				By(fmt.Sprintf("Use heap %d MB of heap on app", heapToUse))
				CurlAppInstance(cfg, appName, 0, fmt.Sprintf("/memory/%d/5", heapToUse))

				By("wait for scale to 2")
				WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)

				By("Drop memory used by app")
				CurlAppInstance(cfg, appName, 0, "/memory/close")

				By("Wait for scale to minimum instances")
				WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
			})
		})
	})

	Context("when scaling by memoryutil", func() {

		Context("when memoryutil", func() {
			BeforeEach(func() {
				//current app resident size is 66mb so 66/128mb is 55%
				policy = GenerateDynamicScaleOutAndInPolicy(1, 2, "memoryutil", 58, 63)
				initialInstanceCount = 1
			})

			It("should scale out and back in", func() {
				heapToUse := min(maxHeapLimitMb, int(float64(cfg.NodeMemoryLimit)*0.80))
				By(fmt.Sprintf("use 80%% or %d MB of memory in app", heapToUse))
				CurlAppInstance(cfg, appName, 0, fmt.Sprintf("/memory/%d/5", heapToUse))

				By("Wait for scale to 2 instances")
				WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)

				By("drop memory used")
				CurlAppInstance(cfg, appName, 0, "/memory/close")

				By("Wait for scale down to 1 instance")
				WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
			})
		})
	})

	Context("when scaling by responsetime", func() {

		JustBeforeEach(func() {
			doneChan = make(chan bool)
			doneAcceptChan = make(chan bool)
		})

		AfterEach(func() {
			close(doneChan)
			Eventually(doneAcceptChan, 10*time.Second).Should(Receive())
		})

		Context("when responsetime is greater than scaling out threshold", func() {

			BeforeEach(func() {
				policy = GenerateDynamicScaleOutPolicy(1, 2, "responsetime", 800)
				initialInstanceCount = 1
			})

			JustBeforeEach(func() {
				// simulate ongoing ~20 slow requests per second
				ticker = time.NewTicker(50 * time.Millisecond)
				appUri := cfh.AppUri(appName, "/responsetime/slow/1000", cfg)
				go func(chan bool) {
					defer GinkgoRecover()
					for {
						select {
						case <-doneChan:
							ticker.Stop()
							doneAcceptChan <- true
							return
						case <-ticker.C:
							wg.Add(1)
							go func() {
								defer wg.Done()
								cfh.Curl(cfg, appUri)
							}()
						}
					}
				}(doneChan)
			})

			It("should scale out", Label(acceptance.LabelSmokeTests), func() {
				WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)
			})

			AfterEach(func() {
				// ginkgo requires all go-routines to be finished before reporting the result.
				// let's wait for all curl executing go-routines to return.
				wg.Wait()
			})
		})

		Context("when responsetime is less than scaling in threshold", func() {

			BeforeEach(func() {
				policy = GenerateDynamicScaleInPolicyBetween("responsetime", 800, 1200)
				initialInstanceCount = 2
			})

			JustBeforeEach(func() {
				// simulate ongoing ~20 fast requests per second
				ticker = time.NewTicker(50 * time.Millisecond)
				appUri := cfh.AppUri(appName, "/responsetime/slow/1000", cfg)
				go func(chan bool) {
					defer GinkgoRecover()
					for {
						select {
						case <-doneChan:
							ticker.Stop()
							doneAcceptChan <- true
							return
						case <-ticker.C:
							wg.Add(1)
							go func() {
								defer wg.Done()
								cfh.Curl(cfg, appUri)
							}()
						}
					}
				}(doneChan)
			})

			It("should scale in", func() {
				WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
			})

			AfterEach(func() {
				// ginkgo requires all go-routines to be finished before reporting the result.
				// let's wait for all curl executing go-routines to return.
				wg.Wait()
			})
		})

	})

	Context("when scaling by throughput", func() {
		JustBeforeEach(func() {
			doneChan = make(chan bool)
			doneAcceptChan = make(chan bool)
		})

		AfterEach(func() {
			close(doneChan)
			Eventually(doneAcceptChan, 10*time.Second).Should(Receive())
		})

		Context("when throughput is greater than scaling out threshold", func() {

			BeforeEach(func() {
				policy = GenerateDynamicScaleOutPolicy(1, 2, "throughput", 15)
				initialInstanceCount = 1
			})

			JustBeforeEach(func() {
				// simulate ongoing ~20 requests per second
				ticker = time.NewTicker(50 * time.Millisecond)
				appUri := cfh.AppUri(appName, "/responsetime/fast", cfg)
				go func(chan bool) {
					defer GinkgoRecover()
					for {
						select {
						case <-doneChan:
							ticker.Stop()
							doneAcceptChan <- true
							return
						case <-ticker.C:
							wg.Add(1)
							go func() {
								defer wg.Done()
								cfh.Curl(cfg, appUri)
							}()
						}
					}
				}(doneChan)
			})

			It("should scale out", func() {
				WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)
			})

			AfterEach(func() {
				// ginkgo requires all go-routines to be finished before reporting the result.
				// let's wait for all curl executing go-routines to return.
				wg.Wait()
			})

		})

		Context("when throughput is less than scaling in threshold", func() {

			BeforeEach(func() {
				policy = GenerateDynamicScaleInPolicyBetween("throughput", 15, 25)
				initialInstanceCount = 2
			})

			JustBeforeEach(func() {
				// simulate ongoing ~20 requests per second
				ticker = time.NewTicker(50 * time.Millisecond)
				appUri := cfh.AppUri(appName, "/responsetime/fast", cfg)
				go func(chan bool) {
					defer GinkgoRecover()
					for {
						select {
						case <-doneChan:
							ticker.Stop()
							doneAcceptChan <- true
							return
						case <-ticker.C:
							wg.Add(1)
							go func() {
								defer wg.Done()
								cfh.Curl(cfg, appUri)
							}()
						}
					}
				}(doneChan)
			})

			It("should scale in", func() {
				WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
			})

			AfterEach(func() {
				// ginkgo requires all go-routines to be finished before reporting the result.
				// let's wait for all curl executing go-routines to return.
				wg.Wait()
			})
		})
	})

	// To check existing aggregated cpu metrics do: cf asm APP_NAME cpu
	Context("when scaling by cpu", func() {

		BeforeEach(func() {
			policy = GenerateDynamicScaleOutAndInPolicy(1, 2, "cpu", int64(float64(cfg.CPUUpperThreshold)*0.2), int64(float64(cfg.CPUUpperThreshold)*0.4))
			initialInstanceCount = 1
		})

		It("when cpu is greater than scaling out threshold", func() {
			By("should scale out to 2 instances")
			StartCPUUsage(cfg, appName, int(float64(cfg.CPUUpperThreshold)*0.9), 5)
			WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)

			By("should scale in to 1 instance after cpu usage is reduced")
			//only hit the one instance that was asked to run hot.
			StopCPUUsage(cfg, appName, 0)

			WaitForNInstancesRunning(appGUID, 1, 10*time.Minute)
		})
	})

	Context("when there is a scaling policy for cpuutil", func() {
		BeforeEach(func() {
			policy = GenerateDynamicScaleOutAndInPolicy(1, 2, "cpuutil", 40, 80)
			initialInstanceCount = 1
		})

		It("should scale out and in", func() {
			// this test depends on
			//   - Diego cell size (CPU and RAM)
			//   - CPU entitlements per share configured in ci/operations/set-cpu-entitlement-per-share.yaml
			//   - app memory configured via cfg.CPUUtilScalingPolicyTest.AppMemory
			//   - app CPU entitlement configured via cfg.CPUUtilScalingPolicyTest.AppCPUEntitlement
			//
			// the following gives an example how to calculate an app CPU entitlement:
			//   - Diego cell size = 8 CPU 32GB RAM
			//   - total shares = 1024 * 32[GB host ram] / 8[upper limit of app memory in GB] = 4096
			//   - CPU entitlement per share = 8[number host CPUs] * 100/ 4096[total shares] = 0,1953%
			//   - app memory = 1GB
			//   - app CPU entitlement = 4096[total shares] / (32[GB host ram] * 1024) * (1[app memory in GB] * 1024) * 0,1953 ~= 25%

			ScaleMemory(cfg, appName, cfg.CPUUtilScalingPolicyTest.AppMemory)

			// cpuutil will be 100% if cpu usage is reaching the value of cpu entitlement
			maxCPUUsage := cfg.CPUUtilScalingPolicyTest.AppCPUEntitlement
			StartCPUUsage(cfg, appName, maxCPUUsage, 5)
			WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)

			// only hit the one instance that was asked to run hot
			StopCPUUsage(cfg, appName, 0)
			WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
		})
	})

	Context("when there is a scaling policy for diskutil", func() {
		BeforeEach(func() {
			policy = GenerateDynamicScaleOutAndInPolicy(1, 2, "diskutil", 30, 60)
			initialInstanceCount = 1
		})

		It("should scale out and in", func() {
			ScaleDisk(cfg, appName, "1GB")

			StartDiskUsage(cfg, appName, 800, 5)
			WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)

			// only hit the one instance that was asked to occupy disk space
			StopDiskUsage(cfg, appName, 0)
			WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
		})
	})

	Context("when there is a scaling policy for disk", func() {
		BeforeEach(func() {
			policy = GenerateDynamicScaleOutAndInPolicy(1, 2, "disk", 300, 600)
			initialInstanceCount = 1
		})

		It("should scale out and in", func() {
			ScaleDisk(cfg, appName, "1GB")

			StartDiskUsage(cfg, appName, 800, 5)
			WaitForNInstancesRunning(appGUID, 2, 5*time.Minute)

			// only hit the one instance that was asked to occupy disk space
			StopDiskUsage(cfg, appName, 0)
			WaitForNInstancesRunning(appGUID, 1, 5*time.Minute)
		})
	})
})

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}
