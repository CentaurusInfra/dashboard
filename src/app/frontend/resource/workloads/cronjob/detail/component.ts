// Copyright 2017 The Kubernetes Authors.
// Copyright 2020 Authors of Arktos - file modified.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {CronJobDetail} from '@api/backendapi';
import {Subscription} from 'rxjs/Subscription';

import {ActionbarService, ResourceMeta} from '../../../../common/services/global/actionbar';
import {NotificationsService} from '../../../../common/services/global/notifications';
import {EndpointManager, Resource} from '../../../../common/services/resource/endpoint';
import {NamespacedResourceService} from '../../../../common/services/resource/resource';

@Component({
  selector: 'kd-cron-job-detail',
  templateUrl: './template.html',
})
export class CronJobDetailComponent implements OnInit, OnDestroy {
  private cronJobSubscription_: Subscription;
  private readonly endpoint_ = EndpointManager.resource(Resource.cronJob, true, true);
  cronJob: CronJobDetail;
  isInitialized = false;
  eventListEndpoint: string;
  activeJobsEndpoint: string;
  inactiveJobsEndpoint: string;

  constructor(
    private readonly cronJob_: NamespacedResourceService<CronJobDetail>,
    private readonly actionbar_: ActionbarService,
    private readonly activatedRoute_: ActivatedRoute,
    private readonly notifications_: NotificationsService,
  ) {}

  ngOnInit(): void {
    const resourceName = this.activatedRoute_.snapshot.params.resourceName;
    const resourceNamespace = this.activatedRoute_.snapshot.params.resourceNamespace;

    this.eventListEndpoint = this.endpoint_.child(resourceName, Resource.event, resourceNamespace);
    this.activeJobsEndpoint = this.endpoint_.child(resourceName, Resource.job, resourceNamespace);
    this.inactiveJobsEndpoint = this.activeJobsEndpoint + `?active=false`;

    this.cronJobSubscription_ = this.cronJob_
      .get(this.endpoint_.detail(), resourceName, resourceNamespace)
      .subscribe((d: CronJobDetail) => {
        this.cronJob = d;
        this.notifications_.pushErrors(d.errors);
        this.actionbar_.onInit.emit(new ResourceMeta('Cron Job', d.objectMeta, d.typeMeta));
        this.isInitialized = true;
      });
  }

  ngOnDestroy(): void {
    this.cronJobSubscription_.unsubscribe();
    this.actionbar_.onDetailsLeave.emit();
  }
}
