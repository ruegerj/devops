#!/bin/bash

if ! command -v ansible >/dev/null 2>&1; then
  echo "ansible not installed!"
  exit 127
fi

if ! command -v vagrant >/dev/null 2>&1; then
  echo "vagrant not installed!"
  exit 127
fi

if [[ ! -f ./id_vagrant || ! -f ./id_vagrant.pub ]]; then
    ssh-keygen -t ed25519 -f ./id_vagrant -N "" -q
fi

ansible-galaxy install -r ansible/requirements.yml

vagrant up

# just wait a little bit for thing to calm down
sleep 10

ansible-playbook ansible/site.yml -i ansible/inventory.dev.yml

