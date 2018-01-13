#!/bin/bash

## change workerdir
cd /home/worker/scripts

## get ip address
ipv4=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)

###modifi kubelet config
sed -i "s/0.0.0.0/$ipv4/g" /etc/sysconfig/kubernetes/kubelet

## modify kubelet-proxy config
sed -i "s/0.0.0.0/$ipv4/g" /etc/sysconfig/kubernetes/proxy

### reload systemctl
systemctl daemon-reload

## restart flanneld
systemctl enable flanneld && systemctl restart flanneld
## restart docker
systemctl enable docker && systemctl restart docker 
## restart kubelet
systemctl enable kubelet && systemctl restart kubelet && systemctl status kubelet
## restart kube-proxy
systemctl enable kube-proxy && systemctl restart kube-proxy && systemctl status kube-proxy
