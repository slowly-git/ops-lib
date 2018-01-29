# 运行前请先手动关闭selinux并重启服务器
sed -i "s@SELINUX=enforcing@SELINUX=disabled@g" /etc/selinux/config 
# 请不要设置swap分区
################################初始环境部分################################
yum update -y
systemctl stop firewalld.service
systemctl disable  firewalld.service
################################内核优化部分################################
#修改最大文件数
echo '* soft nofile 32768' >> /etc/security/limits.conf
echo '* hard nofile 65536' >> /etc/security/limits.conf
echo 'mysql soft nofile 65535' >> /etc/security/limits.conf
echo 'mysql hard nofile 65535' >> /etc/security/limits.conf
#内核参数调优
echo 'net.ipv4.tcp_syn_retries = 1' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_synack_retries = 1' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_keepalive_time = 600' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_keepalive_probes = 3' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_keepalive_intvl =15' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_retries2 = 5' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_fin_timeout = 2' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_max_tw_buckets = 36000' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_tw_recycle = 0' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_tw_reuse = 1' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_max_orphans = 32768' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_max_syn_backlog = 16384' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_mem = 8388608 8388608 8388608' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_rmem = 4096 87380 8388608' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_wmem = 4096 65535 8388608' >> /etc/sysctl.conf
echo 'net.core.rmem_max = 16777216' >> /etc/sysctl.conf
echo 'net.core.wmem_max = 16777216' >> /etc/sysctl.conf
echo 'net.core.wmem_default = 8388608' >> /etc/sysctl.conf
echo 'net.core.rmem_default = 8388608' >> /etc/sysctl.conf
echo 'net.core.optmem_max = 40960' >> /etc/sysctl.conf
echo 'net.core.netdev_max_backlog = 3000' >> /etc/sysctl.conf
echo 'net.ipv4.ip_local_port_range = 1024 65000' >> /etc/sysctl.conf
echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf
echo 'net.ipv4.ip_forward_use_pmtu = 0' >> /etc/sysctl.conf
sysctl -p

################################docker 安装################################
yum remove docker docker-common docker-selinux docker-engine git -y
yum install -y yum-utils device-mapper-persistent-data lvm2 wget
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum-config-manager --disable docker-ce-edge
yum install docker-ce -y
systemctl start docker
systemctl enable docker
docker pull cn-bj-public-ops-harbor-1.360in.com/pinguo-open/pause-amd64:3.0

################################k8s 安装################################
wget https://s3.cn-north-1.amazonaws.com.cn/pinguo-dev/g-ops-client/pingup_k8s_client.tgz
mkdir /var/lib/kubelet
tar -zxvPf pingup_k8s_client.tgz
systemctl daemon-reload

HOST=$(echo cn-kanny-k8s-client-$RANDOM)
sed -i "s@--hostname-override=localhost@--hostname-override=$HOST@g" /etc/sysconfig/kubernetes/kubelet 

systemctl start kubelet
systemctl enable kubelet
