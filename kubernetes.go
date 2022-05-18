package main

import (
	"bytes"
	"context"
	"text/template"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type ContainerDefinition struct {
	Name     string
	Image    string
	Password string
	Shell    string
}

func (h *Handler) getImageData(imageName string) (Image, error) {
	var image Image
	if result := h.db.Where("image_name = ?", imageName).First(&image); result.Error != nil {
		return Image{}, result.Error
	}
	return image, nil
}

func getResources(username string, h *Handler) (int, error) {
	serviceClient := h.client.CoreV1().Services(h.namespace)
	service, err := serviceClient.Get(context.TODO(), username, metav1.GetOptions{})
	return int(service.Spec.Ports[0].NodePort), err
}

func createResource(filename string, containerDefinition ContainerDefinition, data interface{}) {
	template, err := template.New(filename).ParseFiles("templates/" + filename)
	if err != nil {
		panic(err)
	}

	var manifest bytes.Buffer
	if err := template.Execute(&manifest, containerDefinition); err != nil {
		panic(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(&manifest, 10000).Decode(data)
	if err != nil {
		panic(err)
	}

}

func createResources(username string, imageName string, password string, h *Handler) (*Container, error) {

	// Gather image info from database
	image, err := h.getImageData(imageName)
	if err != nil {
		return &Container{}, err
	}

	containerDefinition := ContainerDefinition{
		Name:     username,
		Image:    image.Repository,
		Password: password,
		Shell:    image.Shell,
	}

	role := &rbac.Role{}
	createResource("role.yaml", containerDefinition, role)
	roleClient := h.client.RbacV1().Roles(h.namespace)
	_, err = roleClient.Create(context.TODO(), role, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	roleBinding := &rbac.RoleBinding{}
	createResource("rolebinding.yaml", containerDefinition, roleBinding)
	roleBindingClient := h.client.RbacV1().RoleBindings(h.namespace)
	_, err = roleBindingClient.Create(context.TODO(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	serviceAccount := &corev1.ServiceAccount{}
	createResource("serviceaccount.yaml", containerDefinition, serviceAccount)
	serviceAccountClient := h.client.CoreV1().ServiceAccounts(h.namespace)
	_, err = serviceAccountClient.Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	secret := &corev1.Secret{}
	createResource("secret.yaml", containerDefinition, secret)
	secretClient := h.client.CoreV1().Secrets(h.namespace)
	_, err = secretClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	statefulSetSSH := &appsv1.StatefulSet{}
	statefulSet := &appsv1.StatefulSet{}
	createResource("StatefulSet-ssh.yaml", containerDefinition, statefulSetSSH)
	createResource("StatefulSet.yaml", containerDefinition, statefulSet)
	statefulSetClient := h.client.AppsV1().StatefulSets(h.namespace)
	_, err = statefulSetClient.Create(context.TODO(), statefulSetSSH, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	_, err = statefulSetClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	service := &corev1.Service{}
	createResource("service.yaml", containerDefinition, service)
	serviceClient := h.client.CoreV1().Services(h.namespace)
	_, err = serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	newContainer := &Container{
		Username:  username,
		ImageName: imageName,
	}

	return newContainer, nil

}

func deleteResources(username string, h *Handler) error {
	roleClient := h.client.RbacV1().Roles(h.namespace)
	err := roleClient.Delete(context.TODO(), username, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	roleBindingClient := h.client.RbacV1().RoleBindings(h.namespace)
	err = roleBindingClient.Delete(context.TODO(), username, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	serviceAccountClient := h.client.CoreV1().ServiceAccounts(h.namespace)
	err = serviceAccountClient.Delete(context.TODO(), username, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	serviceClient := h.client.CoreV1().Services(h.namespace)
	err = serviceClient.Delete(context.TODO(), username, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	statefulSetClient := h.client.AppsV1().StatefulSets(h.namespace)
	err = statefulSetClient.Delete(context.TODO(), username, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = statefulSetClient.Delete(context.TODO(), username+"-ssh", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	secretClient := h.client.CoreV1().Secrets(h.namespace)
	err = secretClient.Delete(context.TODO(), username, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
