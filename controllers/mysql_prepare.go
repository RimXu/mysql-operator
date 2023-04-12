package controllers

import (
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	api "mysql-operator/pkg/api"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type MysqlPrepare struct {
}

func (d MysqlPrepare) Create(evt event.CreateEvent) bool {
	//fmt.Println("Create")
	//fmt.Println(evt.Object.(*appsv1.Deployment).ObjectMeta.Namespace)
	//fmt.Println(evt.Object.GetName())
	//evt.Object.(*nginxv1.Nginx).ObjectMeta
	// 只有return true 才能将数据加入到workqueue中
	return true
}

func (d MysqlPrepare) Delete(evt event.DeleteEvent) bool {
	return true
}

// Watches Job Update事件,通过事件状态更新CRD Spec.phase
func (d MysqlPrepare) Update(evt event.UpdateEvent) bool {
	objs := evt.ObjectNew.GetOwnerReferences()
	for id := range objs {
		evt_kind := objs[id].Kind
		var job_status string
		if evt_kind == constants.Kind {
			// 断言类型
			job, _ := evt.ObjectNew.(*batchv1.Job)
			if job != nil {
				job_name := job.Name
				job_namespace := job.Namespace
				job_owner := objs[id].Name
				job_status_suc := job.Status.Succeeded
				job_status_act := job.Status.Active
				job_status_fai := job.Status.Failed
				if job_status_suc > 0 {
					job_status = "JobSuccess"
				} else if job_status_act > 0 {
					job_status = "JobRunning"
				} else if job_status_fai > 0 {
					job_status = "JobFailed"
				}
				err := api.UpdateMysqlStatus(job_namespace, job_owner, job_status)
				if err != nil {
					logrus.Errorf("Status update failed: { namespace:'%s', name:'%s' }", job_namespace, job_name)
				}
			}
		}
	}
	return true
}

func (d MysqlPrepare) Generic(evt event.GenericEvent) bool {
	return false
}
