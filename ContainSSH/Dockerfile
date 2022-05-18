FROM alpine:3.15.2
RUN rm /etc/motd
ADD https://storage.googleapis.com/kubernetes-release/release/v1.22.0/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

# Setup OpenSSH
RUN apk add --update openssh \
    && /usr/bin/ssh-keygen -A \
    && ssh-keygen -t rsa -b 4096 -f  /etc/ssh/ssh_host_key

ENTRYPOINT adduser $USERNAME -D \
           && echo "$USERNAME:$PASSWORD" | chpasswd \
           && kubectl config set-cluster cfc --server=https://kubernetes.default --certificate-authority=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
           && kubectl config set-context cfc --cluster=cfc \
           && kubectl config set-credentials user --token=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token) \
           && kubectl config set-context cfc --user=user \
           && kubectl config use-context cfc \
           && mv /root/.kube /home/$USERNAME/.kube \
           && chown -R $USERNAME:$USERNAME /home/$USERNAME/.kube \
           && echo "kubectl exec -it $USERNAME-0 -n $(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace) -c main -- $CONTAINERSHELL; exit" >> /home/$USERNAME/.profile \
           && /usr/sbin/sshd -D
