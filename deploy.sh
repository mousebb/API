#!/bin/bash

GROUP=api-group
TEMPLATE_BASE=api-template
ZONE=us-central1-a

function start_docker(){
	if ! boot2docker status | grep -q running
	then
		boot2docker start
		DOCKER_HOST="tcp://$(boot2docker config | grep LowerIP | awk 'sub(/LowerIP = "/, ""){print $1}' | awk 'sub(/"/,""){print $1}'):2376"
		DOCKER_CERT_PATH="$(boot2docker config | grep Dir | awk 'sub(/Dir = "/, ""){print $1}' | awk 'sub(/"/,""){print $1}')/certs/boot2docker-vm"
		DOCKER_TLS_VERIFY=1
	fi
}

function push_image(){
	docker build -t gcr.io/curt_groups/api .
	gcloud preview docker push gcr.io/curt_groups/api
}

function restart_templates(){

	current_time=$(date "+%Y-%m-%d-%H-%M-%S")
	new_template=$(echo "$TEMPLATE_BASE-$current_time")

	old_instance_name=$(gcloud compute instance-templates list | awk 'NR>1 {print $1}')
	gcloud compute instance-templates describe $old_instance_name | shyaml get-value properties.metadata.items.0.value >> manifest.yaml
	gcloud compute instance-templates create $new_template --image container-vm --metadata-from-file google-container-manifest=manifest.yaml --scopes https://www.googleapis.com/auth/devstorage.full_control https://www.googleapis.com/auth/compute https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/replicapool

	gcloud preview managed-instance-groups --zone $ZONE set-template $GROUP --template $new_template
	gcloud preview rolling-updates --zone $ZONE start --group $GROUP --template $new_template

	old_templates=$(gcloud compute instance-templates list | awk '$1 ~ /$TEMPLATE_BASE/ { print $1}' | grep -v $new_template | tr " " "\n")
	for tmp in $old_templates
	do
		gcloud compute instance-templates delete $tmp --quiet
	done

	rm -rf manifest.yaml

}

push_image
restart_templates