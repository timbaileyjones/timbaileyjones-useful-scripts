#!/usr/bin/env python
import calendar
import datetime as dt
import os
import time

import boto3
import jsondate

cp_client = boto3.client('codepipeline')
cb_client = boto3.client('codebuild')

list_pipelines_paginator = cp_client.get_paginator('list_pipelines')
list_pipelines_iterator = list_pipelines_paginator.paginate()

for list_pipelines_page in list_pipelines_iterator:
    print('list_pipelines_page', jsondate.dumps(list_pipelines_page, indent=2))

    for pipeline in list_pipelines_page['pipelines']:
        name = pipeline['name']
        get_pipeline_response = cp_client.get_pipeline(name=name)
        updated = str(pipeline['updated'])[:19]
        print(f'get_pipeline_response for {name}')
        os.makedirs(name, exist_ok=True)
        filepath = f'{name}/{name}.json'
        with open(filepath, 'w') as fp:
            fp.write(jsondate.dumps(get_pipeline_response, indent=2))

        modify_time = calendar.timegm(time.strptime(updated, '%Y-%m-%d %H:%M:%S'))
        os.utime(filepath, (modify_time, modify_time))
        os.utime(name, (modify_time, modify_time))

        print(jsondate.dumps(get_pipeline_response, indent=2))
        stages = get_pipeline_response['pipeline']['stages']
        for idx, stage in enumerate(stages):
            stage_name = stage['name']
            actions = stage['actions']
            for action in actions:
                if 'actionTypeId' in action:
                    actionTypeId = action['actionTypeId']
                    if 'provider' in actionTypeId and actionTypeId['provider'] == 'CodeBuild':
                        filepath = f'{name}/stage-{idx}-{action["name"]}.json'
                        codebuild_projects = cb_client.batch_get_projects(names=[action['name']])
                        for codebuild_project in codebuild_projects:
                            with open(filepath, 'w') as fp:
                                fp.write(jsondate.dumps(codebuild_project, indent=2))
                            os.utime(filepath, (modify_time, modify_time))
