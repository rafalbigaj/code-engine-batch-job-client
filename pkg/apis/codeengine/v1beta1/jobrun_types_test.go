/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
)

var _ = Describe("Jobrun Types", func() {
	Context("when UpdateFailedIndices", func() {
		var jobRun *cv1b1.JobRun
		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{
				TypeMeta: metav1.TypeMeta{Kind: "JobRun", APIVersion: cv1b1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "unit-test-utils",
					Name:      "fake-job-run",
				},
			}
			jobRun.Spec.JobDefinitionSpec.Template.Containers = []corev1.Container{{}}
		})
		Describe("when UpdateFailedIndices is called", func() {
			It("calculates the correct failed indices for one not yet created pod", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{0: corev1.PodUnknown}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(*jr.Status.FailedIndices).To(Equal("0"))

			})

			It("calculates the correct failed indices for one failed pod", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{0: corev1.PodFailed}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(*jr.Status.FailedIndices).To(Equal("0"))

			})

			It("calculates the correct failed indices for one succeeded pod", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{0: corev1.PodSucceeded}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(jr.Status.FailedIndices).To(BeNil())

			})

			It("calculates the correct failed indices for two failed and one succeeded pod", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{
					0: corev1.PodSucceeded,
					3: corev1.PodFailed,
					4: corev1.PodFailed,
				}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(*jr.Status.FailedIndices).To(Equal("3-4"))

			})

			It("calculates the correct failed indices for two failed and two succeeded pod in non-consecutive order", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{
					5:  corev1.PodSucceeded,
					18: corev1.PodFailed,
					3:  corev1.PodSucceeded,
					17: corev1.PodFailed,
				}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(*jr.Status.FailedIndices).To(Equal("17-18"))
			})

			It("calculates the correct failed indices for all succeeded pods in non-consecutive order", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{
					5:  corev1.PodSucceeded,
					19: corev1.PodSucceeded,
					66: corev1.PodSucceeded,
					42: corev1.PodSucceeded,
					28: corev1.PodSucceeded,
				}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(jr.Status.FailedIndices).To(BeNil())

			})

			It("calculates the correct failed indices for all non-succeeded pods in arbitrary order", func() {
				jr := jobRun.DeepCopy()
				podSnapshots := map[int64]corev1.PodPhase{
					5:  corev1.PodFailed,
					19: corev1.PodPending,
					66: corev1.PodFailed,
					42: corev1.PodFailed,
					28: corev1.PodRunning,
					41: corev1.PodFailed,
					43: corev1.PodUnknown,
					18: corev1.PodFailed,
				}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(*jr.Status.FailedIndices).To(Equal("5,18-19,28,41-43,66"))
			})

			It("will not count failed indices when running in daemon mode", func() {
				jr := jobRun.DeepCopy()
				jr.Spec.JobDefinitionSpec.Template.Containers = []corev1.Container{
					{
						Env: []corev1.EnvVar{
							{
								Name:  cv1b1.CEExecutionMode,
								Value: cv1b1.CEExecutionModeValue,
							},
						},
					},
				}

				podSnapshots := map[int64]corev1.PodPhase{
					5:  corev1.PodFailed,
					19: corev1.PodPending,
					66: corev1.PodFailed,
					42: corev1.PodFailed,
					28: corev1.PodRunning,
					41: corev1.PodFailed,
					43: corev1.PodUnknown,
					18: corev1.PodFailed,
				}
				jr.UpdateFailedIndices(podSnapshots)
				Expect(jr.Status.FailedIndices).To(BeNil())
			})
		})
	})

	Context("when UpdateSucceededIndices", func() {
		var jobRun *cv1b1.JobRun
		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{
				TypeMeta: metav1.TypeMeta{Kind: "JobRun", APIVersion: cv1b1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "unit-test-utils",
					Name:      "fake-job-run",
				},
			}
		})

		It("calculates the correct succeeded indices for one succeeded pod", func() {
			jr := jobRun.DeepCopy()
			podSnapshots := map[int64]corev1.PodPhase{0: corev1.PodSucceeded}
			jr.UpdateSucceededIndices(podSnapshots)
			Expect(*jr.Status.SucceededIndices).To(Equal("0"))

		})

		It("calculates the correct succeeded indices for one failed pod", func() {
			jr := jobRun.DeepCopy()
			podSnapshots := map[int64]corev1.PodPhase{0: corev1.PodFailed}
			jr.UpdateSucceededIndices(podSnapshots)
			Expect(jr.Status.SucceededIndices).To(BeNil())

		})

		It("calculates the correct succeeded indices for one failed and two succeeded pod", func() {
			jr := jobRun.DeepCopy()
			podSnapshots := map[int64]corev1.PodPhase{
				0: corev1.PodFailed,
				3: corev1.PodSucceeded,
				4: corev1.PodSucceeded,
			}
			jr.UpdateSucceededIndices(podSnapshots)
			Expect(*jr.Status.SucceededIndices).To(Equal("3-4"))

		})

		It("calculates the correct succeeded indices for two failed and two succeeded pod in non-consecutive order", func() {
			jr := jobRun.DeepCopy()
			podSnapshots := map[int64]corev1.PodPhase{
				5:  corev1.PodSucceeded,
				18: corev1.PodFailed,
				3:  corev1.PodSucceeded,
				17: corev1.PodFailed,
			}
			jr.UpdateSucceededIndices(podSnapshots)
			Expect(*jr.Status.SucceededIndices).To(Equal("3,5"))

		})
	})

	Context("when UpdateStatusCounts", func() {
		var (
			jobRun    *cv1b1.JobRun
			podStatus []corev1.PodPhase
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{
				TypeMeta: metav1.TypeMeta{Kind: "JobRun", APIVersion: cv1b1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "unit-test-utils",
					Name:      "fake-job-run",
				},
			}
			podStatus = []corev1.PodPhase{
				corev1.PodPending,
				corev1.PodRunning,
				corev1.PodSucceeded,
				corev1.PodFailed,
				corev1.PodUnknown,
			}
		})

		It("counts job pods status", func() {
			jr := jobRun.DeepCopy()
			total := int64(5)

			jr.UpdateStatusCounts(total, podStatus)
			Expect(jr.Status.Pending).To(Equal(int64(1)))
			Expect(jr.Status.Running).To(Equal(int64(1)))
			Expect(jr.Status.Succeeded).To(Equal(int64(1)))
			Expect(jr.Status.Failed).To(Equal(int64(1)))
			Expect(jr.Status.Unknown).To(Equal(int64(1)))
			Expect(jr.Status.Requested).To(Equal(int64(0)))
		})

		It("counts job pods with uncreated pods", func() {
			jr := jobRun.DeepCopy()
			total := int64(6)

			jr.UpdateStatusCounts(total, podStatus)
			Expect(jr.Status.Pending).To(Equal(int64(1)))
			Expect(jr.Status.Running).To(Equal(int64(1)))
			Expect(jr.Status.Succeeded).To(Equal(int64(1)))
			Expect(jr.Status.Failed).To(Equal(int64(1)))
			Expect(jr.Status.Unknown).To(Equal(int64(1)))
			Expect(jr.Status.Requested).To(Equal(int64(1)))
		})
	})

	Context("handling conditions", func() {
		var (
			jobRunStatus    *cv1b1.JobRunStatus
			jobrunCondition cv1b1.JobRunCondition
		)
		BeforeEach(func() {
			jobRunStatus = &cv1b1.JobRunStatus{}
			jobrunCondition = cv1b1.JobRunCondition{
				Type: cv1b1.JobPending,
			}
		})

		When("there are no conditions for the jobrun", func() {
			It("cannot get the desired condition", func() {
				Expect(jobRunStatus.GetCondition(cv1b1.JobComplete)).To(BeNil())
			})
			It("cannot get the latest condition", func() {
				Expect(jobRunStatus.GetLatestCondition()).To(BeNil())
			})
			It("can add the condition", func() {
				jobRunStatus.AddCondition(jobrunCondition)
				Expect(len(jobRunStatus.Conditions)).To(Equal(1))
			})
		})
		When("there exist one pending conditions for the jobrun", func() {
			BeforeEach(func() {
				jobRunStatus.Conditions = []cv1b1.JobRunCondition{jobrunCondition}
			})
			It("can get the desired condition", func() {
				Expect(jobRunStatus.GetCondition(cv1b1.JobPending)).To(Equal(&jobrunCondition))
			})
			It("can get the latest condition", func() {
				Expect(jobRunStatus.GetLatestCondition()).To(Equal(&jobrunCondition))
			})
			It("ignore the same the condition", func() {
				jobRunStatus.AddCondition(jobrunCondition)
				Expect(len(jobRunStatus.Conditions)).To(Equal(1))
			})
		})

	})
	Context("JobRunDaemonMode", func() {
		var (
			jobRun      *cv1b1.JobRun
			envVarName  string
			envVarValue string
		)

		JustBeforeEach(func() {
			jobRun = &cv1b1.JobRun{
				TypeMeta: metav1.TypeMeta{Kind: "JobRun", APIVersion: cv1b1.SchemeGroupVersion.String()},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "unit-test-utils",
					Name:      "fake-job-run",
				},
				Spec: cv1b1.JobRunSpec{
					JobDefinitionSpec: cv1b1.JobDefinitionSpec{
						Template: cv1b1.JobPodTemplate{
							Containers: []corev1.Container{
								{
									Name:  "bla",
									Image: "busybox",
									Env: []corev1.EnvVar{
										{
											Name:  envVarName,
											Value: envVarValue,
										},
									},
								},
							},
						},
					},
				},
			}
		})

		Context("with a JobRun that specifies CE Execution Mode", func() {
			BeforeEach(func() {
				envVarName = cv1b1.CEExecutionMode
				envVarValue = cv1b1.CEExecutionModeValue
			})

			It("enables daemon mode", func() {
				jr := jobRun.DeepCopy()
				Expect(jr.IsRunningInDaemonMode()).To(BeTrue())
			})
		})

		Context("with a JobRun that does not specify CE Execution Mode", func() {
			BeforeEach(func() {
				envVarName = "someFakeKey"
				envVarValue = "someFakeValue"
			})

			It("doesn't enable daemon mode", func() {
				jr := jobRun.DeepCopy()
				Expect(jr.IsRunningInDaemonMode()).To(BeFalse())
			})
		})

		Context("with a JobRun that specifies correct CE Execution Mode value, incorrect key", func() {
			BeforeEach(func() {
				envVarName = "someFakeKey"
				envVarValue = cv1b1.CEExecutionModeValue
			})

			It("doesn't enable daemon mode", func() {
				jr := jobRun.DeepCopy()
				Expect(jr.IsRunningInDaemonMode()).To(BeFalse())
			})
		})

		Context("with a JobRun that specifies correct CE Execution Mode key, incorrect value", func() {
			BeforeEach(func() {
				envVarName = cv1b1.CEExecutionMode
				envVarValue = "someFakeValue"
			})

			It("doesn't enable daemon mode", func() {
				jr := jobRun.DeepCopy()
				Expect(jr.IsRunningInDaemonMode()).To(BeFalse())
			})
		})
	})
})
